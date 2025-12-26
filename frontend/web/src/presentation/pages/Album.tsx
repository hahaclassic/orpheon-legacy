import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { 
  Box, 
  Container, 
  Typography, 
  Grid, 
  Card, 
  CardMedia, 
  Chip, 
  Button,
  List,
  ListItem,
  ListItemText,
  Divider,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Select,
  MenuItem,
  Menu,
} from '@mui/material';
import { ArrowBack, PlayArrow, Pause, Add as AddIcon, Check, MoreVert } from '@mui/icons-material';
import { albumService } from '../../core/infrastructure/services/albumService';
import type { Album, Artist, Genre } from '../../core/infrastructure/services/albumService';
import { apiService } from '../services/api';
import { usePlayerContext } from '../contexts/PlayerContext';
import TrackList from '../components/TrackList';

interface Track {
  id: string;
  name: string;
  duration: number;
  track_number: number;
  artists: Artist[];
}

interface Playlist {
  id: string;
  name: string;
  tracks: Track[];
}

const AlbumPage = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [album, setAlbum] = useState<Album | null>(null);
  const [tracks, setTracks] = useState<Track[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [coverUrl, setCoverUrl] = useState<string | null>(null);
  const { state, controls } = usePlayerContext();
  const { currentTrack, isPlaying } = state;
  const { startPlayback, togglePlay } = controls;
  const [playlists, setPlaylists] = useState<Playlist[]>([]);
  const [menuAnchorEl, setMenuAnchorEl] = useState<null | HTMLElement>(null);
  const [selectedTrack, setSelectedTrack] = useState<Track | null>(null);
  const [trackMenuAnchorEl, setTrackMenuAnchorEl] = useState<null | HTMLElement>(null);
  const [selectedTrackForMenu, setSelectedTrackForMenu] = useState<Track | null>(null);

  useEffect(() => {
    const fetchAlbumData = async () => {
      if (!id) {
        setError('Album ID is missing');
        setLoading(false);
        return;
      }
      
      try {
        setLoading(true);
        const data = await albumService.getAlbum(id);
        setAlbum(data);
        
        // Получаем обложку альбома
        try {
          const coverResponse = await apiService.get(`/albums/${id}/cover`, {
            responseType: 'blob'
          });
          const coverUrl = URL.createObjectURL(coverResponse);
          setCoverUrl(coverUrl);
        } catch (err) {
          console.error('Error fetching album cover:', err);
          setCoverUrl('/default-album.png');
        }

        // Получаем треки альбома
        try {
          const tracksResponse = await apiService.get(`/albums/${id}/tracks`);
          setTracks(tracksResponse);
        } catch (err) {
          console.error('Error fetching album tracks:', err);
          setTracks([]);
        }

        setError(null);
      } catch (error) {
        console.error('Error fetching album data:', error);
        setError('Failed to load album data');
        setAlbum(null);
      } finally {
        setLoading(false);
      }
    };

    fetchAlbumData();
  }, [id]);

  useEffect(() => {
    const fetchPlaylists = async () => {
      try {
        const response = await apiService.get('/me/playlists');
        setPlaylists(response);
      } catch (err) {
        console.error('Error fetching playlists:', err);
      }
    };

    fetchPlaylists();
  }, []);

  const formatDuration = (seconds: number) => {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
  };

  const handleTrackClick = (trackId: string) => {
    const track = tracks.find(t => t.id === trackId);
    if (track) {
      if (currentTrack?.id === trackId) {
        togglePlay();
      } else {
        startPlayback(track, tracks);
      }
    }
  };

  const handleAddToPlaylist = (event: React.MouseEvent<HTMLElement>, track: Track) => {
    event.stopPropagation();
    setSelectedTrack(track);
    setMenuAnchorEl(event.currentTarget);
  };

  const handlePlaylistSelect = async (playlistId: string) => {
    if (selectedTrack) {
      try {
        const playlist = playlists.find(p => p.id === playlistId);
        const isInPlaylist = playlist && isTrackInPlaylist(playlist, selectedTrack.id);

        if (isInPlaylist) {
          // Удаляем трек из плейлиста
          await apiService.delete(`/playlists/${playlistId}/tracks/${selectedTrack.id}`);
          // Обновляем локальное состояние плейлиста
          setPlaylists(playlists.map(p => {
            if (p.id === playlistId) {
              return {
                ...p,
                tracks: (p.tracks || []).filter(t => t.id !== selectedTrack.id)
              };
            }
            return p;
          }));
        } else {
          // Добавляем трек в плейлист
          await apiService.post(`/playlists/${playlistId}/tracks`, { track_id: selectedTrack.id });
          // Обновляем локальное состояние плейлиста
          setPlaylists(playlists.map(p => {
            if (p.id === playlistId) {
              return {
                ...p,
                tracks: [...(p.tracks || []), selectedTrack]
              };
            }
            return p;
          }));
        }
      } catch (err: any) {
        console.error('Error managing track in playlist:', {
          error: err,
          response: err.response?.data,
          status: err.response?.status,
          playlistId,
          trackId: selectedTrack.id
        });
      }
    }
  };

  const handleMenuClose = () => {
    setMenuAnchorEl(null);
    setSelectedTrack(null);
  };

  const isTrackInPlaylist = (playlist: Playlist, trackId: string) => {
    return playlist.tracks?.some(track => track.id === trackId) || false;
  };

  const handleTrackMenuOpen = (event: React.MouseEvent<HTMLElement>, track: Track) => {
    event.stopPropagation();
    setSelectedTrackForMenu(track);
    setTrackMenuAnchorEl(event.currentTarget);
  };

  const handleTrackMenuClose = () => {
    setTrackMenuAnchorEl(null);
    setSelectedTrackForMenu(null);
  };

  if (loading) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Typography>Loading...</Typography>
      </Container>
    );
  }

  if (error || !album) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 2 }}>
          <Typography color="error">{error || 'Album not found'}</Typography>
          <Button
            startIcon={<ArrowBack />}
            onClick={() => navigate(-1)}
            variant="outlined"
          >
            Назад
          </Button>
        </Box>
      </Container>
    );
  }

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Button
        startIcon={<ArrowBack />}
        onClick={() => navigate(-1)}
        variant="outlined"
        sx={{ mb: 3 }}
      >
        Назад
      </Button>
      
      <Grid container spacing={4}>
        {/* Album Header */}
        <Grid item xs={12} md={4}>
          <Card>
            <CardMedia
              component="img"
              image={coverUrl || '/default-album.png'}
              alt={album.title}
              sx={{ 
                aspectRatio: '1/1',
                width: '100%',
                height: 'auto',
                objectFit: 'cover'
              }}
            />
          </Card>
        </Grid>
        <Grid item xs={12} md={8}>
          <Box sx={{ mb: 2 }}>
            <Typography variant="h3" component="h1" gutterBottom>
              {album.title}
            </Typography>
            
            {/* Artists */}
            <Box sx={{ mb: 2 }}>
              {album.artists.map((artist, index) => (
                <Typography
                  key={artist.id}
                  variant="h6"
                  component="a"
                  href={`/artists/${artist.id}`}
                  sx={{ 
                    display: 'inline-block',
                    textDecoration: 'none', 
                    color: 'inherit',
                    '&:hover': {
                      textDecoration: 'underline',
                    },
                    '&:not(:last-child)::after': {
                      content: '", "',
                      color: 'text.secondary',
                      marginRight: '4px'
                    }
                  }}
                >
                  {artist.name}
                </Typography>
              ))}
            </Box>

            {/* Genres */}
            {album.genres && album.genres.length > 0 && (
              <Box sx={{ mt: 2, display: 'flex', gap: 1, flexWrap: 'wrap' }}>
                {album.genres.map((genre) => (
                  <Chip 
                    key={genre.id} 
                    label={genre.title} 
                    size="small"
                    sx={{ 
                      backgroundColor: 'primary.light',
                      color: 'primary.contrastText',
                      '&:hover': {
                        backgroundColor: 'primary.main',
                      }
                    }}
                  />
                ))}
              </Box>
            )}

            {/* Release Date */}
            {album.release_date && (
              <Typography variant="body2" color="text.secondary" sx={{ mt: 2 }}>
                Дата релиза: {new Date(album.release_date).toLocaleDateString()}
              </Typography>
            )}

            {/* Label */}
            {album.label && (
              <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                Лейбл: {album.label}
              </Typography>
            )}
          </Box>
        </Grid>

        {/* Tracks */}
        <Grid item xs={12}>
          <Typography variant="h5" sx={{ mb: 2 }}>
            Треки
          </Typography>
          <TrackList
            tracks={tracks}
            currentTrackId={currentTrack?.id}
            isPlaying={isPlaying}
            onTrackClick={handleTrackClick}
            onAddToPlaylist={handleAddToPlaylist}
            showTrackNumber={true}
            showAlbumLink={false}
          />
        </Grid>
      </Grid>

      {/* Track Menu */}
      <Menu
        anchorEl={trackMenuAnchorEl}
        open={Boolean(trackMenuAnchorEl)}
        onClose={handleTrackMenuClose}
        PaperProps={{
          sx: { 
            '& .MuiList-root': {
              padding: 0
            }
          }
        }}
        MenuListProps={{
          sx: {
            padding: 0
          }
        }}
      >
        <MenuItem 
          onClick={() => {
            handleTrackMenuClose();
            navigate(`/albums/${id}`);
          }}
        >
          Перейти к альбому
        </MenuItem>
      </Menu>

      <Menu
        anchorEl={menuAnchorEl}
        open={Boolean(menuAnchorEl)}
        onClose={handleMenuClose}
        PaperProps={{
          sx: { 
            maxHeight: '280px',
            '& .MuiList-root': {
              padding: 0
            }
          }
        }}
        MenuListProps={{
          sx: {
            padding: 0
          }
        }}
      >
        {playlists.map((playlist) => (
          <MenuItem 
            key={playlist.id} 
            onClick={() => handlePlaylistSelect(playlist.id)}
            sx={{ 
              minHeight: '40px',
              '&:hover': {
                backgroundColor: 'action.hover'
              },
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center'
            }}
          >
            <span>{playlist.name}</span>
            {selectedTrack && isTrackInPlaylist(playlist, selectedTrack.id) && (
              <Check sx={{ ml: 1, color: 'white' }} />
            )}
          </MenuItem>
        ))}
      </Menu>
    </Container>
  );
};

export default AlbumPage; 