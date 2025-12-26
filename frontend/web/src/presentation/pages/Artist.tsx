import { useEffect, useState, useRef } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { Box, Container, Typography, Grid, Card, CardContent, CardMedia, List, ListItem, ListItemText, ListItemAvatar, Avatar, CircularProgress, Menu, MenuItem } from '@mui/material';
import { MusicNote, Album, Check } from '@mui/icons-material';
import { api, apiService } from '../services/api';
import AlbumCard from '../components/ContentCards/AlbumCard';
import TrackList from '../components/TrackList';
import { usePlayerContext } from '../contexts/PlayerContext';
import type { Track } from '../types';

interface Artist {
  id: string;
  name: string;
  description: string;
  country: string;
}

interface TrackResponse {
  id: string;
  name: string;
  duration: string;
  album: {
    id: string;
    title: string;
    label: string;
    license_id: string;
    release_date: string;
  };
  artists: Artist[];
  license?: {
    id: string;
    title: string;
    description: string;
    url: string;
  };
  total_streams?: number;
}

interface Genre {
  id: string;
  title: string;
}

interface Album {
  id: string;
  title: string;
  label: string;
  releaseDate: string;
  artists: Artist[];
  genres: Genre[];
  coverUrl?: string;
}

interface Playlist {
  id: string;
  name: string;
  tracks: Track[];
}

const ArtistPage = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { state, controls } = usePlayerContext();
  const { currentTrack, isPlaying } = state;
  const { startPlayback, togglePlay } = controls;
  const [artist, setArtist] = useState<Artist | null>(null);
  const [tracks, setTracks] = useState<Track[]>([]);
  const [albums, setAlbums] = useState<Album[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [avatarUrl, setAvatarUrl] = useState<string | null>(null);
  const [avatarError, setAvatarError] = useState(false);
  const albumCoversRef = useRef<Map<string, string>>(new Map());
  const [playlists, setPlaylists] = useState<Playlist[]>([]);
  const [menuAnchorEl, setMenuAnchorEl] = useState<null | HTMLElement>(null);
  const [selectedTrack, setSelectedTrack] = useState<Track | null>(null);

  useEffect(() => {
    const fetchArtistData = async () => {
      try {
        setLoading(true);
        setError(null);
        setAvatarError(false);

        // Fetch artist data
        const artistData = await api.getArtist(id!);
        setArtist(artistData);

        // Fetch artist's tracks
        const tracksData = await api.getArtistTracks(id!);
        
        // Fetch artist's albums to get cover URLs
        const albumsData = await api.getArtistAlbums(id!);
        console.log('Albums data:', albumsData);
        setAlbums(albumsData || []);

        // Fetch album covers for each track
        const albumCovers = new Map<string, string>();
        for (const track of tracksData || []) {
          if (!albumCovers.has(track.album.id)) {
            try {
              const response = await apiService.get(`/albums/${track.album.id}/cover`, {
                responseType: 'blob'
              });
              const url = URL.createObjectURL(response);
              albumCovers.set(track.album.id, url);
            } catch (err) {
              console.error('Error fetching album cover:', err);
            }
          }
        }

        // Combine tracks with album cover URLs
        const tracksWithCovers = (tracksData || []).map((track: TrackResponse) => {
          return {
            id: track.id,
            name: track.name,
            duration: parseInt(track.duration),
            album_id: track.album.id,
            album: track.album,
            track_number: 0,
            artists: track.artists,
            coverUrl: albumCovers.get(track.album.id),
            license: track.license,
            total_streams: track.total_streams
          };
        });
        console.log('Tracks with covers:', tracksWithCovers);
        
        setTracks(tracksWithCovers);

        // Fetch artist's avatar
        try {
          const response = await apiService.get(`/artists/${id}/avatar`, {
            responseType: 'blob'
          });
          const url = URL.createObjectURL(response);
          setAvatarUrl(url);
        } catch (err) {
          console.error('Error fetching artist avatar:', err);
          setAvatarError(true);
        }

        setLoading(false);
      } catch (error) {
        console.error('Error fetching artist data:', error);
        setError('Failed to load artist data');
        setLoading(false);
      }
    };

    const fetchPlaylists = async () => {
      try {
        const response = await apiService.get('/me/playlists');
        setPlaylists(response);
      } catch (err) {
        console.error('Error fetching playlists:', err);
      }
    };

    if (id) {
      fetchArtistData();
      fetchPlaylists();
    }

    // Cleanup function to revoke object URLs
    return () => {
      albumCoversRef.current.forEach((url: string) => URL.revokeObjectURL(url));
      if (avatarUrl) {
        URL.revokeObjectURL(avatarUrl);
      }
    };
  }, [id]);

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

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="60vh">
        <CircularProgress />
      </Box>
    );
  }

  if (error || !artist) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Typography color="error" variant="h5">
          {error || 'Artist not found'}
        </Typography>
      </Container>
    );
  }

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Grid container spacing={4}>
        {/* Artist Header */}
        <Grid item xs={12}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 3 }}>
            {!avatarError ? (
              <CardMedia
                component="img"
                sx={{ width: 200, height: 200, borderRadius: 2 }}
                image={avatarUrl || `/artists/${artist.id}/avatar`}
                alt={artist.name}
                onError={() => setAvatarError(true)}
              />
            ) : (
              <Box
                sx={{
                  width: 200,
                  height: 200,
                  borderRadius: 2,
                  bgcolor: 'grey.200',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center'
                }}
              >
                <Typography variant="h1" color="text.secondary">
                  {artist.name.charAt(0).toUpperCase()}
                </Typography>
              </Box>
            )}
            <Box>
              <Typography variant="h3" component="h1" gutterBottom>
                {artist.name}
              </Typography>
              <Typography variant="body1" color="text.secondary" paragraph>
                {artist.description}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Country: {artist.country}
              </Typography>
            </Box>
          </Box>
        </Grid>

        {/* Tracks Section */}
        <Grid item xs={12}>
          <Typography variant="h5" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <MusicNote /> Tracks
          </Typography>
          <TrackList
            tracks={tracks}
            onTrackClick={handleTrackClick}
            onAddToPlaylist={handleAddToPlaylist}
            showTrackNumber={false}
            showAlbumLink={true}
          />
        </Grid>

        {/* Albums Section */}
        <Grid item xs={12}>
          <Typography variant="h5" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Album /> Albums
          </Typography>
          <Grid container spacing={2}>
            {albums.map((album) => (
              <Grid item xs={12} sm={6} md={4} lg={3} key={album.id}>
                <AlbumCard
                  id={album.id}
                  title={album.title}
                  label={album.label}
                  release_date={album.releaseDate}
                  artists={album.artists}
                  genres={album.genres}
                />
              </Grid>
            ))}
          </Grid>
        </Grid>

        {/* Playlist Menu */}
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
      </Grid>
    </Container>
  );
};

export default ArtistPage; 