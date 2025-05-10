import React, { useState, useEffect } from 'react';
import {
  Box,
  Card,
  CardContent,
  CardHeader,
  Divider,
  TextField,
  Button,
  Grid,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  FormControlLabel,
  Switch,
  Alert,
  Snackbar,
  Typography,
  CircularProgress,
} from '@mui/material';
import { Save as SaveIcon } from '@mui/icons-material';
import axios from 'axios';

// Интерфейс для конфигурации бота
interface BotConfig {
  Token: string;
  Prefix: string;
  BotName: string;
  DefaultLanguage: string;
  WebInterface: {
    Enabled: boolean;
    Host: string;
    Port: number;
  };
}

const Settings: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [config, setConfig] = useState<BotConfig>({
    Token: '',
    Prefix: '!',
    BotName: 'Lapidar Bot',
    DefaultLanguage: 'ru',
    WebInterface: {
      Enabled: true,
      Host: 'localhost',
      Port: 8080,
    },
  });

  // Загрузка конфигурации при монтировании компонента
  useEffect(() => {
    const fetchConfig = async () => {
      try {
        // В реальном приложении здесь будет запрос к API
        // const response = await axios.get('/api/config');
        // setConfig(response.data);
        
        // Имитация загрузки данных
        setTimeout(() => {
          setConfig({
            Token: 'YOUR_DISCORD_TOKEN',
            Prefix: '!',
            BotName: 'Lapidar Bot',
            DefaultLanguage: 'ru',
            WebInterface: {
              Enabled: true,
              Host: 'localhost',
              Port: 8080,
            },
          });
          setLoading(false);
        }, 1000);
      } catch (err) {
        setError('Ошибка загрузки конфигурации');
        setLoading(false);
      }
    };

    fetchConfig();
  }, []);

  // Обработчик изменения полей формы
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value, checked, type } = e.target;
    
    if (name.includes('.')) {
      // Обработка вложенных полей (например, WebInterface.Enabled)
      const [parent, child] = name.split('.');
      setConfig({
        ...config,
        [parent]: {
          ...config[parent as keyof BotConfig],
          [child]: type === 'checkbox' ? checked : value,
        },
      });
    } else {
      // Обработка полей верхнего уровня
      setConfig({
        ...config,
        [name]: type === 'checkbox' ? checked : value,
      });
    }
  };

  // Обработчик отправки формы
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    setError(null);

    try {
      await apiService.saveConfig(config);
      setSaving(false);
      setSuccess(true);
      
      // Скрываем сообщение об успехе через 3 секунды
      setTimeout(() => setSuccess(false), 3000);
    } catch (err) {
      setSaving(false);
      setError('Ошибка сохранения конфигурации');
    }
  };

  // Обработчик закрытия уведомлений
  const handleCloseAlert = () => {
    setSuccess(false);
    setError(null);
  };

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
        Настройки бота
      </Typography>
      
      <form onSubmit={handleSubmit}>
        <Grid container spacing={3}>
          {/* Основные настройки */}
          <Grid item xs={12} md={6}>
            <Card sx={{ backgroundColor: '#2f3136' }}>
              <CardHeader title="Основные настройки" />
              <Divider sx={{ backgroundColor: '#40444b' }} />
              <CardContent>
                <Grid container spacing={2}>
                  <Grid item xs={12}>
                    <TextField
                      fullWidth
                      label="Discord Token"
                      name="Token"
                      value={config.Token}
                      onChange={handleChange}
                      margin="normal"
                      type="password"
                      required
                      helperText="Токен вашего Discord бота"
                    />
                  </Grid>
                  <Grid item xs={12} sm={6}>
                    <TextField
                      fullWidth
                      label="Префикс команд"
                      name="Prefix"
                      value={config.Prefix}
                      onChange={handleChange}
                      margin="normal"
                      required
                      helperText="Префикс для команд бота (например, !)"
                    />
                  </Grid>
                  <Grid item xs={12} sm={6}>
                    <TextField
                      fullWidth
                      label="Имя бота"
                      name="BotName"
                      value={config.BotName}
                      onChange={handleChange}
                      margin="normal"
                      required
                      helperText="Отображаемое имя бота"
                    />
                  </Grid>
                  <Grid item xs={12}>
                    <FormControl fullWidth margin="normal">
                      <InputLabel>Язык по умолчанию</InputLabel>
                      <Select
                        name="DefaultLanguage"
                        value={config.DefaultLanguage}
                        onChange={handleChange as any}
                        label="Язык по умолчанию"
                      >
                        <MenuItem value="ru">Русский</MenuItem>
                        <MenuItem value="en">English</MenuItem>
                        <MenuItem value="de">Deutsch</MenuItem>
                        <MenuItem value="uk">Українська</MenuItem>
                        <MenuItem value="zh">中文</MenuItem>
                      </Select>
                    </FormControl>
                  </Grid>
                </Grid>
              </CardContent>
            </Card>
          </Grid>
          
          {/* Настройки веб-интерфейса */}
          <Grid item xs={12} md={6}>
            <Card sx={{ backgroundColor: '#2f3136' }}>
              <CardHeader title="Настройки веб-интерфейса" />
              <Divider sx={{ backgroundColor: '#40444b' }} />
              <CardContent>
                <Grid container spacing={2}>
                  <Grid item xs={12}>
                    <FormControlLabel
                      control={
                        <Switch
                          checked={config.WebInterface.Enabled}
                          onChange={handleChange}
                          name="WebInterface.Enabled"
                          color="primary"
                        />
                      }
                      label="Включить веб-интерфейс"
                    />
                  </Grid>
                  <Grid item xs={12} sm={6}>
                    <TextField
                      fullWidth
                      label="Хост"
                      name="WebInterface.Host"
                      value={config.WebInterface.Host}
                      onChange={handleChange}
                      margin="normal"
                      disabled={!config.WebInterface.Enabled}
                      helperText="Хост для веб-интерфейса"
                    />
                  </Grid>
                  <Grid item xs={12} sm={6}>
                    <TextField
                      fullWidth
                      label="Порт"
                      name="WebInterface.Port"
                      value={config.WebInterface.Port}
                      onChange={handleChange}
                      margin="normal"
                      type="number"
                      disabled={!config.WebInterface.Enabled}
                      helperText="Порт для веб-интерфейса"
                    />
                  </Grid>
                </Grid>
              </CardContent>
            </Card>
          </Grid>
          
          {/* Кнопка сохранения */}
          <Grid item xs={12}>
            <Box sx={{ display: 'flex', justifyContent: 'flex-end', mt: 2 }}>
              <Button
                type="submit"
                variant="contained"
                color="primary"
                startIcon={<SaveIcon />}
                disabled={saving}
                sx={{ minWidth: 150 }}
              >
                {saving ? <CircularProgress size={24} /> : 'Сохранить'}
              </Button>
            </Box>
          </Grid>
        </Grid>
      </form>
      
      {/* Уведомления */}
      <Snackbar open={success} autoHideDuration={6000} onClose={handleCloseAlert}>
        <Alert onClose={handleCloseAlert} severity="success" sx={{ width: '100%' }}>
          Настройки успешно сохранены!
        </Alert>
      </Snackbar>
      
      <Snackbar open={!!error} autoHideDuration={6000} onClose={handleCloseAlert}>
        <Alert onClose={handleCloseAlert} severity="error" sx={{ width: '100%' }}>
          {error}
        </Alert>
      </Snackbar>
    </Box>
  );
};

export default Settings;