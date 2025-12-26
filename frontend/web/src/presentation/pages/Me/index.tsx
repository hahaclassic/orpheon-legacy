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
  IconButton,
  Divider,
} from '@mui/material';
import LockIcon from '@mui/icons-material/Lock';
import EditIcon from '@mui/icons-material/Edit';
import FavoriteIcon from '@mui/icons-material/Favorite';
import PlaylistPlayIcon from '@mui/icons-material/PlaylistPlay';
import PersonIcon from '@mui/icons-material/Person';
import api from '../../../core/infrastructure/services/api';

interface UserProfile {
  username: string;
  name?: string;
  email: string;
  avatar?: string;
  birthDate?: string;
  createdAt?: string;
}

interface Playlist {
  id: string;
  name: string;
  isPublic: boolean;
  coverImage?: string;
}

const Me = () => {
  const [activeTab, setActiveTab] = useState(0);
  const [passwordDialogOpen, setPasswordDialogOpen] = useState(false);
  const [passwordData, setPasswordData] = useState({
    old: '',
    new: '',
    confirm: '',
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [profileLoading, setProfileLoading] = useState(true);
  const [playlists, setPlaylists] = useState<Playlist[]>([]);
  const [playlistsLoading, setPlaylistsLoading] = useState(true);
  const [favorites, setFavorites] = useState<Playlist[]>([]);
  const [favoritesLoading, setFavoritesLoading] = useState(true);

  const navigate = useNavigate();

  useEffect(() => {
    setProfileLoading(true);
    api.get('/me')
      .then(res => setProfile(res.data))
      .catch(() => setProfile(null))
      .finally(() => setProfileLoading(false));
  }, []);

  useEffect(() => {
    setPlaylistsLoading(true);
    api.get('/me/playlists')
      .then(res => setPlaylists(res.data))
      .catch(() => setPlaylists([]))
      .finally(() => setPlaylistsLoading(false));
  }, []);

  useEffect(() => {
    setFavoritesLoading(true);
    api.get('/me/favorites')
      .then(res => setFavorites(res.data))
      .catch(() => setFavorites([]))
      .finally(() => setFavoritesLoading(false));
  }, []);

  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setPasswordData((prev) => ({ ...prev, [name]: value }));
  };

  const handlePasswordSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccess(null);
    if (passwordData.new !== passwordData.confirm) {
      setError('Новые пароли не совпадают');
      return;
    }
    setLoading(true);
    try {
      await api.post('/auth/password/update', { old: passwordData.old, new: passwordData.new });
      setSuccess('Пароль успешно изменён!');
      setPasswordData({ old: '', new: '', confirm: '' });
      setPasswordDialogOpen(false);
    } catch (err: any) {
      setError(err?.response?.data?.message || 'Ошибка при смене пароля');
    } finally {
      setLoading(false);
    }
  };

  const handlePlaylistClick = (playlistId: string) => {
    navigate(`/playlists/${playlistId}`);
  };

  const renderPlaylistGrid = (playlists: Playlist[], loading: boolean) => (
    <Grid container spacing={3} sx={{ width: '100%', m: 0 }}>
      {loading ? (
        <Typography color="text.secondary" sx={{ ml: 2 }}>Загрузка...</Typography>
      ) : playlists?.length === 0 ? (
        <Typography color="text.secondary" sx={{ ml: 2 }}>Нет плейлистов</Typography>
      ) : (
        (playlists || []).map((playlist) => (
          <Grid item xs={12} sm={6} md={4} lg={3} key={playlist.id}>
            <Paper 
              sx={{ 
                p: 2,
                height: '100%',
                display: 'flex',
                flexDirection: 'column',
                cursor: 'pointer',
                transition: 'all 0.2s ease-in-out',
                '&:hover': {
                  transform: 'translateY(-4px)',
                  boxShadow: 4,
                },
              }}
              onClick={() => handlePlaylistClick(playlist.id)}
            >
              <Box
                sx={{
                  width: '100%',
                  aspectRatio: '1',
                  mb: 2,
                  bgcolor: 'primary.dark',
                  borderRadius: 2,
                  overflow: 'hidden',
                  position: 'relative',
                }}
              >
                {playlist.coverImage ? (
                  <Box
                    component="img"
                    src={playlist.coverImage}
                    alt={playlist.name}
                    sx={{
                      width: '100%',
                      height: '100%',
                      objectFit: 'cover',
                    }}
                  />
                ) : (
                  <Box
                    sx={{
                      width: '100%',
                      height: '100%',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      bgcolor: 'primary.main',
                    }}
                  >
                    <PlaylistPlayIcon sx={{ fontSize: 48, color: 'white' }} />
                  </Box>
                )}
                {!playlist.isPublic && (
                  <Box
                    sx={{
                      position: 'absolute',
                      top: 8,
                      right: 8,
                      bgcolor: 'rgba(0, 0, 0, 0.6)',
                      borderRadius: 1,
                      px: 1,
                      py: 0.5,
                    }}
                  >
                    <Typography variant="caption" color="white">
                      Приватный
                    </Typography>
                  </Box>
                )}
              </Box>
              <Typography variant="h6" noWrap>{playlist.name}</Typography>
            </Paper>
          </Grid>
        ))
      )}
    </Grid>
  );

  return (
    <Box sx={{ 
      width: '100%', 
      minHeight: '100vh', 
      bgcolor: 'background.default',
      display: 'flex',
      flexDirection: 'column'
    }}>
      {/* Header */}
      <Box
        sx={{
          bgcolor: 'primary.dark',
          py: 6,
          position: 'relative',
          overflow: 'hidden',
          '&::before': {
            content: '""',
            position: 'absolute',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            background: 'linear-gradient(45deg, rgba(0,0,0,0.2) 0%, rgba(0,0,0,0) 100%)',
            zIndex: 1,
          },
        }}
      >
        <Container maxWidth="xl" sx={{ position: 'relative', zIndex: 2 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 4 }}>
            <Avatar
              sx={{
                width: 150,
                height: 150,
                fontSize: 64,
                border: '4px solid',
                borderColor: 'white',
                bgcolor: 'primary.main',
                boxShadow: 4,
              }}
              src={profile?.avatar}
            >
              {profile?.username?.[0]}
            </Avatar>
            <Box sx={{ color: 'white' }}>
              <Typography variant="h3" fontWeight={700} gutterBottom>
                {profileLoading ? 'Загрузка...' : profile?.name || profile?.username || '—'}
              </Typography>
              <Typography variant="h6" sx={{ opacity: 0.9 }}>
                {profileLoading ? '' : profile?.email || ''}
              </Typography>
              {profile?.birthDate && (
                <Typography variant="body1" sx={{ opacity: 0.8, mt: 1 }}>
                  Дата рождения: {new Date(profile.birthDate).toLocaleDateString()}
                </Typography>
              )}
              {profile?.createdAt && (
                <Typography variant="body1" sx={{ opacity: 0.8 }}>
                  Дата регистрации: {new Date(profile.createdAt).toLocaleDateString()}
                </Typography>
              )}
            </Box>
          </Box>
        </Container>
      </Box>

      {/* Tabs */}
      <Box sx={{ borderBottom: 1, borderColor: 'divider', bgcolor: 'background.paper' }}>
        <Container maxWidth="xl">
          <Tabs 
            value={activeTab} 
            onChange={(_, newValue) => setActiveTab(newValue)}
            sx={{ 
              minHeight: 72,
              '& .MuiTab-root': {
                minHeight: 72,
                fontSize: '1rem',
                textTransform: 'none',
                fontWeight: 500,
              },
            }}
          >
            <Tab 
              icon={<PersonIcon />} 
              iconPosition="start" 
              label="Основная информация" 
            />
            <Tab 
              icon={<PlaylistPlayIcon />} 
              iconPosition="start" 
              label="Мои плейлисты" 
            />
            <Tab 
              icon={<FavoriteIcon />} 
              iconPosition="start" 
              label="Избранное" 
            />
          </Tabs>
        </Container>
      </Box>

      {/* Content */}
      <Box sx={{ flex: 1, py: 4 }}>
        <Container maxWidth="xl">
          {activeTab === 0 && (
            <Box>
              <Paper sx={{ p: 3 }}>
                <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 3 }}>
                  <Typography variant="h5">Основная информация</Typography>
                  <IconButton color="primary">
                    <EditIcon />
                  </IconButton>
                </Box>
                <Divider sx={{ mb: 3 }} />
                <Grid container spacing={3}>
                  <Grid item xs={12} md={6}>
                    <Typography variant="subtitle2" color="text.secondary" gutterBottom>
                      Имя пользователя
                    </Typography>
                    <Typography variant="body1">
                      {profile?.username || '—'}
                    </Typography>
                  </Grid>
                  <Grid item xs={12} md={6}>
                    <Typography variant="subtitle2" color="text.secondary" gutterBottom>
                      Email
                    </Typography>
                    <Typography variant="body1">
                      {profile?.email || '—'}
                    </Typography>
                  </Grid>
                  {profile?.birthDate && (
                    <Grid item xs={12} md={6}>
                      <Typography variant="subtitle2" color="text.secondary" gutterBottom>
                        Дата рождения
                      </Typography>
                      <Typography variant="body1">
                        {new Date(profile.birthDate).toLocaleDateString()}
                      </Typography>
                    </Grid>
                  )}
                  {profile?.createdAt && (
                    <Grid item xs={12} md={6}>
                      <Typography variant="subtitle2" color="text.secondary" gutterBottom>
                        Дата регистрации
                      </Typography>
                      <Typography variant="body1">
                        {new Date(profile.createdAt).toLocaleDateString()}
                      </Typography>
                    </Grid>
                  )}
                </Grid>
                <Box sx={{ mt: 4 }}>
                  <Button
                    variant="outlined"
                    startIcon={<LockIcon />}
                    onClick={() => setPasswordDialogOpen(true)}
                  >
                    Изменить пароль
                  </Button>
                </Box>
              </Paper>
            </Box>
          )}

          {activeTab === 1 && (
            <Box>
              <Typography variant="h5" gutterBottom>Мои плейлисты</Typography>
              {renderPlaylistGrid(playlists, playlistsLoading)}
            </Box>
          )}

          {activeTab === 2 && (
            <Box>
              <Typography variant="h5" gutterBottom>Избранное</Typography>
              {renderPlaylistGrid(favorites, favoritesLoading)}
            </Box>
          )}
        </Container>
      </Box>

      {/* Password Change Dialog */}
      <Dialog 
        open={passwordDialogOpen} 
        onClose={() => setPasswordDialogOpen(false)}
        maxWidth="xs"
        fullWidth
      >
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
            <Button type="submit" variant="contained" disabled={loading}>
              {loading ? 'Сохранение...' : 'Сохранить'}
            </Button>
          </DialogActions>
        </form>
      </Dialog>
    </Box>
  );
};

export default Me; 