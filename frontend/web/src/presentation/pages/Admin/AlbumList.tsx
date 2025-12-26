import { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Paper,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Alert,
  CircularProgress,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Chip,
  OutlinedInput,
  // Stack,
  Switch,
  FormControlLabel,
  Grid,
} from '@mui/material';
import { Edit as EditIcon, Delete as DeleteIcon, Add as AddIcon, Remove as RemoveIcon } from '@mui/icons-material';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import { api, apiService } from '../../../presentation/services/api';

interface License {
  id: string;
  title: string;
  description: string;
  url: string;
}

interface Track {
  id?: string;
  name: string;
  genre_id: string;
  duration: number;
  explicit: boolean;
  license_id: string;
  track_number: number;
  audioFile?: File;
  hasAudio?: boolean;
  additionalArtists?: string[];
}

interface Album {
  id: string;
  title: string;
  label: string;
  license: License;
  release_date: string;
  artists: Artist[];
  genres: Genre[];
  coverFile?: File;
  coverUrl?: string;
  tracks: Track[];
}

interface Artist {
  id: string;
  name: string;
  description: string;
  country: string;
}

interface Genre {
  id: string;
  title: string;
}

const ITEM_HEIGHT = 48;
const ITEM_PADDING_TOP = 8;
const MenuProps = {
  PaperProps: {
    style: {
      maxHeight: ITEM_HEIGHT * 4.5 + ITEM_PADDING_TOP,
      width: 250,
    },
  },
};

// Добавляем функции конвертации
const secondsToMMSS = (seconds: number): string => {
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = seconds % 60;
  return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
};

const mmssToSeconds = (mmss: string): number => {
  const [minutes, seconds] = mmss.split(':').map(Number);
  return minutes * 60 + seconds;
};

const AlbumList = () => {
  const [albums, setAlbums] = useState<Album[]>([]);
  const [artists, setArtists] = useState<Artist[]>([]);
  const [genres, setGenres] = useState<Genre[]>([]);
  const [licenses, setLicenses] = useState<{ id: string; title: string; description: string; url: string; }[]>([]);
  const [open, setOpen] = useState(false);
  const [editingAlbum, setEditingAlbum] = useState<Album | null>(null);
  const [formData, setFormData] = useState({
    title: '',
    label: '',
    releaseDate: '',
    licenseId: '',
    selectedArtists: [] as string[],
    selectedGenres: [] as string[],
    coverFile: null as File | null,
    coverUrl: undefined as string | undefined,
    tracks: [] as Track[],
  });
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  const fetchAlbums = async () => {
    try {
      setLoading(true);
      const response = await api.getAlbums();
      const albumsData = Array.isArray(response) ? response : [];
      setAlbums(albumsData);
      setError(null);
    } catch (err) {
      setError('Ошибка при загрузке альбомов');
      console.error('Error fetching albums:', err);
      setAlbums([]);
    } finally {
      setLoading(false);
    }
  };

  const fetchArtists = async () => {
    try {
      const response = await api.getArtists();
      const artistsData = Array.isArray(response) ? response : [];
      setArtists(artistsData);
    } catch (err) {
      console.error('Error fetching artists:', err);
    }
  };

  const fetchGenres = async () => {
    try {
      const response = await api.getGenres();
      console.log('Genres response:', response);
      const genresData = Array.isArray(response) ? response : [];
      console.log('Processed genres data:', genresData);
      setGenres(genresData);
    } catch (err) {
      console.error('Error fetching genres:', err);
    }
  };

  const fetchLicenses = async () => {
    try {
      const response = await api.getLicenses();
      const licensesData = Array.isArray(response) ? response : [];
      setLicenses(licensesData as any);
    } catch (err) {
      console.error('Error fetching licenses:', err);
    }
  };

  useEffect(() => {
    fetchAlbums();
    fetchArtists();
    fetchGenres();
    fetchLicenses();
  }, []);

  const handleOpen = async (album?: Album) => {
    if (album) {
      setEditingAlbum(album);
      // Загружаем треки альбома
      try {
        const tracks = await apiService.get(`/albums/${album.id}/tracks`);
        console.log('Fetched tracks:', tracks);
        
        // Проверяем и валидируем данные треков
        const validatedTracks = tracks.map((track: any) => {
          // Проверяем, что genre_id существует и соответствует одному из доступных жанров
          const validGenreId = track.genre_id && genres.some(g => g.id === track.genre_id) 
            ? track.genre_id 
            : genres.length > 0 ? genres[0].id : '';

          // Получаем дополнительных артистов из TrackMeta
          const additionalArtists = track.artists
            ?.filter((artist: any) => !album.artists.some(a => a.id === artist.id))
            .map((artist: any) => artist.id) || [];

          return {
            id: track.id,
            name: track.name,
            genre_id: validGenreId,
            duration: track.duration || 0,
            explicit: track.explicit || false,
            license_id: track.license.id,
            track_number: track.track_number || 1,
            hasAudio: true,
            additionalArtists: additionalArtists
          };
        });

        setFormData({
          title: album.title,
          label: album.label,
          releaseDate: album.release_date,
          licenseId: album.license.id,
          selectedArtists: album.artists.map(artist => artist.id),
          selectedGenres: album.genres.map(genre => genre.id),
          coverFile: null,
          coverUrl: album.coverUrl,
          tracks: validatedTracks
        });
      } catch (err) {
        console.error('Error fetching album tracks:', err);
        setError('Ошибка при загрузке треков альбома');
      }
    } else {
      setEditingAlbum(null);
      setFormData({
        title: '',
        label: '',
        releaseDate: '',
        licenseId: '',
        selectedArtists: [],
        selectedGenres: [],
        coverFile: null,
        coverUrl: undefined,
        tracks: [],
      });
    }
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
    setEditingAlbum(null);
    setFormData({
      title: '',
      label: '',
      releaseDate: '',
      licenseId: '',
      selectedArtists: [],
      selectedGenres: [],
      coverFile: null,
      coverUrl: undefined,
      tracks: [],
    });
    setError(null);
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSelectChange = (e: any) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleArtistsChange = (event: any) => {
    const {
      target: { value },
    } = event;
    setFormData({
      ...formData,
      selectedArtists: typeof value === 'string' ? value.split(',') : value,
    });
  };

  const handleGenresChange = (event: any) => {
    const {
      target: { value },
    } = event;
    setFormData({
      ...formData,
      selectedGenres: typeof value === 'string' ? value.split(',') : value,
    });
  };

  const handleDateChange = (date: Date | null) => {
    if (date) {
      const isoDate = date.toISOString();
      setFormData({
        ...formData,
        releaseDate: isoDate,
      });
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    // Проверяем длительность всех треков
    const invalidTrack = formData.tracks.find(track => track.duration <= 0);
    if (invalidTrack) {
      setError('Длительность трека должна быть больше 0');
      return;
    }

    try {
      // 1. Создаем альбом
      const albumData = {
        title: formData.title.trim(),
        label: formData.label?.trim() || '',
        release_date: formData.releaseDate,
        license_id: formData.licenseId.trim(),
      };

      let albumId;
      if (editingAlbum) {
        await apiService.put(`/albums/${editingAlbum.id}`, albumData);
        albumId = editingAlbum.id;

        // Для существующего альбома находим только новые жанры и артистов
        const existingGenreIds = editingAlbum.genres.map(g => g.id);
        const existingArtistIds = editingAlbum.artists.map(a => a.id);

        const newGenreIds = formData.selectedGenres.filter(id => !existingGenreIds.includes(id));
        const newArtistIds = formData.selectedArtists.filter(id => !existingArtistIds.includes(id));

        // Назначаем только новые жанры
        if (newGenreIds.length > 0) {
          for (const genreId of newGenreIds) {
            await apiService.post(`/albums/${albumId}/genres/${genreId}`);
          }
        }

        // Назначаем только новых артистов
        if (newArtistIds.length > 0) {
          for (const artistId of newArtistIds) {
            await apiService.post(`/artists/${artistId}/albums/${albumId}`);
          }
        }
      } else {
        const response = await apiService.post('/albums', albumData);
        albumId = response.id;

        // Для нового альбома назначаем все выбранные жанры и артистов
        if (formData.selectedGenres.length > 0) {
          for (const genreId of formData.selectedGenres) {
            await apiService.post(`/albums/${albumId}/genres/${genreId}`);
          }
        }

        if (formData.selectedArtists.length > 0) {
          for (const artistId of formData.selectedArtists) {
            await apiService.post(`/artists/${artistId}/albums/${albumId}`);
          }
        }
      }

      // 2. Загружаем обложку, если она есть
      if (formData.coverFile) {
        const formDataCover = new FormData();
        formDataCover.append('cover', formData.coverFile, formData.coverFile.name);
        console.log('Uploading cover:', {
          fileName: formData.coverFile.name,
          fileType: formData.coverFile.type,
          fileSize: formData.coverFile.size,
        });
        try {
          await apiService.post(`/albums/${albumId}/cover`, formDataCover, {
            headers: {
              'Content-Type': 'multipart/form-data',
            },
          });
        } catch (err) {
          console.error('Error uploading cover:', err);
          throw err;
        }
      }

      // 3. Создаем и загружаем треки
      for (const track of formData.tracks) {
        // Создаем трек только если у него нет id (новый трек)
        var trackID = track.id;
        if (!track.id) {
          const trackData = {
            name: track.name,
            genre_id: track.genre_id,
            duration: track.duration,
            explicit: track.explicit,
            license_id: track.license_id,
            album_id: albumId,
            track_number: track.track_number,
          };

          const trackResponse = await apiService.post('/tracks', trackData);
          trackID = trackResponse.id;
        }

        // Связываем трек со всеми артистами (и основными артистами альбома, и дополнительными)
        const allArtists = [...formData.selectedArtists];
        if (track.additionalArtists) {
          allArtists.push(...track.additionalArtists);
        }

        // Удаляем дубликаты
        const uniqueArtists = [...new Set(allArtists)];

        for (const artistId of uniqueArtists) {
          try {
            await apiService.post(`/artists/${artistId}/tracks/${trackID}`);
          } catch (err) {
            console.error(`Error adding artist ${artistId} to track ${trackID}:`, err);
            // Продолжаем выполнение даже если не удалось добавить артиста
          }
        }

        if (track.audioFile) {
          const formDataAudio = new FormData();
          formDataAudio.append('audio', track.audioFile, track.audioFile.name);
          console.log('Uploading audio:', {
            fileName: track.audioFile.name,
            fileType: track.audioFile.type,
            fileSize: track.audioFile.size,
            trackID: trackID,
          });
          try {
            await apiService.post(`/tracks/${trackID}/audio`, formDataAudio, {
              headers: {
                'Content-Type': 'multipart/form-data',
              },
            });
          } catch (err) {
            console.error('Error uploading audio:', err);
            throw err;
          }
        }
      }

      // Закрываем форму и обновляем список
      await fetchAlbums();
      setOpen(false);
      setEditingAlbum(null);
      setFormData({
        title: '',
        label: '',
        releaseDate: '',
        licenseId: '',
        selectedArtists: [],
        selectedGenres: [],
        coverFile: null,
        coverUrl: undefined,
        tracks: [],
      });
    } catch (err) {
      console.error('Error in handleSubmit:', err);
      setError(err instanceof Error ? err.message : 'Ошибка при сохранении альбома');
    }
  };

  const handleCoverChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setFormData({ ...formData, coverFile: e.target.files[0] });
    }
  };

  const handleAddTrack = () => {
    // При добавлении нового трека устанавливаем первый доступный жанр по умолчанию
    const defaultGenreId = genres.length > 0 ? genres[0].id : '';
    
    setFormData({
      ...formData,
      tracks: [
        ...formData.tracks,
        {
          name: '',
          genre_id: defaultGenreId,
          duration: 0,
          explicit: false,
          license_id: licenses[0].id,
          track_number: formData.tracks.length + 1,
          hasAudio: false,
          additionalArtists: [],
        },
      ],
    });
  };

  const handleRemoveTrack = async (index: number) => {
    const track = formData.tracks[index];
    
    // Если у трека есть id, значит он существует на сервере и его нужно удалить
    if (track.id) {
      try {
        await apiService.delete(`/tracks/${track.id}`);
      } catch (err) {
        console.error('Error deleting track:', err);
        setError('Ошибка при удалении трека');
        return;
      }
    }

    // Удаляем трек из локального состояния
    const newTracks = [...formData.tracks];
    newTracks.splice(index, 1);
    // Обновляем номера треков
    newTracks.forEach((track, idx) => {
      track.track_number = idx + 1;
    });
    setFormData({ ...formData, tracks: newTracks });
  };

  const handleTrackChange = (index: number, field: keyof Track, value: any) => {
    const newTracks = [...formData.tracks];
    if (field === 'duration') {
      // Если значение в формате MM:SS, конвертируем в секунды
      if (typeof value === 'string' && value.includes(':')) {
        const seconds = mmssToSeconds(value);
        if (seconds <= 0) {
          setError('Длительность трека должна быть больше 0');
          return;
        }
        newTracks[index] = { ...newTracks[index], [field]: seconds };
      } else {
        if (value <= 0) {
          setError('Длительность трека должна быть больше 0');
          return;
        }
        newTracks[index] = { ...newTracks[index], [field]: value };
      }
    } else {
      newTracks[index] = { ...newTracks[index], [field]: value };
    }
    setFormData({ ...formData, tracks: newTracks });
  };

  const handleTrackAudioChange = (index: number, e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const newTracks = [...formData.tracks];
      newTracks[index] = { ...newTracks[index], audioFile: e.target.files[0] };
      setFormData({ ...formData, tracks: newTracks });
    }
  };

  const handleDeleteTrackAudio = async (trackId: string) => {
    try {
      await apiService.delete(`/tracks/${trackId}/audio`);
      // Обновляем список треков после удаления аудио
      if (editingAlbum) {
        const tracks = await apiService.get(`/albums/${editingAlbum.id}/tracks`);
        setFormData(prev => ({ ...prev, tracks }));
      }
    } catch (err) {
      console.error('Error deleting track audio:', err);
      setError('Ошибка при удалении аудиофайла');
    }
  };

  const handleDelete = async (id: string) => {
    if (window.confirm('Вы уверены, что хотите удалить этот альбом?')) {
      try {
        await api.deleteAlbum(id);
        fetchAlbums();
      } catch (err) {
        setError('Ошибка при удалении альбома');
        console.error('Error deleting album:', err);
      }
    }
  };

  return (
    <Box sx={{ p: 4, maxWidth: 1200, mx: 'auto' }}>
      <Paper sx={{ p: 4 }} elevation={3}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3, gap: 2 }}>
          <Typography variant="h5" fontWeight={700}>
            Управление альбомами
          </Typography>
          <Button variant="contained" color="primary" onClick={() => handleOpen()}>
            Добавить альбом
          </Button>
        </Box>

        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        {loading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
            <CircularProgress />
          </Box>
        ) : (
          <TableContainer>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>ID</TableCell>
                  <TableCell>Название</TableCell>
                  <TableCell>Исполнители</TableCell>
                  <TableCell>Жанры</TableCell>
                  <TableCell>Лицензия</TableCell>
                  <TableCell>Дата выпуска</TableCell>
                  <TableCell align="right">Действия</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {albums.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={7} align="center">
                      Нет доступных альбомов
                    </TableCell>
                  </TableRow>
                ) : (
                  albums.map((album) => (
                    <TableRow key={album.id}>
                      <TableCell>{album.id}</TableCell>
                      <TableCell>{album.title}</TableCell>
                      <TableCell>
                        {album.artists?.map(artist => artist.name).join(', ') || 'Нет исполнителей'}
                      </TableCell>
                      <TableCell>
                        {album.genres?.map(genre => genre.title).join(', ') || 'Нет жанров'}
                      </TableCell>
                      <TableCell>
                        {album.license?.title || 'Нет лицензии'}
                      </TableCell>
                      <TableCell>{new Date(album.release_date).toLocaleDateString()}</TableCell>
                      <TableCell align="right">
                        <IconButton onClick={() => handleOpen(album)} color="primary">
                          <EditIcon />
                        </IconButton>
                        <IconButton onClick={() => handleDelete(album.id)} color="error">
                          <DeleteIcon />
                        </IconButton>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </TableContainer>
        )}
      </Paper>

      <Dialog open={open} onClose={handleClose} maxWidth="md" fullWidth>
        <DialogTitle>
          {editingAlbum ? 'Редактировать альбом' : 'Добавить альбом'}
        </DialogTitle>
        <form onSubmit={handleSubmit}>
          <DialogContent>
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <TextField
                  autoFocus
                  margin="dense"
                  label="Название"
                  fullWidth
                  value={formData.title}
                  onChange={handleChange}
                  name="title"
                  required
                />
              </Grid>
              <Grid item xs={12}>
                <TextField
                  margin="dense"
                  label="Лейбл"
                  fullWidth
                  value={formData.label}
                  onChange={handleChange}
                  name="label"
                />
              </Grid>
              <Grid item xs={12}>
                <FormControl fullWidth>
                  <InputLabel>Артисты</InputLabel>
                  <Select
                    multiple
                    value={formData.selectedArtists}
                    onChange={handleArtistsChange}
                    input={<OutlinedInput label="Артисты" />}
                    renderValue={(selected) => (
                      <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                        {selected.map((value) => (
                          <Chip
                            key={value}
                            label={artists.find(artist => artist.id === value)?.name || value}
                          />
                        ))}
                      </Box>
                    )}
                    MenuProps={MenuProps}
                  >
                    {artists.map((artist) => (
                      <MenuItem key={artist.id} value={artist.id}>
                        {artist.name}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Grid>
              <Grid item xs={12}>
                <FormControl fullWidth>
                  <InputLabel>Жанры</InputLabel>
                  <Select
                    multiple
                    value={formData.selectedGenres}
                    onChange={handleGenresChange}
                    input={<OutlinedInput label="Жанры" />}
                    renderValue={(selected) => (
                      <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                        {selected.map((value) => (
                          <Chip
                            key={value}
                            label={genres.find(genre => genre.id === value)?.title || value}
                          />
                        ))}
                      </Box>
                    )}
                    MenuProps={MenuProps}
                  >
                    {genres.map((genre) => (
                      <MenuItem key={genre.id} value={genre.id}>
                        {genre.title}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Grid>
              <Grid item xs={12}>
                <FormControl fullWidth>
                  <InputLabel>Лицензия</InputLabel>
                  <Select
                    value={formData.licenseId}
                    onChange={handleSelectChange}
                    name="licenseId"
                    label="Лицензия"
                    required
                  >
                    {licenses.map((license) => (
                      <MenuItem key={license.id} value={license.id}>
                        {license.title}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Grid>
              <Grid item xs={12}>
                <DatePicker
                  label="Дата выпуска"
                  value={formData.releaseDate ? new Date(formData.releaseDate) : null}
                  onChange={handleDateChange}
                  slotProps={{
                    textField: {
                      fullWidth: true,
                      margin: 'dense',
                      required: true,
                    }
                  }}
                />
              </Grid>
              <Grid item xs={12}>
                <input
                  accept="image/*"
                  type="file"
                  id="cover-upload"
                  onChange={handleCoverChange}
                  style={{ display: 'none' }}
                />
                <label htmlFor="cover-upload">
                  <Button variant="outlined" component="span">
                    {formData.coverUrl ? 'Изменить обложку' : 'Загрузить обложку'}
                  </Button>
                </label>
                {formData.coverFile && (
                  <Typography variant="body2" sx={{ mt: 1 }}>
                    Выбрана новая обложка: {formData.coverFile.name}
                  </Typography>
                )}
                {formData.coverUrl && !formData.coverFile && (
                  <Box sx={{ mt: 2, maxWidth: 200 }}>
                    <img
                      src={formData.coverUrl}
                      alt="Album cover"
                      style={{ width: '100%', height: 'auto' }}
                    />
                  </Box>
                )}
              </Grid>
              <Grid item xs={12}>
                <Typography variant="h6" sx={{ mb: 2 }}>
                  Треки
                </Typography>
                {formData.tracks.map((track, index) => (
                  <Paper key={track.id || index} sx={{ p: 2, mb: 2 }}>
                    <Grid container spacing={2} alignItems="center">
                      <Grid item xs={12} sm={6}>
                        <TextField
                          label="Название трека"
                          fullWidth
                          value={track.name}
                          onChange={(e) => handleTrackChange(index, 'name', e.target.value)}
                          required
                        />
                      </Grid>
                      <Grid item xs={12} sm={6}>
                        <FormControl fullWidth>
                          <InputLabel>Жанр</InputLabel>
                          <Select
                            value={track.genre_id}
                            onChange={(e) => handleTrackChange(index, 'genre_id', e.target.value)}
                            label="Жанр"
                            required
                          >
                            {genres.map((genre) => (
                              <MenuItem key={genre.id} value={genre.id}>
                                {genre.title}
                              </MenuItem>
                            ))}
                          </Select>
                        </FormControl>
                      </Grid>
                      <Grid item xs={12} sm={6}>
                        <TextField
                          label="Длительность (ММ:СС)"
                          fullWidth
                          value={secondsToMMSS(track.duration)}
                          onChange={(e) => handleTrackChange(index, 'duration', e.target.value)}
                          placeholder="3:45"
                          required
                          inputProps={{
                            pattern: '^[0-9]+:[0-5][0-9]$',
                            title: 'Формат: минуты:секунды (например, 3:45)'
                          }}
                        />
                      </Grid>
                      <Grid item xs={12} sm={6}>
                        <FormControl fullWidth>
                          <InputLabel>Лицензия</InputLabel>
                          <Select
                            value={track.license_id}
                            onChange={(e) => handleTrackChange(index, 'license_id', e.target.value)}
                            label="Лицензия"
                            required
                          >
                            {licenses.map((license) => (
                              <MenuItem key={license.id} value={license.id}>
                                {license.title}
                              </MenuItem>
                            ))}
                          </Select>
                        </FormControl>
                      </Grid>
                      <Grid item xs={12} sm={6}>
                        <FormControlLabel
                          control={
                            <Switch
                              checked={track.explicit}
                              onChange={(e) => handleTrackChange(index, 'explicit', e.target.checked)}
                            />
                          }
                          label="Explicit"
                        />
                      </Grid>
                      <Grid item xs={12}>
                        <FormControl fullWidth>
                          <InputLabel>Дополнительные артисты</InputLabel>
                          <Select
                            multiple
                            value={track.additionalArtists || []}
                            onChange={(e) => handleTrackChange(index, 'additionalArtists', e.target.value)}
                            input={<OutlinedInput label="Дополнительные артисты" />}
                            renderValue={(selected) => (
                              <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                                {selected.map((value) => (
                                  <Chip
                                    key={value}
                                    label={artists.find(artist => artist.id === value)?.name || value}
                                  />
                                ))}
                              </Box>
                            )}
                            MenuProps={MenuProps}
                          >
                            {artists
                              .filter(artist => !formData.selectedArtists.includes(artist.id)) // Исключаем артистов альбома
                              .map((artist) => (
                                <MenuItem key={artist.id} value={artist.id}>
                                  {artist.name}
                                </MenuItem>
                              ))}
                          </Select>
                        </FormControl>
                      </Grid>
                      <Grid item xs={12} sm={6}>
                        <input
                          accept="audio/*"
                          type="file"
                          id={`track-audio-${index}`}
                          onChange={(e) => handleTrackAudioChange(index, e)}
                          style={{ display: 'none' }}
                        />
                        <label htmlFor={`track-audio-${index}`}>
                          <Button variant="outlined" component="span">
                            {track.audioFile ? 'Изменить аудио' : 'Загрузить аудио'}
                          </Button>
                        </label>
                        {track.audioFile && (
                          <Typography variant="body2" sx={{ mt: 1 }}>
                            Выбран новый файл: {track.audioFile.name}
                          </Typography>
                        )}
                        {track.id && !track.audioFile && (
                          <Box sx={{ mt: 1 }}>
                            <Typography variant="body2" sx={{ mb: 1 }}>
                              {track.hasAudio ? 'Аудиофайл загружен' : 'Аудиофайл не загружен'}
                            </Typography>
                            {track.hasAudio && (
                              <Button
                                variant="outlined"
                                color="error"
                                size="small"
                                onClick={() => handleDeleteTrackAudio(track.id!)}
                              >
                                Удалить аудио
                              </Button>
                            )}
                          </Box>
                        )}
                      </Grid>
                      <Grid item xs={12}>
                        <Button
                          variant="outlined"
                          color="error"
                          startIcon={<RemoveIcon />}
                          onClick={() => handleRemoveTrack(index)}
                        >
                          Удалить трек
                        </Button>
                      </Grid>
                    </Grid>
                  </Paper>
                ))}
                <Button
                  variant="outlined"
                  startIcon={<AddIcon />}
                  onClick={handleAddTrack}
                  sx={{ mt: 2 }}
                >
                  Добавить трек
                </Button>
              </Grid>
            </Grid>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose}>Отмена</Button>
            <Button type="submit" variant="contained" color="primary">
              {editingAlbum ? 'Сохранить' : 'Добавить'}
            </Button>
          </DialogActions>
        </form>
      </Dialog>
    </Box>
  );
};

export default AlbumList; 