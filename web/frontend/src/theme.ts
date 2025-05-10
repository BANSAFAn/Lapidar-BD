import { createTheme } from '@mui/material/styles';

// Создаем темную тему в стиле Discord для приложения
const theme = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: '#7289da', // Discord цвет
      light: '#8ea1e1',
      dark: '#5b6eae',
    },
    secondary: {
      main: '#43b581', // Discord зеленый
      light: '#4fd592',
      dark: '#3a9c6f',
    },
    error: {
      main: '#f04747', // Discord красный
    },
    warning: {
      main: '#faa61a', // Discord оранжевый
    },
    info: {
      main: '#7289da', // Discord синий
    },
    success: {
      main: '#43b581', // Discord зеленый
    },
    background: {
      default: '#36393f', // Discord фон
      paper: '#2f3136',  // Discord фон элементов
    },
    text: {
      primary: '#ffffff',
      secondary: '#b9bbbe',
    },
    divider: '#40444b',
  },
  typography: {
    fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
    h4: {
      fontWeight: 600,
      fontSize: '1.75rem',
    },
    h5: {
      fontWeight: 600,
      fontSize: '1.25rem',
    },
    h6: {
      fontWeight: 600,
      fontSize: '1rem',
    },
    subtitle1: {
      fontWeight: 600,
    },
  },
  components: {
    MuiAppBar: {
      styleOverrides: {
        root: {
          backgroundColor: '#36393f',
          boxShadow: 'none',
          borderBottom: '1px solid #40444b',
        },
      },
    },
    MuiDrawer: {
      styleOverrides: {
        paper: {
          backgroundColor: '#2f3136',
          borderRight: 'none',
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          backgroundColor: '#2f3136',
          borderRadius: '8px',
        },
      },
    },
    MuiCardHeader: {
      styleOverrides: {
        root: {
          padding: '16px 20px',
        },
        title: {
          fontWeight: 600,
        },
      },
    },
    MuiCardContent: {
      styleOverrides: {
        root: {
          padding: '16px 20px',
          '&:last-child': {
            paddingBottom: '20px',
          },
        },
      },
    },
    MuiButton: {
      styleOverrides: {
        root: {
          textTransform: 'none',
          fontWeight: 600,
          borderRadius: '4px',
        },
      },
    },
    MuiListItem: {
      styleOverrides: {
        root: {
          borderRadius: '4px',
          '&.Mui-selected': {
            backgroundColor: 'rgba(114, 137, 218, 0.1)',
          },
          '&.Mui-selected:hover': {
            backgroundColor: 'rgba(114, 137, 218, 0.2)',
          },
        },
      },
    },
    MuiListItemButton: {
      styleOverrides: {
        root: {
          borderRadius: '4px',
          '&.Mui-selected': {
            backgroundColor: 'rgba(114, 137, 218, 0.1)',
          },
          '&.Mui-selected:hover': {
            backgroundColor: 'rgba(114, 137, 218, 0.2)',
          },
        },
      },
    },
    MuiTextField: {
      styleOverrides: {
        root: {
          '& .MuiOutlinedInput-root': {
            '& fieldset': {
              borderColor: '#40444b',
            },
            '&:hover fieldset': {
              borderColor: '#7289da',
            },
          },
        },
      },
    },
    MuiChip: {
      styleOverrides: {
        root: {
          borderRadius: '4px',
        },
      },
    },
  },
});

export default theme;