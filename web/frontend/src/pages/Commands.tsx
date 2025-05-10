import React, { useState, useEffect } from 'react';
import {
  Box,
  Card,
  CardContent,
  CardHeader,
  Divider,
  List,
  ListItem,
  ListItemText,
  ListItemSecondaryAction,
  Switch,
  Typography,
  TextField,
  InputAdornment,
  IconButton,
  Chip,
  Grid,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  CircularProgress,
  Alert,
  Snackbar,
} from '@mui/material';
import {
  Search as SearchIcon,
  FilterList as FilterListIcon,
} from '@mui/icons-material';
import apiService, { Command } from '../services/api';

const Commands: React.FC = () => {
  const [commands, setCommands] = useState<Command[]>([]);
  const [filteredCommands, setFilteredCommands] = useState<Command[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const [categoryFilter, setCategoryFilter] = useState('all');
  const [snackbar, setSnackbar] = useState({
    open: false,
    message: '',
    severity: 'success' as 'success' | 'error',
  });

  // Загрузка списка команд при монтировании компонента
  useEffect(() => {
    const fetchCommands = async () => {
      try {
        const data = await apiService.getCommands();
        setCommands(data);
        setFilteredCommands(data);
        setLoading(false);
      } catch (error) {
        console.error('Ошибка при загрузке команд:', error);
        setLoading(false);
        setSnackbar({
          open: true,
          message: 'Ошибка при загрузке списка команд',
          severity: 'error',
        });
      }
    };

    fetchCommands();
  }, []);

  // Фильтрация команд при изменении поискового запроса или категории
  useEffect(() => {
    let result = commands;

    // Фильтрация по поисковому запросу
    if (searchTerm) {
      result = result.filter(
        (command) =>
          command.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
          command.description.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }

    // Фильтрация по категории
    if (categoryFilter !== 'all') {
      result = result.filter((command) => command.category === categoryFilter);
    }

    setFilteredCommands(result);
  }, [searchTerm, categoryFilter, commands]);

  // Получение уникальных категорий для фильтра
  const categories = ['all', ...new Set(commands.map((cmd) => cmd.category))];

  // Обработчик изменения статуса команды
  const handleToggleCommand = async (command: Command) => {
    try {
      const updatedCommand = { ...command, enabled: !command.enabled };
      await apiService.updateCommand(updatedCommand);

      // Обновляем локальное состояние
      setCommands(
        commands.map((cmd) =>
          cmd.name === command.name ? updatedCommand : cmd
        )
      );

      setSnackbar({
        open: true,
        message: `Команда ${command.name} ${updatedCommand.enabled ? 'включена' : 'отключена'}`,
        severity: 'success',
      });
    } catch (error) {
      console.error('Ошибка при обновлении команды:', error);
      setSnackbar({
        open: true,
        message: 'Ошибка при обновлении статуса команды',
        severity: 'error',
      });
    }
  };

  // Обработчик закрытия уведомления
  const handleCloseSnackbar = () => {
    setSnackbar({ ...snackbar, open: false });
  };

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Управление командами
      </Typography>

      <Card sx={{ mb: 4 }}>
        <CardHeader title="Фильтры" />
        <Divider />
        <CardContent>
          <Grid container spacing={2}>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Поиск команд"
                variant="outlined"
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <SearchIcon />
                    </InputAdornment>
                  ),
                }}
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <FormControl fullWidth variant="outlined">
                <InputLabel id="category-filter-label">Категория</InputLabel>
                <Select
                  labelId="category-filter-label"
                  value={categoryFilter}
                  onChange={(e) => setCategoryFilter(e.target.value)}
                  label="Категория"
                  startAdornment={
                    <InputAdornment position="start">
                      <FilterListIcon />
                    </InputAdornment>
                  }
                >
                  {categories.map((category) => (
                    <MenuItem key={category} value={category}>
                      {category === 'all' ? 'Все категории' : category}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
          </Grid>
        </CardContent>
      </Card>

      <Card>
        <CardHeader title="Список команд" />
        <Divider />
        <CardContent>
          {loading ? (
            <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
              <CircularProgress />
            </Box>
          ) : filteredCommands.length === 0 ? (
            <Typography variant="body1" sx={{ textAlign: 'center', p: 3 }}>
              Команды не найдены
            </Typography>
          ) : (
            <List>
              {filteredCommands.map((command) => (
                <React.Fragment key={command.name}>
                  <ListItem>
                    <ListItemText
                      primary={
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                          <Typography variant="subtitle1">{command.name}</Typography>
                          <Chip
                            label={command.category}
                            size="small"
                            color="primary"
                            variant="outlined"
                          />
                        </Box>
                      }
                      secondary={
                        <>
                          <Typography variant="body2" color="text.secondary">
                            {command.description}
                          </Typography>
                          <Typography variant="body2" color="text.secondary">
                            Использование: <code>{command.usage}</code>
                          </Typography>
                        </>
                      }
                    />
                    <ListItemSecondaryAction>
                      <Switch
                        edge="end"
                        checked={command.enabled}
                        onChange={() => handleToggleCommand(command)}
                      />
                    </ListItemSecondaryAction>
                  </ListItem>
                  <Divider variant="inset" component="li" />
                </React.Fragment>
              ))}
            </List>
          )}
        </CardContent>
      </Card>

      <Snackbar
        open={snackbar.open}
        autoHideDuration={6000}
        onClose={handleCloseSnackbar}
      >
        <Alert
          onClose={handleCloseSnackbar}
          severity={snackbar.severity}
          sx={{ width: '100%' }}
        >
          {snackbar.message}
        </Alert>
      </Snackbar>
    </Box>
  );
};

export default Commands;