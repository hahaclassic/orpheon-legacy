import { useState, useEffect, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Box,
  Typography,
  Paper,
  IconButton,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  CircularProgress,
  Alert,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Tooltip,
  Menu,
  MenuItem,
  Switch,
  Divider,
  Container,
  Grid,
  Card,
  CardMedia,
  List,
  ListItem,
  ListItemText,
  FormControlLabel,
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import FavoriteIcon from '@mui/icons-material/Favorite';
import FavoriteBorderIcon from '@mui/icons-material/FavoriteBorder';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import PhotoCameraIcon from '@mui/icons-material/PhotoCamera';
import DeleteIcon from '@mui/icons-material/Delete';
import { ArrowBack, PlayArrow, Pause, Add as AddIcon, Check, MoreVert } from '@mui/icons-material';
import axios from 'axios';
import type { AxiosResponse } from 'axios';
import { useAuthContext } from '../../contexts/AuthContext';
import { usePlayerContext } from '../../contexts/PlayerContext';
import api from '../../../core/infrastructure/services/api';
import TrackList from '../../components/TrackList';
import type { Track } from '../../types';

interface Artist {
  id: string;
  name: string;
}

interface Playlist {
  id: string;
  name: string;
  description?: string;
  is_private: boolean;
  coverImage?: string;
  created_at: string;
  updated_at: string;
  rating: number;
  is_favorite: boolean;
  owner: {
    id: string;
    name: string;
  };
  tracks: Track[];
}

const formatDuration = (seconds: number) => {
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = seconds % 60;
  return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
};

const PlaylistPage = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { user } = useAuthContext();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [playlist, setPlaylist] = useState<Playlist | null>(null);
  const [loading, setLoading] = useState(true);
  const [tracksLoading, setTracksLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [updatingPrivacy, setUpdatingPrivacy] = useState(false);
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [editForm, setEditForm] = useState({
    name: '',
    description: '',
  });
  const [updatingFavorite, setUpdatingFavorite] = useState(false);
  const [uploadingCover, setUploadingCover] = useState(false);
  const [menuAnchor, setMenuAnchor] = useState<null | HTMLElement>(null);
  const { state, controls } = usePlayerContext();
  const { currentTrack, isPlaying } = state;
  const { startPlayback, togglePlay } = controls;
  const [coverUrl, setCoverUrl] = useState<string | null>(null);
  const [playlists, setPlaylists] = useState<Playlist[]>([]);
  const [menuAnchorEl, setMenuAnchorEl] = useState<null | HTMLElement>(null);
  const [selectedTrack, setSelectedTrack] = useState<Track | null>(null);
  const [trackMenuAnchorEl, setTrackMenuAnchorEl] = useState<null | HTMLElement>(null);
  const [selectedTrackForMenu, setSelectedTrackForMenu] = useState<Track | null>(null);
  const albumCoversRef = useRef<Map<string, string>>(new Map());

  const isOwner = Boolean(user && playlist && user.id === playlist.owner.id);
  const isAuthenticated = Boolean(user);

  useEffect(() => {
    const fetchPlaylistData = async () => {
      if (!id) {
        setError('Playlist ID is missing');
        setLoading(false);
        return;
      }
      
      try {
        setLoading(true);
        const response: AxiosResponse<Playlist> = await api.get(`/playlists/${id}`);
        const data = response.data;
        if (!data) {
          throw new Error('Failed to fetch playlist data');
        }
        setPlaylist(data);
        
        // Get playlist cover
        try {
          const coverResponse: AxiosResponse<Blob> = await api.get(`/playlists/${id}/cover`, {
            responseType: 'blob'
          });
          if (coverResponse.data) {
            const coverUrl = URL.createObjectURL(coverResponse.data);
            setCoverUrl(coverUrl);
          }
        } catch (err) {
          console.error('Error fetching playlist cover:', err);
          setCoverUrl('/default-playlist.png');
        }
        
        // Get playlist tracks
        try {
          setTracksLoading(true);
          const tracksResponse: AxiosResponse<Track[]> = await api.get(`/playlists/${id}/tracks`);
          const tracks = tracksResponse.data;
          
          if (!tracks || !Array.isArray(tracks)) {
            throw new Error('Invalid tracks response');
          }

          // Fetch album covers for each track
          const albumCovers = new Map<string, string>();
          for (const track of tracks) {
            if (track?.album?.id && !albumCovers.has(track.album.id)) {
              try {
                const albumCoverResponse: AxiosResponse<Blob> = await api.get(`/albums/${track.album.id}/cover`, {
                  responseType: 'blob'
                });
                if (albumCoverResponse.data) {
                  const url = URL.createObjectURL(albumCoverResponse.data);
                  albumCovers.set(track.album.id, url);
                }
              } catch (err) {
                console.error('Error fetching album cover:', err);
              }
            }
          }

          // Combine tracks with album cover URLs
          const tracksWithCovers = tracks.map((track: Track) => ({
            ...track,
            coverUrl: track?.album?.id ? albumCovers.get(track.album.id) : undefined
          }));

          setPlaylist(prev => prev ? { ...prev, tracks: tracksWithCovers } : null);
          albumCoversRef.current = albumCovers;
        } catch (err) {
          console.error('Error fetching playlist tracks:', err);
          setPlaylist(prev => prev ? { ...prev, tracks: [] } : null);
        } finally {
          setTracksLoading(false);
        }

        setError(null);
      } catch (error) {
        console.error('Error fetching playlist data:', error);
        setError('Failed to load playlist data');
        setPlaylist(null);
      } finally {
        setLoading(false);
      }
    };

    fetchPlaylistData();

    // Cleanup function to revoke object URLs
    return () => {
      if (albumCoversRef.current) {
        albumCoversRef.current.forEach((url: string) => URL.revokeObjectURL(url));
      }
      if (coverUrl) {
        URL.revokeObjectURL(coverUrl);
      }
    };
  }, [id]);

  useEffect(() => {
    if (isAuthenticated) {
      const fetchPlaylists = async () => {
        try {
          const response: AxiosResponse<Playlist[]> = await api.get('/me/playlists');
          setPlaylists(response.data);
        } catch (err) {
          console.error('Error fetching playlists:', err);
        }
      };

      fetchPlaylists();
    }
  }, [isAuthenticated]);

  const handlePrivacyChange = async (_: React.ChangeEvent<HTMLInputElement>) => {
    if (!playlist) return;

    const newPrivacyValue = !playlist.is_private;

    try {
      setUpdatingPrivacy(true);
      await api.patch(`/playlists/${id}/privacy`, { is_private: newPrivacyValue });
      setPlaylist(prev => prev ? { ...prev, is_private: newPrivacyValue } : null);
    } catch (err) {
      setError('Не удалось изменить настройки приватности');
    } finally {
      setUpdatingPrivacy(false);
    }
  };

  const handleDeletePlaylist = async () => {
    if (!playlist || !window.confirm('Вы уверены, что хотите удалить этот плейлист?')) return;

    try {
      await api.delete(`/playlists/${id}`);
      navigate('/library');
    } catch (err) {
      setError('Не удалось удалить плейлист');
    }
  };

  const handleEditSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!playlist) return;

    try {
      const response = await api.put(`/playlists/${id}`, {
        name: editForm.name,
        description: editForm.description,
        is_private: playlist.is_private
      });
      setPlaylist(response.data);
      setEditDialogOpen(false);
    } catch (err) {
      setError('Не удалось обновить информацию о плейлисте');
    }
  };

  const handleFavoriteClick = async () => {
    if (!playlist || updatingFavorite) return;

    try {
      setUpdatingFavorite(true);
      if (playlist.is_favorite) {
        await api.delete(`/me/favorites/${id}`);
        setPlaylist(prev => prev ? {
          ...prev,
          is_favorite: false,
          rating: prev.rating - 1
        } : null);
      } else {
        await api.post(`/me/favorites/${id}`);
        setPlaylist(prev => prev ? {
          ...prev,
          is_favorite: true,
          rating: prev.rating + 1
        } : null);
      }
    } catch (err) {
      setError('Не удалось обновить статус избранного');
    } finally {
      setUpdatingFavorite(false);
    }
  };

  const handleCoverClick = () => {
    fileInputRef.current?.click();
  };

  const handleCoverChange = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file || !id) return;

    try {
      setUploadingCover(true);
      const formData = new FormData();
      formData.append('cover', file);

      await api.post(`/playlists/${id}/cover`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });

      // Reload the cover
      const coverResponse = await api.get(`/playlists/${id}/cover`, {
        responseType: 'blob'
      });
      const coverUrl = URL.createObjectURL(coverResponse.data);
      setCoverUrl(coverUrl);
    } catch (err) {
      setError('Не удалось загрузить обложку');
    } finally {
      setUploadingCover(false);
      event.target.value = '';
    }
  };

  const handleDeleteCover = async () => {
    if (!id) return;

    try {
      setUploadingCover(true);
      await api.delete(`/playlists/${id}/cover`);

      // Освобождаем URL если он был
      if (playlist?.coverImage?.startsWith('blob:')) {
        URL.revokeObjectURL(playlist.coverImage);
      }

      setCoverUrl('/default-playlist.png');
    } catch (err) {
      setError('Не удалось удалить обложку');
    } finally {
      setUploadingCover(false);
    }
  };

  const handleAddToPlaylist = (event: React.MouseEvent<HTMLElement>, track: Track) => {
    event.stopPropagation();
    setSelectedTrack(track);
    setMenuAnchorEl(event.currentTarget);
  };

  const handlePlaylistMenuClose = () => {
    setMenuAnchorEl(null);
    setSelectedTrack(null);
  };

  const handlePlaylistSelect = async (playlistId: string) => {
    if (selectedTrack) {
      try {
        const playlist = playlists.find(p => p.id === playlistId);
        const isInPlaylist = playlist && isTrackInPlaylist(playlist, selectedTrack.id);

        if (isInPlaylist) {
          await api.delete(`/playlists/${playlistId}/tracks/${selectedTrack.id}`);
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
          await api.post(`/playlists/${playlistId}/tracks`, { track_id: selectedTrack.id });
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
    handlePlaylistMenuClose();
  };

  const handleMenuOpen = (e: React.MouseEvent<HTMLElement>) => {
    e.stopPropagation();
    setMenuAnchor(e.currentTarget);
  };

  const handleMenuClose = () => {
    setMenuAnchor(null);
  };

  const handleTrackMenuOpen = (event: React.MouseEvent<HTMLElement>, track: Track) => {
    event.stopPropagation();
    console.log('Track data:', track);
    setSelectedTrackForMenu(track);
    setTrackMenuAnchorEl(event.currentTarget);
  };

  const handleTrackMenuClose = () => {
    setTrackMenuAnchorEl(null);
    setSelectedTrackForMenu(null);
  };

  const handleTrackClick = (trackId: string) => {
    const track = playlist?.tracks.find(t => t.id === trackId);
    if (track) {
      if (currentTrack?.id === trackId) {
        togglePlay();
      } else {
        // Добавляем необходимые поля для трека
        const trackWithRequiredFields = {
          ...track,
          audioUrl: `${api.defaults.baseURL}/tracks/${track.id}/audio`,
          album_id: track.album?.id || '',
          album: track.album || {
            id: '',
            title: '',
            label: '',
            license_id: '',
            release_date: ''
          }
        };
        startPlayback(trackWithRequiredFields, playlist?.tracks.map(t => ({
          ...t,
          audioUrl: `${api.defaults.baseURL}/tracks/${t.id}/audio`,
          album_id: t.album?.id || '',
          album: t.album || {
            id: '',
            title: '',
            label: '',
            license_id: '',
            release_date: ''
          }
        })) || []);
      }
    }
  };

  const isTrackInPlaylist = (playlist: Playlist, trackId: string) => {
    return playlist.tracks?.some(track => track.id === trackId) || false;
  };

  const handleTrackReorder = async (sourceIndex: number, destinationIndex: number) => {
    if (!playlist || !isOwner) return;

    try {
      const track = playlist.tracks[sourceIndex];
      await api.patch(`/playlists/${playlist.id}/tracks/${track.id}/position`, {
        position: destinationIndex + 1
      });

      // Update local state
      const newTracks = [...playlist.tracks];
      const [movedTrack] = newTracks.splice(sourceIndex, 1);
      newTracks.splice(destinationIndex, 0, movedTrack);
      setPlaylist(prev => prev ? { ...prev, tracks: newTracks } : null);
    } catch (error) {
      console.error('Error reordering track:', error);
      // Optionally show an error message to the user
    }
  };

  if (loading) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Typography>Loading...</Typography>
      </Container>
    );
  }

  if (error || !playlist) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 2 }}>
          <Typography color="error">{error || 'Playlist not found'}</Typography>
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
        {/* Playlist Header */}
        <Grid item xs={12} md={4}>
          <Card>
            <Box sx={{ position: 'relative' }}>
              <CardMedia
                component="img"
                image={coverUrl || '/default-playlist.png'}
                alt={playlist.name}
                sx={{ 
                  aspectRatio: '1/1',
                  width: '100%',
                  height: 'auto',
                  objectFit: 'cover'
                }}
              />
              {isOwner && (
                <Box sx={{ position: 'absolute', bottom: 8, right: 8, display: 'flex', gap: 1 }}>
                  <IconButton
                    sx={{
                      bgcolor: 'rgba(0, 0, 0, 0.6)',
                      '&:hover': {
                        bgcolor: 'rgba(0, 0, 0, 0.8)',
                      },
                    }}
                    onClick={handleCoverClick}
                    disabled={uploadingCover}
                  >
                    <PhotoCameraIcon sx={{ color: 'white' }} />
                  </IconButton>
                  {coverUrl && coverUrl !== '/default-playlist.png' && (
                    <IconButton
                      sx={{
                        bgcolor: 'rgba(0, 0, 0, 0.6)',
                        '&:hover': {
                          bgcolor: 'rgba(0, 0, 0, 0.8)',
                        },
                      }}
                      onClick={handleDeleteCover}
                      disabled={uploadingCover}
                    >
                      <DeleteIcon sx={{ color: 'white' }} />
                    </IconButton>
                  )}
                </Box>
              )}
            </Box>
          </Card>
          <input
            type="file"
            accept="image/*"
            hidden
            ref={fileInputRef}
            onChange={handleCoverChange}
          />
        </Grid>
        <Grid item xs={12} md={8}>
          <Box sx={{ mb: 2 }}>
            <Typography variant="caption" color="text.secondary" sx={{ display: 'block', mb: 0.5 }}>
              плейлист
            </Typography>
            <Typography 
              variant="h2" 
              component="h1" 
              sx={{ 
                fontWeight: 700,
                fontSize: { xs: '2.5rem', md: '3.5rem' },
                mb: 1
              }}
            >
              {playlist?.name}
            </Typography>

            {/* Description */}
            <Box 
              onClick={() => {
                if (isOwner) {
                  setEditForm({
                    name: playlist?.name || '',
                    description: playlist?.description || ''
                  });
                  setEditDialogOpen(true);
                }
              }}
              sx={{ 
                mb: 3,
                cursor: isOwner ? 'pointer' : 'default',
                '&:hover': isOwner ? {
                  opacity: 0.8
                } : {}
              }}
            >
              {playlist?.description ? (
                <Typography variant="body1" color="text.secondary">
                  {playlist.description}
                </Typography>
              ) : isOwner ? (
                <Typography variant="body1" color="text.secondary" sx={{ fontStyle: 'italic' }}>
                  Добавить описание
                </Typography>
              ) : null}
            </Box>
            
            {/* Owner and Meta information */}
            <Box sx={{ mb: 3 }}>
              <Typography
                variant="h6"
                component="a"
                href={`/users/${playlist.owner.id}`}
                sx={{ 
                  display: 'inline-block',
                  textDecoration: 'none', 
                  color: 'inherit',
                  mb: 0.5,
                  '&:hover': {
                    textDecoration: 'underline',
                  }
                }}
              >
                {playlist.owner.name}
              </Typography>

              <Typography variant="body2" color="text.secondary">
                Создан: {new Date(playlist.created_at).toLocaleString()}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Обновлен: {new Date(playlist.updated_at).toLocaleString()}
              </Typography>
            </Box>
            
            {/* Actions */}
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <IconButton onClick={handleFavoriteClick} disabled={updatingFavorite}>
                  {playlist.is_favorite ? (
                    <FavoriteIcon sx={{ color: 'white' }} />
                  ) : (
                    <FavoriteBorderIcon sx={{ color: 'white' }} />
                  )}
                </IconButton>
                <Typography variant="body2" color="text.secondary">
                  {playlist.rating}
                </Typography>
              </Box>
              {isOwner && (
                <IconButton onClick={handleMenuOpen}>
                  <MoreVertIcon />
                </IconButton>
              )}
            </Box>
          </Box>
        </Grid>
      </Grid>

      {/* Tracks List */}
      <Grid item xs={12}>
        <Typography variant="h5" sx={{ mb: 2 }}>
          Треки
        </Typography>
        {tracksLoading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
            <CircularProgress />
          </Box>
        ) : (
          <TrackList
            tracks={playlist.tracks}
            onTrackClick={handleTrackClick}
            onAddToPlaylist={isAuthenticated ? handleAddToPlaylist : undefined}
            showTrackNumber={true}
            showAlbumLink={true}
            onTrackReorder={isOwner ? handleTrackReorder : undefined}
            isDraggable={isOwner}
          />
        )}
      </Grid>

      {/* Edit Dialog */}
      {isOwner && (
        <Dialog open={editDialogOpen} onClose={() => setEditDialogOpen(false)} maxWidth="sm" fullWidth>
          <form onSubmit={handleEditSubmit}>
            <DialogTitle>Редактировать плейлист</DialogTitle>
            <DialogContent sx={{ display: 'flex', flexDirection: 'column', gap: 2, pt: 2 }}>
              <TextField
                label="Название"
                value={editForm.name}
                onChange={(e) => setEditForm(prev => ({ ...prev, name: e.target.value }))}
                required
                fullWidth
              />
              <TextField
                label="Описание"
                value={editForm.description}
                onChange={(e) => setEditForm(prev => ({ ...prev, description: e.target.value }))}
                multiline
                rows={4}
                fullWidth
              />
            </DialogContent>
            <DialogActions>
              <Button onClick={() => setEditDialogOpen(false)}>Отмена</Button>
              <Button type="submit" variant="contained">Сохранить</Button>
            </DialogActions>
          </form>
        </Dialog>
      )}

      {/* Playlist Selection Menu */}
      {isAuthenticated && (
        <Menu
          anchorEl={menuAnchorEl}
          open={Boolean(menuAnchorEl)}
          onClose={handlePlaylistMenuClose}
          onClick={(e) => e.stopPropagation()}
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
      )}

      {/* Track Menu */}
      {isAuthenticated && (
        <Menu
          anchorEl={trackMenuAnchorEl}
          open={Boolean(trackMenuAnchorEl)}
          onClose={handleTrackMenuClose}
          onClick={(e) => e.stopPropagation()}
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
          {selectedTrackForMenu && (
            <MenuItem 
              onClick={() => {
                handleTrackMenuClose();
                navigate(`/albums/${selectedTrackForMenu.album.id}`);
              }}
            >
              Перейти к альбому
            </MenuItem>
          )}
        </Menu>
      )}

      {/* Menu */}
      {isOwner && (
        <Menu
          anchorEl={menuAnchor}
          open={Boolean(menuAnchor)}
          onClose={handleMenuClose}
          onClick={(e) => e.stopPropagation()}
        >
          <MenuItem onClick={() => {
            handleMenuClose();
            setEditDialogOpen(true);
          }}>
            Редактировать
          </MenuItem>
          <MenuItem>
            <FormControlLabel
              control={
                <Switch
                  checked={playlist.is_private}
                  onChange={handlePrivacyChange}
                  disabled={updatingPrivacy}
                />
              }
              label="Приватный плейлист"
            />
          </MenuItem>
          <Divider />
          <MenuItem 
            onClick={() => {
              handleMenuClose();
              handleDeletePlaylist();
            }}
            sx={{ color: 'error.main' }}
          >
            Удалить плейлист
          </MenuItem>
        </Menu>
      )}
    </Container>
  );
};

export default PlaylistPage;