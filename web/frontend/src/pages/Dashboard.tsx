import React, { useState, useEffect } from 'react';
import {
  Grid,
  Paper,
  Typography,
  Box,
  Card,
  CardContent,
  CardHeader,
  Divider,
  List,
  ListItem,
  ListItemText,
  ListItemAvatar,
  Avatar,
  CircularProgress,
} from '@mui/material';
import {
  People as PeopleIcon,
  Forum as ForumIcon,
  Mic as MicIcon,
  Memory as MemoryIcon,
} from '@mui/icons-material';
import axios from 'axios';

// Интерфейс для статистики бота
interface BotStats {
  servers: number;
  users: number;
  channels: number;
  commands: number;
  uptime: string;
  memoryUsage: string;
}

const Dashboard: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState<BotStats>({
    servers: 0,
    users: 0,
    channels: 0,
    commands: 0,
    uptime: '0 часов',
    memoryUsage: '0 MB',
  });

  // В реальном приложении здесь будет запрос к API
  useEffect(() => {
    // Имитация загрузки данных
    const timer = setTimeout(() => {
      setStats({
        servers: 15,
        users: 1250,
        channels: 87,
        commands: 432,
        uptime: '24 часа 15 минут',
        memoryUsage: '128 MB',
      });
      setLoading(false);
    }, 1000);

    return () => clearTimeout(timer);
  }, []);

  // Данные для статистических карточек
  const statCards = [
    { title: 'Серверов', value: stats.servers, icon: <PeopleIcon fontSize="large" color="primary" /> },
    { title: 'Пользователей', value: stats.users, icon: <PeopleIcon fontSize="large" color="secondary" /> },
    { title: 'Каналов', value: stats.channels, icon: <ForumIcon fontSize="large" color="info" /> },
    { title: 'Команд', value: stats.commands, icon: <MicIcon fontSize="large" color="success" /> },
  ];

  // Последние события (в реальном приложении будут загружаться с сервера)
  const recentEvents = [
    { id: 1, type: 'command', user: 'User1', content: '/help', time: '5 минут назад' },
    { id: 2, type: 'join', user: 'User2', content: 'присоединился к серверу', time: '10 минут назад' },
    { id: 3, type: 'message', user: 'User3', content: 'отправил сообщение', time: '15 минут назад' },
    { id: 4, type: 'voice', user: 'User4', content: 'присоединился к голосовому каналу', time: '20 минут назад' },
  ];

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100%' }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Панель управления
      </Typography>
      
      {/* Статистические карточки */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        {statCards.map((card, index) => (
          <Grid item xs={12} sm={6} md={3} key={index}>
            <Paper elevation={2} sx={{ p: 2, height: '100%', backgroundColor: '#2f3136' }}>
              <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                <Box>
                  <Typography variant="h3" component="div">
                    {card.value}
                  </Typography>
                  <Typography variant="subtitle1" color="text.secondary">
                    {card.title}
                  </Typography>
                </Box>
                {card.icon}
              </Box>
            </Paper>
          </Grid>
        ))}
      </Grid>

      {/* Информация о системе и последние события */}
      <Grid container spacing={3}>
        <Grid item xs={12} md={6}>
          <Card sx={{ backgroundColor: '#2f3136' }}>
            <CardHeader title="Информация о системе" />
            <Divider sx={{ backgroundColor: '#40444b' }} />
            <CardContent>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
                <Typography variant="body1">Время работы:</Typography>
                <Typography variant="body1">{stats.uptime}</Typography>
              </Box>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
                <Typography variant="body1">Использование памяти:</Typography>
                <Typography variant="body1">{stats.memoryUsage}</Typography>
              </Box>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
                <Typography variant="body1">Версия бота:</Typography>
                <Typography variant="body1">1.0.0</Typography>
              </Box>
              <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                <Typography variant="body1">Статус:</Typography>
                <Typography variant="body1" color="success.main">Онлайн</Typography>
              </Box>
            </CardContent>
          </Card>
        </Grid>
        
        <Grid item xs={12} md={6}>
          <Card sx={{ backgroundColor: '#2f3136' }}>
            <CardHeader title="Последние события" />
            <Divider sx={{ backgroundColor: '#40444b' }} />
            <CardContent sx={{ p: 0 }}>
              <List>
                {recentEvents.map((event) => (
                  <React.Fragment key={event.id}>
                    <ListItem>
                      <ListItemAvatar>
                        <Avatar sx={{ bgcolor: '#7289da' }}>
                          {event.type === 'command' && <MicIcon />}
                          {event.type === 'join' && <PeopleIcon />}
                          {event.type === 'message' && <ForumIcon />}
                          {event.type === 'voice' && <MicIcon />}
                        </Avatar>
                      </ListItemAvatar>
                      <ListItemText
                        primary={event.user}
                        secondary={
                          <React.Fragment>
                            <Typography component="span" variant="body2" color="text.primary">
                              {event.content}
                            </Typography>
                            {` — ${event.time}`}
                          </React.Fragment>
                        }
                      />
                    </ListItem>
                    <Divider variant="inset" component="li" sx={{ backgroundColor: '#40444b' }} />
                  </React.Fragment>
                ))}
              </List>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
};

export default Dashboard;