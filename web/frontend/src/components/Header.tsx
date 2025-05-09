import React from 'react';
import {
  AppBar,
  Toolbar,
  Typography,
  IconButton,
  Badge,
  Box,
  Chip,
} from '@mui/material';
import {
  Notifications as NotificationsIcon,
  CheckCircle as CheckCircleIcon,
  Refresh as RefreshIcon,
} from '@mui/icons-material';

const Header: React.FC = () => {
  // В реальном приложении здесь будет состояние, получаемое от API
  const botStatus = 'online';

  return (
    <AppBar position="sticky" sx={{ backgroundColor: '#36393f', boxShadow: 'none', borderBottom: '1px solid #40444b' }}>
      <Toolbar>
        <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
          Панель управления
        </Typography>
        
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
          <Chip
            icon={<CheckCircleIcon />}
            label={botStatus === 'online' ? 'Онлайн' : 'Оффлайн'}
            color={botStatus === 'online' ? 'success' : 'error'}
            size="small"
            variant="outlined"
          />
          
          <IconButton color="inherit" size="small">
            <RefreshIcon />
          </IconButton>
          
          <IconButton color="inherit" size="small">
            <Badge badgeContent={4} color="error">
              <NotificationsIcon />
            </Badge>
          </IconButton>
        </Box>
      </Toolbar>
    </AppBar>
  );
};

export default Header;