import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import {
  Box,
  Drawer,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Divider,
  Typography,
} from '@mui/material';
import {
  Dashboard as DashboardIcon,
  Settings as SettingsIcon,
  Message as MessageIcon,
  Group as GroupIcon,
  VoiceChat as VoiceChatIcon,
  Code as CodeIcon,
} from '@mui/icons-material';

const drawerWidth = 240;

const Sidebar: React.FC = () => {
  const location = useLocation();

  const menuItems = [
    { text: 'Панель управления', icon: <DashboardIcon />, path: '/' },
    { text: 'Настройки', icon: <SettingsIcon />, path: '/settings' },
    { text: 'Команды', icon: <CodeIcon />, path: '/commands' },
    { text: 'Сообщения', icon: <MessageIcon />, path: '/messages' },
    { text: 'Пользователи', icon: <GroupIcon />, path: '/users' },
    { text: 'Голосовые каналы', icon: <VoiceChatIcon />, path: '/voice' },
  ];

  return (
    <Drawer
      variant="permanent"
      sx={{
        width: drawerWidth,
        flexShrink: 0,
        '& .MuiDrawer-paper': {
          width: drawerWidth,
          boxSizing: 'border-box',
          backgroundColor: '#2f3136', // Discord sidebar color
          color: '#ffffff',
        },
      }}
    >
      <Box sx={{ p: 2, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        <Typography variant="h6" component="div" sx={{ fontWeight: 'bold' }}>
          Lapidar BD
        </Typography>
      </Box>
      <Divider sx={{ backgroundColor: '#40444b' }} />
      <List>
        {menuItems.map((item) => (
          <ListItem key={item.text} disablePadding>
            <ListItemButton
              component={Link}
              to={item.path}
              selected={location.pathname === item.path}
              sx={{
                '&.Mui-selected': {
                  backgroundColor: '#393c43',
                  '&:hover': {
                    backgroundColor: '#42464D',
                  },
                },
                '&:hover': {
                  backgroundColor: '#42464D',
                },
              }}
            >
              <ListItemIcon sx={{ color: '#b9bbbe' }}>{item.icon}</ListItemIcon>
              <ListItemText primary={item.text} />
            </ListItemButton>
          </ListItem>
        ))}
      </List>
      <Box sx={{ flexGrow: 1 }} />
      <Box sx={{ p: 2, textAlign: 'center' }}>
        <Typography variant="caption" color="text.secondary">
          Версия 1.0.0
        </Typography>
      </Box>
    </Drawer>
  );
};

export default Sidebar;