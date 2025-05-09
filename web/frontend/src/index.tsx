import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import CssBaseline from '@mui/material/CssBaseline';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import App from './App';

// Создаем темную тему для приложения
const darkTheme = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: '#7289da', // Discord цвет
    },
    secondary: {
      main: '#43b581', // Discord зеленый
    },
    background: {
      default: '#36393f', // Discord фон
      paper: '#2f3136',  // Discord фон элементов
    },
  },
  typography: {
    fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
  },
});

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);

root.render(
  <React.StrictMode>
    <BrowserRouter>
      <ThemeProvider theme={darkTheme}>
        <CssBaseline />
        <App />
      </ThemeProvider>
    </BrowserRouter>
  </React.StrictMode>
);