import React from 'react';
import { Box, Typography, Grid, Card, CardContent } from '@mui/material';
import { useNavigate } from 'react-router-dom';

export const AdminPanel: React.FC = () => {
  const navigate = useNavigate();

  const menuItems = [
    {
      title: 'Управление альбомами',
      description: 'Создание, редактирование и управление музыкальными альбомами',
      path: '/admin/albums',
    },
    {
      title: 'Управление артистами',
      description: 'Создание и управление профилями артистов',
      path: '/admin/artists',
    },
    {
      title: 'Управление жанрами',
      description: 'Добавление, редактирование и удаление музыкальных жанров',
      path: '/admin/genres',
    },
    {
      title: 'Управление лицензиями',
      description: 'Управление музыкальными лицензиями и ценами',
      path: '/admin/licenses',
    },
  ];

  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h4" gutterBottom>
        Панель администратора
      </Typography>
      <Grid container spacing={3} sx={{ ml: 0 }}>
        {menuItems.map((item) => (
          <Grid item xs={12} sm={6} md={3} key={item.path}>
            <Card 
              sx={{ 
                height: '100%',
                cursor: 'pointer',
                '&:hover': {
                  boxShadow: 6,
                },
              }}
              onClick={() => navigate(item.path)}
            >
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  {item.title}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  {item.description}
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>
    </Box>
  );
}; 