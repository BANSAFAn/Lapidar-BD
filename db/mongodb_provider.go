package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBProvider представляет провайдер для работы с MongoDB
type MongoDBProvider struct {
	client     *mongo.Client
	db         *mongo.Database
	reports    *mongo.Collection
	bans       *mongo.Collection
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// Initialize инициализирует соединение с базой данных MongoDB
func (p *MongoDBProvider) Initialize(config DatabaseConfig) error {
	// Формируем строку подключения
	connStr := fmt.Sprintf("mongodb://%s:%s@%s:%d",
		config.User, config.Password, config.Host, config.Port)

	// Создаем контекст с возможностью отмены
	p.ctx, p.cancelFunc = context.WithCancel(context.Background())

	// Настраиваем клиент MongoDB
	clientOptions := options.Client().ApplyURI(connStr)

	// Подключаемся к MongoDB
	client, err := mongo.Connect(p.ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("ошибка подключения к MongoDB: %w", err)
	}

	// Проверяем соединение
	err = client.Ping(p.ctx, nil)
	if err != nil {
		return fmt.Errorf("ошибка проверки соединения с MongoDB: %w", err)
	}

	p.client = client

	// Выбираем базу данных
	dbName := config.Database
	if dbName == "" {
		dbName = "discord_bot"
	}
	p.db = client.Database(dbName)

	// Инициализируем коллекции
	p.reports = p.db.Collection("reports")
	p.bans = p.db.Collection("bans")

	return nil
}

// Close закрывает соединение с базой данных
func (p *MongoDBProvider) Close() error {
	if p.cancelFunc != nil {
		p.cancelFunc()
	}

	if p.client != nil {
		return p.client.Disconnect(context.Background())
	}
	return nil
}

// AddReport добавляет новый репорт в базу данных
func (p *MongoDBProvider) AddReport(reportedUserID, reporterID, reason string) (int64, error) {
	report := bson.M{
		"reported_user_id": reportedUserID,
		"reporter_id":      reporterID,
		"reason":           reason,
		"timestamp":        time.Now(),
		"confirmed":        false,
		"confirmed_by":     "",
	}

	result, err := p.reports.InsertOne(p.ctx, report)
	if err != nil {
		return 0, err
	}

	// Преобразуем ObjectID в int64 для совместимости с интерфейсом
	id := result.InsertedID.(primitive.ObjectID)
	timestamp := id.Timestamp().Unix()
	return timestamp, nil
}

// ConfirmReport подтверждает репорт администратором
func (p *MongoDBProvider) ConfirmReport(reportID int64, adminID string) error {
	// В MongoDB мы используем ObjectID, но для совместимости с интерфейсом
	// мы принимаем int64. Здесь мы ищем по временной метке.
	filter := bson.M{"timestamp": time.Unix(reportID, 0)}
	update := bson.M{
		"$set": bson.M{
			"confirmed":    true,
			"confirmed_by": adminID,
		},
	}

	_, err := p.reports.UpdateOne(p.ctx, filter, update)
	return err
}

// GetReportsByUser получает все репорты на указанного пользователя
func (p *MongoDBProvider) GetReportsByUser(userID string) ([]Report, error) {
	filter := bson.M{"reported_user_id": userID}

	cursor, err := p.reports.Find(p.ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(p.ctx)

	var reports []Report
	for cursor.Next(p.ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}

		var report Report
		report.ReportedUserID = doc["reported_user_id"].(string)
		report.ReporterID = doc["reporter_id"].(string)
		report.Reason = doc["reason"].(string)
		report.Timestamp = doc["timestamp"].(primitive.DateTime).Time()
		report.Confirmed = doc["confirmed"].(bool)

		if confirmedBy, ok := doc["confirmed_by"].(string); ok && confirmedBy != "" {
			report.ConfirmedBy = confirmedBy
		}

		// Используем временную метку как ID для совместимости
		report.ID = report.Timestamp.Unix()

		reports = append(reports, report)
	}

	return reports, nil
}

// GetReportCount получает количество подтвержденных репортов на пользователя
func (p *MongoDBProvider) GetReportCount(userID string) (int, error) {
	filter := bson.M{
		"reported_user_id": userID,
		"confirmed":        true,
	}

	count, err := p.reports.CountDocuments(p.ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// AddBan добавляет новый бан в базу данных
func (p *MongoDBProvider) AddBan(userID, reason, adminID string, duration *time.Duration) error {
	var expiresAt *time.Time
	if duration != nil {
		expires := time.Now().Add(*duration)
		expiresAt = &expires
	}

	ban := bson.M{
		"user_id":    userID,
		"reason":     reason,
		"admin_id":   adminID,
		"timestamp":  time.Now(),
		"expires_at": expiresAt,
	}

	_, err := p.bans.InsertOne(p.ctx, ban)
	return err
}

// GetActiveBan проверяет, есть ли активный бан у пользователя
func (p *MongoDBProvider) GetActiveBan(userID string) (*Ban, error) {
	filter := bson.M{
		"user_id": userID,
		"$or": []bson.M{
			{"expires_at": nil},
			{"expires_at": bson.M{"$gt": time.Now()}},
		},
	}

	var doc bson.M
	err := p.bans.FindOne(p.ctx, filter).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var ban Ban
	ban.UserID = doc["user_id"].(string)
	ban.Reason = doc["reason"].(string)
	ban.AdminID = doc["admin_id"].(string)
	ban.Timestamp = doc["timestamp"].(primitive.DateTime).Time()

	// Используем временную метку как ID для совместимости
	ban.ID = ban.Timestamp.Unix()

	if expiresAt, ok := doc["expires_at"].(primitive.DateTime); ok {
		expTime := expiresAt.Time()
		ban.ExpiresAt = &expTime
	}

	return &ban, nil
}

// GetType возвращает тип базы данных
func (p *MongoDBProvider) GetType() string {
	return "mongodb"
}
