import { useState, useEffect } from "react";
import {
  Box,
  Container,
  Typography,
  TextField,
  InputAdornment,
  Grid,
  Button,
  CircularProgress,
  Alert,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Autocomplete,
  Card,
  CardMedia,
  CardContent,
  IconButton,
  Chip,
  Menu,
} from "@mui/material";
import { useNavigate, useSearchParams } from "react-router-dom";
import SearchIcon from "@mui/icons-material/Search";
import ImageIcon from "@mui/icons-material/Image";
import FavoriteIcon from '@mui/icons-material/Favorite';
import FavoriteBorderIcon from '@mui/icons-material/FavoriteBorder';
import CheckIcon from '@mui/icons-material/Check';
import { apiService } from "../../../presentation/services/api";
import LoadingSpinner from "../../components/LoadingSpinner";
import ErrorBoundary from "../../components/ErrorBoundary";
import AlbumCard from "../../components/ContentCards/AlbumCard";
import ArtistCard from "../../components/ContentCards/ArtistCard";
import PlaylistCard from "../../components/ContentCards/PlaylistCard";
import TrackList from "../../components/TrackList";
import { usePlayerContext } from "../../contexts/PlayerContext";
import type { Track as TrackType } from '../../types';
import { useAuthContext } from "../../contexts/AuthContext";
import axios from 'axios';

type ContentType = "track" | "album" | "playlist" | "artist";

interface Genre {
  id: string;
  title: string;
}

interface License {
  id: string;
  title: string;
  description: string;
  url: string;
}

interface Artist {
  id: string;
  name: string;
  description: string;
  country: string;
}

interface Album {
  id: string;
  title: string;
  label: string;
  license_id: string;
  release_date: string;
}

interface Track {
  id: string;
  name: string;
  duration: number;
  explicit: boolean;
  track_number: number;
  total_streams: number;
  genre: Genre;
  license: License;
  album: Album;
  artists: Artist[];
  albumCoverUrl?: string;
}

const API_URL = import.meta.env.VITE_API_URL;

const Search = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const { state, controls } = usePlayerContext();
  const { currentTrack, isPlaying } = state;
  const { startPlayback, togglePlay } = controls;
  const [searchQuery, setSearchQuery] = useState("");
  const [country, setCountry] = useState(searchParams.get('country') || '');
  const [genre, setGenre] = useState<Genre | null>(null);
  const [contentType, setContentType] = useState<ContentType>("track");
  const [albumCovers, setAlbumCovers] = useState<Map<string, string>>(new Map());
  const [playlistCovers, setPlaylistCovers] = useState<Map<string, string>>(new Map());
  const [artistAvatars, setArtistAvatars] = useState<Map<string, string>>(new Map());
  const [genres, setGenres] = useState<Genre[]>([]);
  const [loadingGenres, setLoadingGenres] = useState(false);
  const [results, setResults] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();
  const [playlists, setPlaylists] = useState<any[]>([]);
  const [menuAnchorEl, setMenuAnchorEl] = useState<null | HTMLElement>(null);
  const [selectedTrack, setSelectedTrack] = useState<TrackType | null>(null);
  const { isAuthenticated } = useAuthContext();

  const handleTrackClick = (trackId: string) => {
    const track = results.find((t: TrackType) => t.id === trackId);
    if (track) {
      if (currentTrack?.id === trackId) {
        togglePlay();
      } else {
        startPlayback(track, results);
      }
    }
  };

  // Функция для очистки URL-объектов
  const cleanupUrls = (urlMap: Map<string, string>) => {
    urlMap.forEach((url) => URL.revokeObjectURL(url));
  };

  // Очистка всех URL-объектов при размонтировании компонента
  useEffect(() => {
    return () => {
      cleanupUrls(albumCovers);
      cleanupUrls(playlistCovers);
      cleanupUrls(artistAvatars);
    };
  }, []);

  // Очистка URL-объектов при изменении параметров поиска
  useEffect(() => {
    cleanupUrls(albumCovers);
    cleanupUrls(playlistCovers);
    cleanupUrls(artistAvatars);
    
    setAlbumCovers(new Map());
    setPlaylistCovers(new Map());
    setArtistAvatars(new Map());
    setResults([]);

    if (searchQuery || contentType) {
      handleSearch();
    }
  }, [searchQuery, country, genre, contentType]);

  useEffect(() => {
    const fetchGenres = async () => {
      try {
        console.log('[Search] Fetching genres...');
        setLoadingGenres(true);
        const response = await axios.get(`${API_URL}/genres`, {
          headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
          }
        });
        console.log('[Search] Genres received:', response.data);
        setGenres(response.data);
        
        const genreId = searchParams.get('genre_id');
        if (genreId) {
          const foundGenre = response.data.find((g: Genre) => g.id === genreId);
          if (foundGenre) {
            setGenre(foundGenre);
          }
        }
      } catch (err: any) {
        console.error('[Search] Error fetching genres:', {
          message: err.message,
          status: err.response?.status,
          data: err.response?.data,
          headers: err.response?.headers,
          config: err.config
        });
      } finally {
        setLoadingGenres(false);
      }
    };

    fetchGenres();
  }, [searchParams]);

  const handleSearch = async (e?: React.FormEvent) => {
    if (e) {
      e.preventDefault();
    }

    try {
      console.log('[Search] Starting search with params:', {
        query: searchQuery,
        country,
        genre: genre?.id,
        contentType
      });

      setLoading(true);
      setError(null);

      const params = new URLSearchParams();
      if (searchQuery) params.append('query', searchQuery);
      if (country) params.append('country', country);
      if (genre?.id) params.append('genre_id', genre.id);
      if (contentType) params.append('type', contentType);

      console.log('[Search] Making request to:', `${API_URL}/search?${params.toString()}`);
      const response = await axios.get(`${API_URL}/search`, {
        params,
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        }
      });

      console.log('[Search] Response received:', response.data);
      setResults(response.data);

      if (contentType === 'track') {
        await fetchAlbumCovers(response.data);
      } else if (contentType === 'playlist') {
        await fetchPlaylistCovers(response.data);
      } else if (contentType === 'artist') {
        await fetchArtistAvatars(response.data);
      }
    } catch (err: any) {
      console.error('[Search] Error:', err);
      setError(err.response?.data?.message || 'Произошла ошибка при поиске');
    } finally {
      setLoading(false);
    }
  };

  const fetchAlbumCovers = async (tracks: Track[]) => {
    for (const track of tracks) {
      if (track.album?.id && !albumCovers.has(track.album.id)) {
        try {
          const response = await axios.get(`${API_URL}/albums/${track.album.id}/cover`, {
            responseType: 'blob'
          });
          const url = URL.createObjectURL(response.data);
          setAlbumCovers(prev => new Map(prev).set(track.album.id, url));
        } catch (err) {
          console.error(`Failed to load cover for album ${track.album.id}:`, err);
        }
      }
    }
  };

  const fetchPlaylistCovers = async (playlists: any[]) => {
    for (const playlist of playlists) {
      if (playlist.id && !playlistCovers.has(playlist.id)) {
        try {
          const response = await axios.get(`${API_URL}/playlists/${playlist.id}/cover`, {
            responseType: 'blob'
          });
          const url = URL.createObjectURL(response.data);
          setPlaylistCovers(prev => new Map(prev).set(playlist.id, url));
        } catch (err) {
          console.error(`Failed to load cover for playlist ${playlist.id}:`, err);
        }
      }
    }
  };

  const fetchArtistAvatars = async (artists: any[]) => {
    for (const artist of artists) {
      if (artist.id && !artistAvatars.has(artist.id)) {
        try {
          const response = await axios.get(`${API_URL}/artists/${artist.id}/avatar`, {
            responseType: 'blob'
          });
          const url = URL.createObjectURL(response.data);
          setArtistAvatars(prev => new Map(prev).set(artist.id, url));
        } catch (err) {
          console.error(`Failed to load avatar for artist ${artist.id}:`, err);
        }
      }
    }
  };

  const handleAddToPlaylist = (event: React.MouseEvent<HTMLElement>, track: TrackType) => {
    event.stopPropagation();
    if (!isAuthenticated) {
      navigate('/login', {
        state: {
          from: window.location.pathname,
          message: 'Чтобы добавить трек в плейлист, необходимо войти',
        },
      });
      return;
    }
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
                tracks: (p.tracks || []).filter((t: TrackType) => t.id !== selectedTrack.id)
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

  const isTrackInPlaylist = (playlist: any, trackId: string) => {
    return playlist.tracks?.some((track: TrackType) => track.id === trackId) || false;
  };

  const renderResults = () => {
    if (loading) {
      return <LoadingSpinner />;
    }

    if (error) {
      return <Alert severity="error">{error}</Alert>;
    }

    if (results.length === 0) {
      return (
        <Box sx={{ textAlign: 'center', py: 4 }}>
          <Typography variant="h6" color="text.secondary">
            Ничего не найдено
          </Typography>
        </Box>
      );
    }

    switch (contentType) {
      case 'track':
        const tracksWithCovers = results.map((track: TrackType) => ({
          ...track,
          coverUrl: albumCovers.get(track.album.id)
        }));
        return (
          <TrackList
            tracks={tracksWithCovers}
            onTrackClick={handleTrackClick}
            onAddToPlaylist={handleAddToPlaylist}
            showTrackNumber={false}
            showAlbumLink={true}
          />
        );
      case "album":
        return (
          <Grid container spacing={3}>
            {results.map((album: any) => (
              <Grid item xs={12} sm={6} md={3} key={album.id}>
                <AlbumCard
                  id={album.id}
                  title={album.title}
                  label={album.label}
                  release_date={album.release_date}
                  artists={album.artists}
                  genres={album.genres}
                />
              </Grid>
            ))}
          </Grid>
        );
      case "playlist":
        return (
          <Grid container spacing={3}>
            {results.map((playlist: any) => (
              <Grid item xs={12} sm={6} md={3} key={playlist.id}>
                <PlaylistCard
                  id={playlist.id}
                  name={playlist.name}
                  coverUrl={playlistCovers.get(playlist.id)}
                  isFavorite={playlist.is_favorite}
                  rating={playlist.rating || 0}
                  owner={playlist.owner}
                  onFavoriteChange={(isFavorite) => {
                    // Обновляем локальное состояние плейлиста
                    setResults(prev => prev.map(p => 
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
            ))}
          </Grid>
        );
      case "artist":
        return (
          <Grid container spacing={3}>
            {results.map((artist: any) => (
              <Grid item xs={12} sm={6} md={3} key={artist.id}>
                <ArtistCard
                  id={artist.id}
                  name={artist.name}
                  country={artist.country}
                  genre={artist.genre?.title || artist.genre || ''}
                  coverImage={artistAvatars.get(artist.id)}
                />
              </Grid>
            ))}
          </Grid>
        );
    }
  };

  return (
    <ErrorBoundary>
      <Container maxWidth="lg">
        <Box sx={{ mb: 4 }}>
          <Box sx={{ mb: 3 }}>
            <form onSubmit={handleSearch}>
              <Grid container spacing={2} alignItems="center">
                <Grid item xs={12} md={4}>
                  <TextField
                    fullWidth
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    placeholder="Search for tracks, albums, artists, or playlists..."
                    InputProps={{
                      startAdornment: (
                        <InputAdornment position="start">
                          <SearchIcon />
                        </InputAdornment>
                      ),
                    }}
                  />
                </Grid>
                
                <Grid item xs={12} md={2}>
                  <FormControl fullWidth>
                    <InputLabel>Content Type</InputLabel>
                    <Select
                      value={contentType}
                      label="Content Type"
                      onChange={(e) => setContentType(e.target.value as ContentType)}
                    >
                      <MenuItem value="track">Tracks</MenuItem>
                      <MenuItem value="album">Albums</MenuItem>
                      <MenuItem value="playlist">Playlists</MenuItem>
                      <MenuItem value="artist">Artists</MenuItem>
                    </Select>
                  </FormControl>
                </Grid>
                
                <Grid item xs={12} md={2}>
                  <TextField
                    fullWidth
                    label="Country"
                    value={country}
                    onChange={(e) => setCountry(e.target.value)}
                    placeholder="Enter country name"
                  />
                </Grid>
                
                <Grid item xs={12} md={2}>
                  <FormControl fullWidth>
                    <Autocomplete
                      options={genres}
                      getOptionLabel={(option) => option.title}
                      value={genre}
                      onChange={(_, newValue) => setGenre(newValue)}
                      loading={loadingGenres}
                      renderInput={(params) => (
                        <TextField
                          {...params}
                          label="Genre"
                          InputProps={{
                            ...params.InputProps,
                            endAdornment: (
                              <>
                                {loadingGenres ? <CircularProgress color="inherit" size={20} /> : null}
                                {params.InputProps.endAdornment}
                              </>
                            ),
                          }}
                        />
                      )}
                    />
                  </FormControl>
                </Grid>

                <Grid item xs={12} md={2}>
                  <Button 
                    type="submit" 
                    variant="contained" 
                    color="primary"
                    fullWidth
                    disabled={loading}
                  >
                    {loading ? <CircularProgress size={24} /> : "Search"}
                  </Button>
                </Grid>
              </Grid>
            </form>
          </Box>

          {renderResults()}

          {isAuthenticated && (
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
                    <CheckIcon sx={{ ml: 1, color: 'white' }} />
                  )}
                </MenuItem>
              ))}
            </Menu>
          )}
        </Box>
      </Container>
    </ErrorBoundary>
  );
};

export default Search; 