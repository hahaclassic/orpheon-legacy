import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Typography,
  Paper,
  Avatar,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Alert,
  Grid,
  Container,
  Tabs,
  Tab,
  Card,
  CardContent,
  CardMedia,
  IconButton,
  CircularProgress,
} from '@mui/material';
import { DatePicker } from '@mui/x-date-pickers';
import LockIcon from '@mui/icons-material/Lock';
import EditIcon from '@mui/icons-material/Edit';
import ImageIcon from '@mui/icons-material/Image';
import axios from 'axios';
import { apiService } from '../../services/api';
import { useAuthContext } from '../../contexts/AuthContext';
import PlaylistCard from '../../components/ContentCards/PlaylistCard';

interface User {
  id: string;
  name: string;
  registration_date: string;
  birth_date: string;
  access_lvl: number;
}

interface Playlist {
  id: string;
  name: string;
  coverImage?: string;
  trackCount: number;
  is_favorite: boolean;
  rating: number;
  owner: {
    id: string;
    name: string;
  };
}

export const Profile = () => {
  const { user, updateUser } = useAuthContext();
  const [open, setOpen] = useState(false);
  const [formData, setFormData] = useState({
    name: user?.name || '',
    birth_date: user?.birth_date && user.birth_date !== '0001-01-01T00:00:00Z' 
      ? new Date(user.birth_date).toISOString().split('T')[0] 
      : '',
  });
  const [error, setError] = useState<string>('');
  const [activeTab, setActiveTab] = useState(0);
  const [passwordDialogOpen, setPasswordDialogOpen] = useState(false);
  const [passwordData, setPasswordData] = useState({
    old: '',
    new: '',
    confirm: '',
  });
  const [loading, setLoading] = useState(true);
  const [success, setSuccess] = useState<string>('');
  const [myPlaylists, setMyPlaylists] = useState<Playlist[]>([]);
  const [favoritePlaylists, setFavoritePlaylists] = useState<Playlist[]>([]);
  const [editDialogOpen, setEditDialogOpen] = useState(false);

  const navigate = useNavigate();

  useEffect(() => {
    if (user) {
      setFormData({
        name: user.name || '',
        birth_date: user.birth_date && user.birth_date !== '0001-01-01T00:00:00Z'
          ? new Date(user.birth_date).toISOString().split('T')[0]
          : '',
      });
    }
  }, [user]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const [userData, myPlaylistsData, favoritePlaylistsData] = await Promise.all([
          apiService.get('/me'),
          apiService.get('/me/playlists'),
          apiService.get('/me/favorites'),
        ]);

        // Загрузка обложек для плейлистов
        const loadPlaylistCovers = async (playlists: Playlist[]) => {
          return Promise.all(
            playlists.map(async (playlist) => {
              try {
                const response = await apiService.get(`/playlists/${playlist.id}/cover`, {
                  responseType: 'blob'
                });
                const imageUrl = URL.createObjectURL(response);
                return { ...playlist, coverImage: imageUrl };
              } catch (error) {
                console.error(`Error loading cover for playlist ${playlist.id}:`, error);
                return playlist;
              }
            })
          );
        };

        const [myPlaylistsWithCovers, favoritePlaylistsWithCovers] = await Promise.all([
          loadPlaylistCovers(myPlaylistsData),
          loadPlaylistCovers(favoritePlaylistsData)
        ]);

        setMyPlaylists(myPlaylistsWithCovers);
        setFavoritePlaylists(favoritePlaylistsWithCovers);
      } catch (err) {
        console.error('Error fetching profile data:', err);
        setError('Ошибка при загрузке данных профиля');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setPasswordData((prev) => ({ ...prev, [name]: value }));
  };

  const handlePasswordSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    if (passwordData.new !== passwordData.confirm) {
      setError('Новые пароли не совпадают');
      return;
    }
    try {
      await apiService.post('/auth/password/update', {
        old: passwordData.old,
        new: passwordData.new
      });
      setSuccess('Пароль успешно изменён!');
      setPasswordData({ old: '', new: '', confirm: '' });
      setPasswordDialogOpen(false);
    } catch (err: any) {
      setError(err?.response?.data?.message || 'Ошибка при смене пароля');
    }
  };

  const handleEditSubmit = async () => {
    try {
      if (!user) return;
      
      const updatedUser = await updateUser({
        id: user.id,
        name: formData.name,
        registration_date: user.registration_date,
        birth_date: formData.birth_date ? new Date(formData.birth_date).toISOString() : '0001-01-01T00:00:00Z',
        access_lvl: user.access_lvl
      });
      setEditDialogOpen(false);
      setSuccess('Профиль успешно обновлен');
    } catch (err) {
      setError('Ошибка при обновлении профиля');
    }
  };

  const handlePlaylistClick = (playlistId: string) => {
    navigate(`/playlists/${playlistId}`);
  };

  if (loading) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
          <CircularProgress />
        </Box>
      </Container>
    );
  }

  if (error || !user) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Alert severity="error" sx={{ mb: 2 }}>{error || 'Профиль не найден'}</Alert>
      </Container>
    );
  }

  return (
    <Container 
      maxWidth="lg" 
      sx={{ 
        py: 4,
        height: 'calc(100vh - 90px)', // Высота экрана минус высота плеера
        display: 'flex',
        flexDirection: 'column'
      }}
    >
      {/* User Info Section */}
      <Box sx={{ mb: 4 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Typography variant="h4" component="h1">
            Профиль
          </Typography>
          <Box>
            <Button
              variant="outlined"
              startIcon={<EditIcon />}
              onClick={() => setEditDialogOpen(true)}
              sx={{ mr: 2 }}
            >
              Редактировать
            </Button>
            <Button
              variant="outlined"
              onClick={() => setPasswordDialogOpen(true)}
            >
              Изменить пароль
            </Button>
          </Box>
        </Box>

        <Grid container spacing={3}>
          <Grid item xs={12} md={6}>
            <Card>
              <CardContent>
                <Box sx={{ display: 'flex', alignItems: 'center', mb: 3 }}>
                  <Avatar
                    sx={{
                      width: 100,
                      height: 100,
                      bgcolor: 'primary.main',
                      fontSize: '2rem',
                      mr: 2
                    }}
                  >
                    {user.name ? user.name.charAt(0).toUpperCase() : '?'}
                  </Avatar>
                  <Box>
                    <Typography variant="h5" gutterBottom>
                      {user.name || 'Без имени'}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      {user.access_lvl === 1 ? 'Администратор' : 'Пользователь'}
                    </Typography>
                  </Box>
                </Box>
                <Box sx={{ mb: 2 }}>
                  <Typography variant="subtitle2" color="text.secondary">
                    Дата регистрации
                  </Typography>
                  <Typography variant="body1">
                    {new Date(user.registration_date).toLocaleDateString()}
                  </Typography>
                </Box>
                <Box sx={{ mb: 2 }}>
                  <Typography variant="subtitle2" color="text.secondary">
                    Дата рождения
                  </Typography>
                  <Typography variant="body1">
                    {user.birth_date === '0001-01-01T00:00:00Z' 
                      ? 'Не указана' 
                      : new Date(user.birth_date).toLocaleDateString()}
                  </Typography>
                </Box>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      </Box>

      {/* Tabs */}
      <Box sx={{ borderBottom: 1, borderColor: 'divider', bgcolor: 'background.paper' }}>
        <Container maxWidth="xl">
          <Tabs 
            value={activeTab} 
            onChange={(_, newValue) => setActiveTab(newValue)}
            sx={{ minHeight: 64 }}
          >
            <Tab label="Мои плейлисты" />
            <Tab label="Избранные плейлисты" />
          </Tabs>
        </Container>
      </Box>

      {/* Content */}
      <Box sx={{ 
        flex: 1, 
        py: 4,
        overflow: 'auto' // Добавляем прокрутку при необходимости
      }}>
        <Container maxWidth="xl">
          {activeTab === 0 && (
            <Box>
              <Typography variant="h5" gutterBottom>Мои плейлисты</Typography>
              <Grid container spacing={3}>
                {myPlaylists.length === 0 ? (
                  <Grid item xs={12}>
                    <Typography variant="h6" textAlign="center" color="text.secondary">
                      У вас пока нет плейлистов
                    </Typography>
                  </Grid>
                ) : (
                  myPlaylists.map((playlist) => (
                    <Grid item xs={12} sm={6} md={3} key={playlist.id}>
                      <PlaylistCard
                        id={playlist.id}
                        name={playlist.name}
                        coverUrl={playlist.coverImage}
                        isFavorite={playlist.is_favorite}
                        rating={playlist.rating || 0}
                        owner={playlist.owner}
                        onFavoriteChange={(isFavorite) => {
                          // Обновляем локальное состояние плейлиста
                          setMyPlaylists(prev => prev.map(p => 
                            p.id === playlist.id 
                              ? { 
                                  ...p, 
                                  is_favorite: isFavorite,
                                  rating: isFavorite ? (p.rating || 0) + 1 : (p.rating || 0) - 1
                                }
                              : p
                          ));
                        }}
                      />
                    </Grid>
                  ))
                )}
              </Grid>
            </Box>
          )}

          {activeTab === 1 && (
            <Box>
              <Typography variant="h5" gutterBottom>Избранные плейлисты</Typography>
              <Grid container spacing={3}>
                {favoritePlaylists.length === 0 ? (
                  <Grid item xs={12}>
                    <Typography variant="h6" textAlign="center" color="text.secondary">
                      У вас пока нет избранных плейлистов
                    </Typography>
                  </Grid>
                ) : (
                  favoritePlaylists.map((playlist) => (
                    <Grid item xs={12} sm={6} md={3} key={playlist.id}>
                      <PlaylistCard
                        id={playlist.id}
                        name={playlist.name}
                        coverUrl={playlist.coverImage}
                        isFavorite={playlist.is_favorite}
                        rating={playlist.rating || 0}
                        owner={playlist.owner}
                        onFavoriteChange={(isFavorite) => {
                          // Обновляем локальное состояние плейлиста
                          setFavoritePlaylists(prev => prev.map(p => 
                            p.id === playlist.id 
                              ? { 
                                  ...p, 
                                  is_favorite: isFavorite,
                                  rating: isFavorite ? (p.rating || 0) + 1 : (p.rating || 0) - 1
                                }
                              : p
                          ));
                        }}
                      />
                    </Grid>
                  ))
                )}
              </Grid>
            </Box>
          )}
        </Container>
      </Box>

      {/* Password Change Dialog */}
      <Dialog open={passwordDialogOpen} onClose={() => setPasswordDialogOpen(false)}>
        <DialogTitle>Изменение пароля</DialogTitle>
        <form onSubmit={handlePasswordSubmit}>
          <DialogContent>
            {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
            {success && <Alert severity="success" sx={{ mb: 2 }}>{success}</Alert>}
            <TextField
              autoFocus
              margin="dense"
              name="old"
              label="Текущий пароль"
              type="password"
              fullWidth
              value={passwordData.old}
              onChange={handlePasswordChange}
              required
            />
            <TextField
              margin="dense"
              name="new"
              label="Новый пароль"
              type="password"
              fullWidth
              value={passwordData.new}
              onChange={handlePasswordChange}
              required
            />
            <TextField
              margin="dense"
              name="confirm"
              label="Подтвердите новый пароль"
              type="password"
              fullWidth
              value={passwordData.confirm}
              onChange={handlePasswordChange}
              required
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setPasswordDialogOpen(false)}>Отмена</Button>
            <Button type="submit" variant="contained">
              Сохранить
            </Button>
          </DialogActions>
        </form>
      </Dialog>

      {/* Edit Dialog */}
      <Dialog open={editDialogOpen} onClose={() => setEditDialogOpen(false)}>
        <DialogTitle>Редактировать профиль</DialogTitle>
        <DialogContent>
          <Box sx={{ pt: 2 }}>
            <TextField
              fullWidth
              label="Имя пользователя"
              value={formData.name}
              onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
              sx={{ mb: 2 }}
            />
            <DatePicker
              label="Дата рождения"
              value={formData.birth_date ? new Date(formData.birth_date) : null}
              onChange={(date) => {
                setFormData(prev => ({
                  ...prev,
                  birth_date: date ? date.toISOString().split('T')[0] : ''
                }));
              }}
              slotProps={{
                textField: {
                  fullWidth: true,
                  margin: 'dense'
                }
              }}
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setEditDialogOpen(false)}>Отмена</Button>
          <Button onClick={handleEditSubmit} variant="contained">
            Сохранить
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
};

export default Profile; 