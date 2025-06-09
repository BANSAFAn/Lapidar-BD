import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Box, Button, TextField, Typography, Paper, Container, Alert, CircularProgress } from '@mui/material';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import SecurityIcon from '@mui/icons-material/Security';
import apiService from '../services/api';

const Login = () => {
  const navigate = useNavigate();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [code, setCode] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [step, setStep] = useState(1); // 1 - email/password, 2 - 2FA
  const [tempToken, setTempToken] = useState('');

  const handleEmailPasswordSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const data = await apiService.login(email, password);

      if (data.success) {
        if (data.require_2fa) {
          // Переходим к вводу кода 2FA
          setTempToken(data.token);
          setStep(2);
        } else {
          // Если 2FA не требуется (что маловероятно в нашей реализации)
          localStorage.setItem('token', data.token);
          navigate('/dashboard');
        }
      } else {
        setError(data.message || 'Ошибка входа');
      }
    } catch (err) {
      setError('Ошибка сервера: ' + (err.response?.data?.message || err.message));
    } finally {
      setLoading(false);
    }
  };

  const handleTOTPSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const data = await apiService.verifyTOTP(tempToken, code);

      if (data.success) {
        localStorage.setItem('token', data.token);
        navigate('/dashboard');
      } else {
        setError(data.message || 'Неверный код');
      }
    } catch (err) {
      setError('Ошибка сервера: ' + (err.response?.data?.message || err.message));
    } finally {
      setLoading(false);
    }
  };

  return (
    <Container component="main" maxWidth="xs">
      <Paper
        elevation={6}
        sx={{
          marginTop: 8,
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          padding: 4,
        }}
      >
        {step === 1 ? (
          <>
            <LockOutlinedIcon sx={{ fontSize: 40, color: 'primary.main', mb: 2 }} />
            <Typography component="h1" variant="h5">
              Вход в панель управления
            </Typography>
            <Box component="form" onSubmit={handleEmailPasswordSubmit} sx={{ mt: 3 }}>
              {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
              <TextField
                margin="normal"
                required
                fullWidth
                id="email"
                label="Email"
                name="email"
                autoComplete="email"
                autoFocus
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
              <TextField
                margin="normal"
                required
                fullWidth
                name="password"
                label="Пароль"
                type="password"
                id="password"
                autoComplete="current-password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
              <Button
                type="submit"
                fullWidth
                variant="contained"
                sx={{ mt: 3, mb: 2 }}
                disabled={loading}
              >
                {loading ? <CircularProgress size={24} /> : 'Войти'}
              </Button>
            </Box>
          </>
        ) : (
          <>
            <SecurityIcon sx={{ fontSize: 40, color: 'primary.main', mb: 2 }} />
            <Typography component="h1" variant="h5">
              Двухфакторная аутентификация
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mt: 1, mb: 2, textAlign: 'center' }}>
              Введите код из приложения аутентификации
            </Typography>
            <Box component="form" onSubmit={handleTOTPSubmit} sx={{ mt: 1 }}>
              {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
              <TextField
                margin="normal"
                required
                fullWidth
                id="code"
                label="Код подтверждения"
                name="code"
                autoFocus
                value={code}
                onChange={(e) => setCode(e.target.value)}
                inputProps={{ maxLength: 6 }}
              />
              <Button
                type="submit"
                fullWidth
                variant="contained"
                sx={{ mt: 3, mb: 2 }}
                disabled={loading}
              >
                {loading ? <CircularProgress size={24} /> : 'Подтвердить'}
              </Button>
              <Button
                fullWidth
                variant="text"
                onClick={() => setStep(1)}
                sx={{ mt: 1 }}
              >
                Назад
              </Button>
            </Box>
          </>
        )}
      </Paper>
    </Container>
  );
};

export default Login;