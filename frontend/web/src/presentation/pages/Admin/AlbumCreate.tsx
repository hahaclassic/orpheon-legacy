import { useState, useEffect } from 'react';
import {
  Box,
  Paper,
  Typography,
  TextField,
  Button,
  Stack,
  MenuItem,
  InputLabel,
  Select,
  FormControl,
  OutlinedInput,
  Chip,
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
  Alert,
} from '@mui/material';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import { Edit as EditIcon, Delete as DeleteIcon } from '@mui/icons-material';
import { apiService } from '../../../presentation/services/api';

interface Artist {
  id: number;
  name: string;
}

interface Genre {
  id: number;
  name: string;
}

interface License {
  id: number;
  name: string;
}

interface Album {
  id: number;
  title: string;
  label?: string;
  artists: Artist[];
  license: License;
  releaseDate: string;
  coverUrl?: string;
  genres: Genre[];
}

const AlbumCreate = () => {
  const [albums, setAlbums] = useState<Album[]>([]);
  const [artists, setArtists] = useState<Artist[]>([]);
  const [genres, setGenres] = useState<Genre[]>([]);
  const [licenses, setLicenses] = useState<License[]>([]);
  const [open, setOpen] = useState(false);
  const [editingAlbum, setEditingAlbum] = useState<Album | null>(null);
  const [form, setForm] = useState({
    title: '',
    label: '',
    artists: [] as number[],
    license: '',
    releaseDate: null as Date | null,
    cover: null as File | null,
    genres: [] as number[],
  });
  const [coverPreview, setCoverPreview] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const fetchData = async () => {
    try {
      const [albumsRes, artistsRes, genresRes, licensesRes] = await Promise.all([
        apiService.get('/albums'),
        apiService.get('/artists'),
        apiService.get('/genres'),
        apiService.get('/licenses'),
      ]);
      
      // Ensure albums is an array
      const albumsData = Array.isArray(albumsRes) ? albumsRes : [];
      setAlbums(albumsData);
      
      // Ensure other data is also arrays
      setArtists(Array.isArray(artistsRes) ? artistsRes : []);
      setGenres(Array.isArray(genresRes) ? genresRes : []);
      setLicenses(Array.isArray(licensesRes) ? licensesRes : []);
    } catch (err) {
      setError('Ошибка при загрузке данных');
      console.error('Error fetching data:', err);
      // Initialize empty arrays in case of error
      setAlbums([]);
      setArtists([]);
      setGenres([]);
      setLicenses([]);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const handleOpen = (album?: Album) => {
    if (album) {
      setEditingAlbum(album);
      setForm({
        title: album.title,
        label: album.label || '',
        artists: album.artists.map(a => a.id),
        license: album.license.id.toString(),
        releaseDate: new Date(album.releaseDate),
        cover: null,
        genres: album.genres.map(g => g.id),
      });
      setCoverPreview(album.coverUrl || null);
    } else {
      setEditingAlbum(null);
      setForm({
        title: '',
        label: '',
        artists: [],
        license: '',
        releaseDate: null,
        cover: null,
        genres: [],
      });
      setCoverPreview(null);
    }
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
    setEditingAlbum(null);
    setForm({
      title: '',
      label: '',
      artists: [],
      license: '',
      releaseDate: null,
      cover: null,
      genres: [],
    });
    setCoverPreview(null);
    setError(null);
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setForm((prev) => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSelectChange = (name: string, value: any) => {
    setForm((prev) => ({ ...prev, [name]: value }));
  };

  const handleCoverChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0] || null;
    setForm((prev) => ({ ...prev, cover: file }));
    if (file) {
      setCoverPreview(URL.createObjectURL(file));
    } else {
      setCoverPreview(null);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const formDataToSend = new FormData();
      formDataToSend.append('title', form.title.trim());
      formDataToSend.append('label', form.label.trim());
      form.artists.forEach((artistId: number) => {
        formDataToSend.append('artists', artistId.toString());
      });
      formDataToSend.append('license', form.license.toString());
      if (form.releaseDate) {
        formDataToSend.append('releaseDate', form.releaseDate.toISOString());
      }
      form.genres.forEach((genreId: number) => {
        formDataToSend.append('genres', genreId.toString());
      });
      if (coverPreview) {
        formDataToSend.append('cover', coverPreview);
      }

      if (editingAlbum) {
        await apiService.put(`/albums/${editingAlbum.id}`, formDataToSend);
      } else {
        await apiService.post('/albums', formDataToSend);
      }
      handleClose();
      fetchData();
    } catch (err) {
      setError('Ошибка при сохранении альбома');
      console.error('Error saving album:', err);
    }
  };

  const handleDelete = async (id: number) => {
    if (window.confirm('Вы уверены, что хотите удалить этот альбом?')) {
      try {
        await apiService.delete(`/albums/${id}`);
        fetchData();
      } catch (err) {
        setError('Ошибка при удалении альбома');
        console.error('Error deleting album:', err);
      }
    }
  };

  return (
    <Box sx={{ p: 4, maxWidth: 1200, mx: 'auto' }}>
      <Paper sx={{ p: 4 }} elevation={3}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
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

        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Обложка</TableCell>
                <TableCell>Название</TableCell>
                <TableCell>Лейбл</TableCell>
                <TableCell>Артисты</TableCell>
                <TableCell>Жанры</TableCell>
                <TableCell>Дата релиза</TableCell>
                <TableCell align="right">Действия</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {albums.map((album) => (
                <TableRow key={album.id}>
                  <TableCell>
                    {album.coverUrl && (
                      <img
                        src={album.coverUrl}
                        alt={album.title}
                        style={{ width: 50, height: 50, objectFit: 'cover', borderRadius: 4 }}
                      />
                    )}
                  </TableCell>
                  <TableCell>{album.title}</TableCell>
                  <TableCell>{album.label}</TableCell>
                  <TableCell>
                    {album.artists.map(artist => artist.name).join(', ')}
                  </TableCell>
                  <TableCell>
                    {album.genres.map(genre => genre.name).join(', ')}
                  </TableCell>
                  <TableCell>
                    {new Date(album.releaseDate).toLocaleDateString()}
                  </TableCell>
                  <TableCell align="right">
                    <IconButton onClick={() => handleOpen(album)} color="primary">
                      <EditIcon />
                    </IconButton>
                    <IconButton onClick={() => handleDelete(album.id)} color="error">
                      <DeleteIcon />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>

      <Dialog open={open} onClose={handleClose} maxWidth="md" fullWidth>
        <DialogTitle>
          {editingAlbum ? 'Редактировать альбом' : 'Добавить альбом'}
        </DialogTitle>
        <form onSubmit={handleSubmit}>
          <DialogContent>
            <Stack spacing={3}>
              <TextField
                autoFocus
                label="Название альбома"
                name="title"
                value={form.title}
                onChange={handleChange}
                required
                fullWidth
              />
              <TextField
                label="Лейбл (опционально)"
                name="label"
                value={form.label}
                onChange={handleChange}
                fullWidth
              />
              <FormControl fullWidth>
                <InputLabel>Артисты</InputLabel>
                <Select
                  multiple
                  value={form.artists}
                  onChange={e => handleSelectChange('artists', e.target.value)}
                  input={<OutlinedInput label="Артисты" />}
                  renderValue={(selected) => (
                    <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                      {(selected as number[]).map((id) => {
                        const artist = artists.find(a => a.id === id);
                        return artist ? <Chip key={id} label={artist.name} /> : null;
                      })}
                    </Box>
                  )}
                >
                  {artists.map((artist) => (
                    <MenuItem key={artist.id} value={artist.id}>
                      {artist.name}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
              <FormControl fullWidth required>
                <InputLabel>Лицензия</InputLabel>
                <Select
                  value={form.license}
                  onChange={e => handleSelectChange('license', e.target.value)}
                  label="Лицензия"
                >
                  {licenses.map((license) => (
                    <MenuItem key={license.id} value={license.id}>{license.name}</MenuItem>
                  ))}
                </Select>
              </FormControl>
              <DatePicker
                label="Дата релиза"
                value={form.releaseDate}
                onChange={date => handleSelectChange('releaseDate', date)}
                slotProps={{ textField: { fullWidth: true } }}
              />
              <FormControl fullWidth>
                <InputLabel>Жанры</InputLabel>
                <Select
                  multiple
                  value={form.genres}
                  onChange={e => handleSelectChange('genres', e.target.value)}
                  input={<OutlinedInput label="Жанры" />}
                  renderValue={(selected) => (
                    <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                      {(selected as number[]).map((id) => {
                        const genre = genres.find(g => g.id === id);
                        return genre ? <Chip key={id} label={genre.name} /> : null;
                      })}
                    </Box>
                  )}
                >
                  {genres.map((genre) => (
                    <MenuItem key={genre.id} value={genre.id}>
                      {genre.name}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
              <Box>
                <Button variant="outlined" component="label">
                  Загрузить обложку
                  <input type="file" accept="image/*" hidden onChange={handleCoverChange} />
                </Button>
                {coverPreview && (
                  <Box mt={2}>
                    <img src={coverPreview} alt="Обложка" style={{ maxWidth: 200, borderRadius: 8 }} />
                  </Box>
                )}
              </Box>
            </Stack>
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

export default AlbumCreate; 