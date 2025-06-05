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
  Chip,
  IconButton,
  Paper,
  Stack,
  Tooltip,
} from '@mui/material';
import { 
  Save as SaveIcon,
  Add as AddIcon,
  Delete as DeleteIcon,
  Info as InfoIcon,
} from '@mui/icons-material';
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
    AltPorts: number[];
  };
}

const Settings: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [newPort, setNewPort] = useState<string>('');
  const [config, setConfig] = useState<BotConfig>({
    Token: '',
    Prefix: '!',
    BotName: 'Lapidar Bot',
    DefaultLanguage: 'ru',
    WebInterface: {
      Enabled: true,
      Host: 'localhost',
      Port: 8080,
      AltPorts: [3000, 8000],
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
              AltPorts: [3000, 8000],
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

  // Обработчик добавления альтернативного порта
  const handleAddAltPort = () => {
    if (!newPort || isNaN(Number(newPort))) return;
    
    const portNumber = parseInt(newPort, 10);
    if (portNumber <= 0 || portNumber > 65535) {
      setError('Порт должен быть в диапазоне от 1 до 65535');
      return;
    }
    
    // Проверяем, не существует ли уже такой порт
    if (config.WebInterface.AltPorts.includes(portNumber)) {
      setError('Этот порт уже добавлен в список альтернативных');
      return;
    }
    
    // Проверяем, не совпадает ли с основным портом
    if (portNumber === Number(config.WebInterface.Port)) {
      setError('Этот порт уже используется как основной');
      return;
    }
    
    const updatedAltPorts = [...config.WebInterface.AltPorts, portNumber];
    setConfig({
      ...config,
      WebInterface: {
        ...config.WebInterface,
        AltPorts: updatedAltPorts,
      },
    });
    setNewPort('');
  };

  // Обработчик удаления альтернативного порта
  const handleRemoveAltPort = (port: number) => {
    const updatedAltPorts = config.WebInterface.AltPorts.filter(p => p !== port);
    setConfig({
      ...config,
      WebInterface: {
        ...config.WebInterface,
        AltPorts: updatedAltPorts,
      },
    });
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
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '70vh' }}>
        <CircularProgress size={40} thickness={4} />
        <Typography variant="h6" sx={{ ml: 2 }}>
          Загрузка настроек...
        </Typography>
      </Box>
    );
  }

  return (
    <Box>
      <Box sx={{ 
        display: 'flex', 
        justifyContent: 'space-between', 
        alignItems: 'center', 
        mb: 3,
        borderBottom: '1px solid #40444b',
        pb: 2
      }}>
        <Typography variant="h4" sx={{ fontWeight: 600, color: '#fff' }}>
          Настройки бота
        </Typography>
        <Button
          variant="contained"
          color="primary"
          startIcon={<SaveIcon />}
          onClick={handleSubmit}
          disabled={saving}
          sx={{ 
            minWidth: 150,
            py: 1,
            boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)'
          }}
        >
          {saving ? <CircularProgress size={24} color="inherit" /> : 'Сохранить'}
        </Button>
      </Box>
      
      <form onSubmit={handleSubmit}>
        <Grid container spacing={3}>
          {/* Основные настройки */}
          <Grid item xs={12} md={6}>
            <Card sx={{ 
              backgroundColor: '#2f3136',
              boxShadow: '0 4px 10px rgba(0, 0, 0, 0.15)',
              borderRadius: '10px',
              overflow: 'hidden',
              height: '100%',
              display: 'flex',
              flexDirection: 'column'
            }}>
              <CardHeader 
                title="Основные настройки" 
                titleTypographyProps={{ fontWeight: 600 }}
                sx={{ 
                  backgroundColor: 'rgba(114, 137, 218, 0.1)',
                  borderBottom: '1px solid #40444b',
                  py: 2
                }}
              />
              <Divider sx={{ backgroundColor: '#40444b' }} />
              <CardContent sx={{ flexGrow: 1 }}>
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
            <Card sx={{ 
              backgroundColor: '#2f3136',
              boxShadow: '0 4px 10px rgba(0, 0, 0, 0.15)',
              borderRadius: '10px',
              overflow: 'hidden',
              height: '100%',
              display: 'flex',
              flexDirection: 'column'
            }}>
              <CardHeader 
                title="Настройки веб-интерфейса" 
                titleTypographyProps={{ fontWeight: 600 }}
                sx={{ 
                  backgroundColor: 'rgba(114, 137, 218, 0.1)',
                  borderBottom: '1px solid #40444b',
                  py: 2
                }}
              />
              <Divider sx={{ backgroundColor: '#40444b' }} />
              <CardContent sx={{ flexGrow: 1 }}>
                <Grid container spacing={2}>
                  <Grid item xs={12}>
                    <Paper 
                      variant="outlined" 
                      sx={{ 
                        p: 2, 
                        backgroundColor: 'rgba(47, 49, 54, 0.6)', 
                        borderRadius: '8px',
                        mb: 2,
                        display: 'flex',
                        justifyContent: 'space-between',
                        alignItems: 'center'
                      }}
                    >
                      <Box>
                        <Typography variant="subtitle1" sx={{ fontWeight: 500 }}>
                          Веб-интерфейс
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          Включить или отключить веб-интерфейс администратора
                        </Typography>
                      </Box>
                      <Switch
                        checked={config.WebInterface.Enabled}
                        onChange={handleChange}
                        name="WebInterface.Enabled"
                        color="primary"
                        sx={{ ml: 2 }}
                      />
                    </Paper>
                  </Grid>
                  <Grid item xs={12}>
                    <Typography variant="subtitle1" sx={{ mt: 1, mb: 1, fontWeight: 500 }}>
                      Настройки подключения
                    </Typography>
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
                      InputProps={{
                        sx: { borderRadius: '8px' }
                      }}
                      sx={{
                        '& .MuiOutlinedInput-root': {
                          '&.Mui-disabled': {
                            backgroundColor: 'rgba(47, 49, 54, 0.4)',
                          }
                        }
                      }}
                    />
                  </Grid>
                  <Grid item xs={12} sm={6}>
                    <TextField
                      fullWidth
                      label="Основной порт"
                      name="WebInterface.Port"
                      value={config.WebInterface.Port}
                      onChange={handleChange}
                      margin="normal"
                      type="number"
                      disabled={!config.WebInterface.Enabled}
                      helperText="Основной порт для веб-интерфейса"
                      InputProps={{
                        sx: { borderRadius: '8px' }
                      }}
                      sx={{
                        '& .MuiOutlinedInput-root': {
                          '&.Mui-disabled': {
                            backgroundColor: 'rgba(47, 49, 54, 0.4)',
                          }
                        }
                      }}
                    />
                  </Grid>
                  
                  {/* Альтернативные порты */}
                  <Grid item xs={12}>
                    <Typography variant="subtitle1" sx={{ mt: 3, mb: 1, fontWeight: 500, display: 'flex', alignItems: 'center' }}>
                      Альтернативные порты
                      <Tooltip 
                        title="Дополнительные порты, на которых будет доступен веб-интерфейс. Полезно, если основной порт заблокирован." 
                        arrow
                        placement="top"
                      >
                        <IconButton size="small" sx={{ ml: 0.5, color: 'text.secondary' }}>
                          <InfoIcon fontSize="small" />
                        </IconButton>
                      </Tooltip>
                    </Typography>
                    
                    <Paper 
                      variant="outlined" 
                      sx={{ 
                        p: 2, 
                        backgroundColor: 'rgba(47, 49, 54, 0.6)', 
                        borderRadius: '8px',
                        mb: 2
                      }}
                    >
                      <Box sx={{ display: 'flex', alignItems: 'flex-start' }}>
                        <TextField
                          label="Новый порт"
                          value={newPort}
                          onChange={(e) => setNewPort(e.target.value)}
                          type="number"
                          size="small"
                          disabled={!config.WebInterface.Enabled}
                          sx={{ 
                            mr: 1, 
                            width: '150px',
                            '& .MuiOutlinedInput-root': {
                              borderRadius: '8px',
                              '&.Mui-disabled': {
                                backgroundColor: 'rgba(47, 49, 54, 0.4)',
                              }
                            }
                          }}
                          placeholder="Введите порт"
                        />
                        <Button
                          variant="contained"
                          color="primary"
                          onClick={handleAddAltPort}
                          disabled={!config.WebInterface.Enabled || !newPort}
                          startIcon={<AddIcon />}
                          size="small"
                          sx={{ 
                            borderRadius: '8px',
                            textTransform: 'none',
                            fontWeight: 500,
                            boxShadow: '0 2px 5px rgba(0, 0, 0, 0.2)',
                            '&:hover': {
                              boxShadow: '0 4px 8px rgba(0, 0, 0, 0.3)'
                            },
                            '&.Mui-disabled': {
                              backgroundColor: 'rgba(114, 137, 218, 0.3)'
                            }
                          }}
                        >
                          Добавить
                        </Button>
                      </Box>
                    </Paper>
                    
                    <Paper 
                      variant="outlined" 
                      sx={{ 
                        p: 2, 
                        backgroundColor: 'rgba(47, 49, 54, 0.6)', 
                        borderRadius: '8px',
                        minHeight: '60px',
                        display: 'flex',
                        flexWrap: 'wrap',
                        gap: 1,
                        alignItems: 'center'
                      }}
                    >
                      {config.WebInterface.AltPorts.length > 0 ? (
                        config.WebInterface.AltPorts.map((port) => (
                          <Chip
                            key={port}
                            label={port}
                            color="primary"
                            disabled={!config.WebInterface.Enabled}
                            onDelete={() => handleRemoveAltPort(port)}
                            deleteIcon={<DeleteIcon fontSize="small" />}
                            sx={{ 
                              m: 0.5, 
                              borderRadius: '6px',
                              '&:hover': {
                                backgroundColor: 'rgba(114, 137, 218, 0.25)'
                              },
                              '&.Mui-disabled': {
                                opacity: 0.6,
                                backgroundColor: 'rgba(47, 49, 54, 0.4)'
                              }
                            }}
                          />
                        ))
                      ) : (
                        <Typography variant="body2" color="text.secondary" sx={{ fontStyle: 'italic' }}>
                          Нет альтернативных портов
                        </Typography>
                      )}
                    </Paper>
                  </Grid>
                </Grid>
              </CardContent>
            </Card>
          </Grid>
          
          {/* Конец формы */}
        </Grid>
      </form>
      
      {/* Уведомления */}
      <Snackbar 
        open={success} 
        autoHideDuration={4000} 
        onClose={handleCloseAlert}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
      >
        <Alert 
          onClose={handleCloseAlert} 
          severity="success" 
          variant="filled"
          elevation={6}
          sx={{ 
            width: '100%',
            fontWeight: 500,
            '& .MuiAlert-icon': { fontSize: '1.2rem' }
          }}
        >
          Настройки успешно сохранены!
        </Alert>
      </Snackbar>
      
      <Snackbar 
        open={!!error} 
        autoHideDuration={6000} 
        onClose={handleCloseAlert}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
      >
        <Alert 
          onClose={handleCloseAlert} 
          severity="error" 
          variant="filled"
          elevation={6}
          sx={{ 
            width: '100%',
            fontWeight: 500,
            '& .MuiAlert-icon': { fontSize: '1.2rem' }
          }}
        >
          {error}
        </Alert>
      </Snackbar>
    </Box>
  );
};

export default Settings;