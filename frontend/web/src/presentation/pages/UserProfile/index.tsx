import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import {
  Box,
  Typography,
  Paper,
  Avatar,
  Grid,
  Container,
  Tabs,
  Tab,
  Card,
  CardContent,
  CardMedia,
  CircularProgress,
  Alert,
} from '@mui/material';
import { apiService } from '../../services/api';
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
  is_favorite: boolean;
  rating: number;
  owner: {
    id: string;
    name: string;
  };
}

const UserProfile = () => {
  const { id } = useParams<{ id: string }>();
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');
  const [activeTab, setActiveTab] = useState(0);
  const [userPlaylists, setUserPlaylists] = useState<Playlist[]>([]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const [userData, playlistsData] = await Promise.all([
          apiService.get(`/users/${id}`),
          apiService.get(`/users/${id}/playlists`),
        ]);

        // Load playlist covers
        const playlistsWithCovers = await Promise.all(
          playlistsData.map(async (playlist: Playlist) => {
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

        setUser(userData);
        setUserPlaylists(playlistsWithCovers);
      } catch (err) {
        console.error('Error fetching user data:', err);
        setError('Ошибка при загрузке данных пользователя');
      } finally {
        setLoading(false);
      }
    };

    if (id) {
      fetchData();
    }
  }, [id]);

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
        <Alert severity="error" sx={{ mb: 2 }}>{error || 'Пользователь не найден'}</Alert>
      </Container>
    );
  }

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Paper sx={{ p: 3, mb: 3 }}>
        <Grid container spacing={3} alignItems="center">
          <Grid item>
            <Avatar
              sx={{ width: 120, height: 120 }}
              alt={user.name}
            />
          </Grid>
          <Grid item xs>
            <Typography variant="h4" component="h1" gutterBottom>
              {user.name}
            </Typography>
            <Typography variant="body1" color="text.secondary">
              Дата регистрации: {new Date(user.registration_date).toLocaleDateString()}
            </Typography>
            {user.birth_date && user.birth_date !== '0001-01-01T00:00:00Z' && (
              <Typography variant="body1" color="text.secondary">
                Дата рождения: {new Date(user.birth_date).toLocaleDateString()}
              </Typography>
            )}
          </Grid>
        </Grid>
      </Paper>

      <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 3 }}>
        <Tabs value={activeTab} onChange={(_, newValue) => setActiveTab(newValue)}>
          <Tab label="Плейлисты" />
        </Tabs>
      </Box>

      {activeTab === 0 && (
        <Grid container spacing={3}>
          {userPlaylists.map((playlist) => (
            <Grid item xs={12} sm={6} md={4} key={playlist.id}>
              <PlaylistCard
                id={playlist.id}
                name={playlist.name}
                coverUrl={playlist.coverImage}
                isFavorite={playlist.is_favorite}
                rating={playlist.rating}
                owner={playlist.owner}
              />
            </Grid>
          ))}
          {userPlaylists.length === 0 && (
            <Grid item xs={12}>
              <Typography variant="body1" color="text.secondary" align="center">
                У пользователя пока нет плейлистов
              </Typography>
            </Grid>
          )}
        </Grid>
      )}
    </Container>
  );
};

export default UserProfile; 