import { createTheme } from '@mui/material/styles';

export const theme = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: '#B19CD9', // пастельный фиолетовый
      light: '#D4C4E7',
      dark: '#8A7BB5',
    },
    secondary: {
      main: '#E6B8AF', // пастельный розовый
      light: '#F5D5D0',
      dark: '#C99A92',
    },
    background: {
      default: '#1A1A2E', // темно-синий фон
      paper: '#242442', // чуть светлее для карточек
    },
    text: {
      primary: '#E0E0E0',
      secondary: '#B8B8B8',
    },
  },
  typography: {
    fontFamily: '"Inter", "Roboto", "Helvetica", "Arial", sans-serif',
    h4: {
      fontWeight: 600,
      letterSpacing: '-0.5px',
    },
    h5: {
      fontWeight: 500,
      letterSpacing: '-0.25px',
    },
    h6: {
      fontWeight: 500,
    },
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          textTransform: 'none',
          fontWeight: 500,
        },
        contained: {
          boxShadow: 'none',
          '&:hover': {
            boxShadow: '0px 2px 4px rgba(0, 0, 0, 0.2)',
          },
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          borderRadius: 12,
          background: 'rgba(36, 36, 66, 0.8)',
          backdropFilter: 'blur(10px)',
          transition: 'transform 0.2s ease-in-out, box-shadow 0.2s ease-in-out',
          '&:hover': {
            transform: 'translateY(-4px)',
            boxShadow: '0 8px 16px rgba(0, 0, 0, 0.2)',
          },
        },
      },
    },
    MuiPaper: {
      styleOverrides: {
        root: {
          backgroundImage: 'none',
        },
      },
    },
    MuiCssBaseline: {
      styleOverrides: {
        body: {
          scrollbarWidth: 'thin',
          '&::-webkit-scrollbar': {
            width: '8px',
            height: '8px',
          },
          '&::-webkit-scrollbar-track': {
            background: 'rgba(0, 0, 0, 0.1)',
          },
          '&::-webkit-scrollbar-thumb': {
            background: 'rgba(177, 156, 217, 0.5)',
            borderRadius: '4px',
            '&:hover': {
              background: 'rgba(177, 156, 217, 0.7)',
            },
          },
        },
      },
    },
  },
  shape: {
    borderRadius: 8,
  },
}); 