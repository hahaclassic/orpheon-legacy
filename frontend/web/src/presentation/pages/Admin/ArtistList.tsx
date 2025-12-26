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
  Grid,
} from '@mui/material';
import { Edit as EditIcon, Delete as DeleteIcon } from '@mui/icons-material';
import { api, apiService } from '../../../presentation/services/api';

interface Artist {
  id: string;
  name: string;
  description: string;
  country: string;
  avatarUrl?: string;
}

interface ArtistResponse {
  id: string;
  name: string;
  description: string;
  country: string;
  avatar_url?: string;
}

const ArtistList = () => {
  const [artists, setArtists] = useState<Artist[]>([]);
  const [open, setOpen] = useState(false);
  const [editingArtist, setEditingArtist] = useState<Artist | null>(null);
  const [artistAvatars, setArtistAvatars] = useState<{ [key: string]: string }>({});
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    country: '',
    avatarFile: null as File | null,
    avatarUrl: undefined as string | undefined,
  });
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  const fetchArtistAvatar = async (artistId: string) => {
    try {
      const response = await apiService.get(`/artists/${artistId}/avatar`, {
        responseType: 'blob'
      });
      const url = URL.createObjectURL(response);
      setArtistAvatars(prev => ({ ...prev, [artistId]: url }));
    } catch (err) {
      console.error('Error fetching artist avatar:', err);
    }
  };

  const fetchArtists = async () => {
    try {
      setLoading(true);
      const response = await api.getArtists();
      console.log('Raw API response:', response);
      const artistsData = Array.isArray(response) ? response.map((artist: any) => ({
        ...artist,
        avatarUrl: artist.avatar_url ? `/artists/${artist.id}/avatar` : undefined
      })) : [];
      console.log('Processed artists data:', artistsData);
      setArtists(artistsData);
      
      // Загружаем аватары для всех артистов
      artistsData.forEach(artist => {
        if (artist.avatar_url) {
          fetchArtistAvatar(artist.id);
        }
      });
      
      setError(null);
    } catch (err) {
      setError('Ошибка при загрузке артистов');
      console.error('Error fetching artists:', err);
      setArtists([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchArtists();
  }, []);

  const handleOpen = (artist?: Artist) => {
    if (artist) {
      setEditingArtist(artist);
      setFormData({
        name: artist.name,
        description: artist.description,
        country: artist.country,
        avatarFile: null,
        avatarUrl: artist.avatarUrl,
      });
    } else {
      setEditingArtist(null);
      setFormData({
        name: '',
        description: '',
        country: '',
        avatarFile: null,
        avatarUrl: undefined,
      });
    }
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
    setEditingArtist(null);
    setFormData({
      name: '',
      description: '',
      country: '',
      avatarFile: null,
      avatarUrl: undefined,
    });
    setError(null);
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleAvatarChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      
      // Проверка размера файла (10MB)
      if (file.size > 10 * 1024 * 1024) {
        setError('Размер файла превышает 10MB');
        return;
      }

      // Проверка формата файла
      if (file.type !== 'image/jpeg' && file.type !== 'image/png') {
        setError('Допустимы только файлы JPEG и PNG');
        return;
      }

      setFormData({ ...formData, avatarFile: file });
      setError(null);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const trimmedData = {
        name: formData.name.trim(),
        description: formData.description.trim(),
        country: formData.country.trim(),
      };

      let artistId;
      if (editingArtist) {
        await api.updateArtist(editingArtist.id, trimmedData);
        artistId = editingArtist.id;
      } else {
        const response = await api.createArtist(trimmedData);
        artistId = response.id;
      }

      // Загружаем аватар, если он есть
      if (formData.avatarFile) {
        const formDataAvatar = new FormData();
        formDataAvatar.append('avatar', formData.avatarFile);
        try {
          await apiService.post(`/artists/${artistId}/avatar`, formDataAvatar, {
            headers: {
              'Content-Type': 'multipart/form-data',
            },
          });
        } catch (err) {
          console.error('Error uploading avatar:', err);
          throw err;
        }
      }

      handleClose();
      fetchArtists();
    } catch (err) {
      setError('Ошибка при сохранении артиста');
      console.error('Error saving artist:', err);
    }
  };

  const handleDelete = async (id: string) => {
    if (window.confirm('Вы уверены, что хотите удалить этого артиста?')) {
      try {
        await api.deleteArtist(id);
        fetchArtists();
      } catch (err) {
        setError('Ошибка при удалении артиста');
        console.error('Error deleting artist:', err);
      }
    }
  };

  return (
    <Box sx={{ p: 4, maxWidth: 1200, mx: 'auto' }}>
      <Paper sx={{ p: 4 }} elevation={3}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3, gap: 2 }}>
          <Typography variant="h5" fontWeight={700}>
            Управление артистами
          </Typography>
          <Button variant="contained" color="primary" onClick={() => handleOpen()}>
            Добавить артиста
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
                  <TableCell>Аватар</TableCell>
                  <TableCell>Имя</TableCell>
                  <TableCell>Страна</TableCell>
                  <TableCell>Описание</TableCell>
                  <TableCell align="right">Действия</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {artists.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} align="center">
                      Нет доступных артистов
                    </TableCell>
                  </TableRow>
                ) : (
                  artists.map((artist) => (
                    <TableRow key={artist.id}>
                      <TableCell>{artist.id}</TableCell>
                      <TableCell>
                        {artist.avatarUrl ? (
                          artistAvatars[artist.id] ? (
                            <Box sx={{ width: 50, height: 50 }}>
                              <img
                                src={artistAvatars[artist.id]}
                                alt={`${artist.name} avatar`}
                                style={{ width: '100%', height: '100%', borderRadius: '50%', objectFit: 'cover' }}
                              />
                            </Box>
                          ) : (
                            <Box sx={{ width: 50, height: 50, bgcolor: 'grey.200', borderRadius: '50%' }} />
                          )
                        ) : (
                          <Box sx={{ width: 50, height: 50, bgcolor: 'grey.200', borderRadius: '50%' }} />
                        )}
                      </TableCell>
                      <TableCell>{artist.name}</TableCell>
                      <TableCell>{artist.country}</TableCell>
                      <TableCell>{artist.description}</TableCell>
                      <TableCell align="right">
                        <IconButton onClick={() => handleOpen(artist)} color="primary">
                          <EditIcon />
                        </IconButton>
                        <IconButton onClick={() => handleDelete(artist.id)} color="error">
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

      <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
        <DialogTitle>
          {editingArtist ? 'Редактировать артиста' : 'Добавить артиста'}
        </DialogTitle>
        <form onSubmit={handleSubmit}>
          <DialogContent>
            <TextField
              name="name"
              label="Имя"
              value={formData.name}
              onChange={handleChange}
              fullWidth
              margin="normal"
              required
            />
            <TextField
              name="country"
              label="Страна"
              value={formData.country}
              onChange={handleChange}
              fullWidth
              margin="normal"
              required
            />
            <TextField
              name="description"
              label="Описание"
              value={formData.description}
              onChange={handleChange}
              fullWidth
              margin="normal"
              multiline
              rows={4}
            />
            <Grid item xs={12}>
              <input
                accept="image/jpeg,image/png"
                type="file"
                id="avatar-upload"
                onChange={handleAvatarChange}
                style={{ display: 'none' }}
              />
              <label htmlFor="avatar-upload">
                <Button variant="outlined" component="span">
                  {formData.avatarUrl ? 'Изменить аватар' : 'Загрузить аватар'}
                </Button>
              </label>
              {formData.avatarFile && (
                <Typography variant="body2" sx={{ mt: 1 }}>
                  Выбран новый аватар: {formData.avatarFile.name}
                </Typography>
              )}
              {formData.avatarUrl && !formData.avatarFile && (
                <Box sx={{ mt: 2, maxWidth: 200 }}>
                  <img
                    src={formData.avatarUrl}
                    alt="Artist avatar"
                    style={{ width: '100%', height: 'auto', borderRadius: '50%' }}
                  />
                </Box>
              )}
            </Grid>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose}>Отмена</Button>
            <Button type="submit" variant="contained" color="primary">
              {editingArtist ? 'Сохранить' : 'Добавить'}
            </Button>
          </DialogActions>
        </form>
      </Dialog>
    </Box>
  );
};

export default ArtistList; 