import React from 'react';
import { Routes, Route } from 'react-router-dom';
import { Box, Container } from '@mui/material';

// Компоненты
import Dashboard from './pages/Dashboard';
import Settings from './pages/Settings';
import Commands from './pages/Commands';
import Sidebar from './components/Sidebar';
import Header from './components/Header';

const App: React.FC = () => {
  return (
    <Box sx={{ display: 'flex', minHeight: '100vh' }}>
      <Sidebar />
      <Box sx={{ flexGrow: 1, display: 'flex', flexDirection: 'column' }}>
        <Header />
        <Container component="main" sx={{ flexGrow: 1, p: 3, mt: 2 }}>
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/settings" element={<Settings />} />
            <Route path="/commands" element={<Commands />} />
          </Routes>
        </Container>
      </Box>
    </Box>
  );
};

export default App;