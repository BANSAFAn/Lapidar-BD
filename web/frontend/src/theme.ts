import { createTheme } from '@mui/material/styles';
import { alpha } from '@mui/material/styles';

// Создаем улучшенную темную тему в стиле Discord для приложения
const theme = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: '#7289da', // Discord цвет
      light: '#8ea1e1',
      dark: '#5b6eae',
      contrastText: '#ffffff',
    },
    secondary: {
      main: '#43b581', // Discord зеленый
      light: '#4fd592',
      dark: '#3a9c6f',
      contrastText: '#ffffff',
    },
    error: {
      main: '#f04747', // Discord красный
      light: '#f36d6d',
      dark: '#d03c3c',
      contrastText: '#ffffff',
    },
    warning: {
      main: '#faa61a', // Discord оранжевый
      light: '#fbb746',
      dark: '#d48c16',
      contrastText: '#ffffff',
    },
    info: {
      main: '#7289da', // Discord синий
      light: '#8ea1e1',
      dark: '#5b6eae',
      contrastText: '#ffffff',
    },
    success: {
      main: '#43b581', // Discord зеленый
      light: '#4fd592',
      dark: '#3a9c6f',
      contrastText: '#ffffff',
    },
    background: {
      default: '#36393f', // Discord фон
      paper: '#2f3136',  // Discord фон элементов
    },
    text: {
      primary: '#ffffff',
      secondary: '#b9bbbe',
      disabled: '#72767d',
    },
    divider: '#40444b',
    action: {
      active: '#ffffff',
      hover: 'rgba(255, 255, 255, 0.08)',
      selected: 'rgba(114, 137, 218, 0.16)',
      disabled: 'rgba(255, 255, 255, 0.3)',
      disabledBackground: 'rgba(255, 255, 255, 0.12)',
    },
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
    MuiCssBaseline: {
      styleOverrides: {
        body: {
          scrollbarWidth: 'thin',
          scrollbarColor: '#202225 #2e3338',
          '&::-webkit-scrollbar': {
            width: '8px',
            height: '8px',
          },
          '&::-webkit-scrollbar-track': {
            background: '#2e3338',
            borderRadius: '10px',
          },
          '&::-webkit-scrollbar-thumb': {
            backgroundColor: '#202225',
            borderRadius: '10px',
            '&:hover': {
              backgroundColor: '#7289da',
            },
          },
        },
      },
    },
    MuiAppBar: {
      styleOverrides: {
        root: {
          backgroundColor: '#36393f',
          boxShadow: 'none',
          borderBottom: '1px solid #40444b',
          backdropFilter: 'blur(10px)',
          zIndex: 1100,
        },
      },
    },
    MuiDrawer: {
      styleOverrides: {
        paper: {
          backgroundColor: '#2f3136',
          borderRight: 'none',
          boxShadow: '0 0 10px rgba(0, 0, 0, 0.2)',
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          backgroundColor: '#2f3136',
          borderRadius: '8px',
          boxShadow: '0 2px 10px rgba(0, 0, 0, 0.2)',
          transition: 'transform 0.2s, box-shadow 0.2s',
          '&:hover': {
            transform: 'translateY(-2px)',
            boxShadow: '0 5px 15px rgba(0, 0, 0, 0.3)',
          },
          overflow: 'hidden',
        },
      },
    },
    MuiCardHeader: {
      styleOverrides: {
        root: {
          padding: '16px 20px',
          borderBottom: '1px solid #40444b',
        },
        title: {
          fontWeight: 600,
          fontSize: '1.1rem',
        },
        subheader: {
          color: '#b9bbbe',
        },
        action: {
          marginTop: 0,
          marginRight: 0,
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
    MuiCardActions: {
      styleOverrides: {
        root: {
          padding: '12px 20px',
          borderTop: '1px solid #40444b',
        },
      },
    },
    MuiButton: {
      styleOverrides: {
        root: {
          textTransform: 'none',
          fontWeight: 600,
          borderRadius: '4px',
          padding: '8px 16px',
          transition: 'all 0.2s',
          '&:hover': {
            transform: 'translateY(-1px)',
            boxShadow: '0 4px 8px rgba(0, 0, 0, 0.2)',
          },
        },
        contained: {
          boxShadow: '0 2px 5px rgba(0, 0, 0, 0.2)',
        },
        containedPrimary: {
          '&:hover': {
            backgroundColor: '#8ea1e1',
          },
        },
        containedSecondary: {
          '&:hover': {
            backgroundColor: '#4fd592',
          },
        },
        outlined: {
          borderWidth: '1.5px',
          '&:hover': {
            borderWidth: '1.5px',
          },
        },
      },
    },
    MuiIconButton: {
      styleOverrides: {
        root: {
          transition: 'all 0.2s',
          '&:hover': {
            backgroundColor: 'rgba(114, 137, 218, 0.1)',
            transform: 'scale(1.1)',
          },
        },
      },
    },
    MuiListItem: {
      styleOverrides: {
        root: {
          borderRadius: '4px',
          margin: '2px 0',
          '&.Mui-selected': {
            backgroundColor: 'rgba(114, 137, 218, 0.1)',
            borderLeft: '3px solid #7289da',
            paddingLeft: '13px',
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
          transition: 'all 0.2s',
          '&.Mui-selected': {
            backgroundColor: 'rgba(114, 137, 218, 0.1)',
            borderLeft: '3px solid #7289da',
            paddingLeft: '13px',
          },
          '&.Mui-selected:hover': {
            backgroundColor: 'rgba(114, 137, 218, 0.2)',
          },
          '&:hover': {
            backgroundColor: 'rgba(114, 137, 218, 0.08)',
            transform: 'translateX(2px)',
          },
        },
      },
    },
    MuiListItemIcon: {
      styleOverrides: {
        root: {
          color: '#b9bbbe',
          minWidth: '40px',
        },
      },
    },
    MuiListItemText: {
      styleOverrides: {
        primary: {
          fontWeight: 500,
        },
        secondary: {
          color: '#b9bbbe',
        },
      },
    },
    MuiTextField: {
      styleOverrides: {
        root: {
          marginBottom: '16px',
          '& .MuiOutlinedInput-root': {
            borderRadius: '4px',
            transition: 'all 0.2s',
            '& fieldset': {
              borderColor: '#40444b',
              borderWidth: '1.5px',
            },
            '&:hover fieldset': {
              borderColor: '#7289da',
            },
            '&.Mui-focused fieldset': {
              borderColor: '#7289da',
              borderWidth: '2px',
            },
            '&.Mui-error fieldset': {
              borderColor: '#f04747',
            },
          },
          '& .MuiInputLabel-root': {
            '&.Mui-focused': {
              color: '#7289da',
            },
          },
        },
      },
    },
    MuiChip: {
      styleOverrides: {
        root: {
          borderRadius: '4px',
          fontWeight: 500,
          '&.MuiChip-colorPrimary': {
            backgroundColor: alpha('#7289da', 0.2),
            color: '#7289da',
          },
          '&.MuiChip-colorSecondary': {
            backgroundColor: alpha('#43b581', 0.2),
            color: '#43b581',
          },
          '&.MuiChip-colorError': {
            backgroundColor: alpha('#f04747', 0.2),
            color: '#f04747',
          },
          '&.MuiChip-colorWarning': {
            backgroundColor: alpha('#faa61a', 0.2),
            color: '#faa61a',
          },
        },
        deleteIcon: {
          color: 'inherit',
          opacity: 0.7,
          '&:hover': {
            opacity: 1,
            color: 'inherit',
          },
        },
      },
    },
    MuiTooltip: {
      styleOverrides: {
        tooltip: {
          backgroundColor: '#18191c',
          color: '#ffffff',
          border: '1px solid #2f3136',
          borderRadius: '4px',
          padding: '8px 12px',
          fontSize: '0.85rem',
          boxShadow: '0 2px 10px rgba(0, 0, 0, 0.2)',
        },
        arrow: {
          color: '#18191c',
        },
      },
    },
    MuiDivider: {
      styleOverrides: {
        root: {
          borderColor: '#40444b',
        },
      },
    },
    MuiSwitch: {
      styleOverrides: {
        root: {
          width: 42,
          height: 26,
          padding: 0,
          margin: 8,
        },
        switchBase: {
          padding: 1,
          '&.Mui-checked': {
            transform: 'translateX(16px)',
            color: '#fff',
            '& + .MuiSwitch-track': {
              backgroundColor: '#43b581',
              opacity: 1,
              border: 'none',
            },
          },
          '&.Mui-focusVisible .MuiSwitch-thumb': {
            color: '#43b581',
            border: '6px solid #fff',
          },
        },
        thumb: {
          width: 24,
          height: 24,
        },
        track: {
          borderRadius: 26 / 2,
          backgroundColor: '#72767d',
          opacity: 1,
        },
      },
    },
    MuiPaper: {
      styleOverrides: {
        root: {
          backgroundImage: 'none',
        },
        elevation1: {
          boxShadow: '0 2px 10px rgba(0, 0, 0, 0.2)',
        },
      },
    },
  },
});

export default theme;