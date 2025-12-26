import { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Paper,
  TextField,
  Button,
  Stack,
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
import { Edit as EditIcon, Delete as DeleteIcon } from '@mui/icons-material';
import { apiService } from '../../../presentation/services/api';

interface Artist {
  id: number;
  name: string;
  description: string;
  imageUrl?: string;
}

const ArtistCreate = () => {
  const [artists, setArtists] = useState<Artist[]>([]);
  const [open, setOpen] = useState(false);
  const [editingArtist, setEditingArtist] = useState<Artist | null>(null);
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    imageUrl: '',
  });
  const [error, setError] = useState<string | null>(null);

  const fetchArtists = async () => {
    try {
      const response = await apiService.get('/artists');
      const artistsData = Array.isArray(response) ? response : [];
      setArtists(artistsData);
    } catch (err) {
      setError('Ошибка при загрузке артистов');
      console.error('Error fetching artists:', err);
      setArtists([]);
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
        imageUrl: artist.imageUrl || '',
      });
    } else {
      setEditingArtist(null);
      setFormData({
        name: '',
        description: '',
        imageUrl: '',
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
      imageUrl: '',
    });
    setError(null);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (editingArtist) {
        await apiService.put(`/artists/${editingArtist.id}`, formData);
      } else {
        await apiService.post('/artists', formData);
      }
      handleClose();
      fetchArtists();
    } catch (err) {
      setError('Ошибка при сохранении артиста');
      console.error('Error saving artist:', err);
    }
  };

  const handleDelete = async (id: number) => {
    if (window.confirm('Вы уверены, что хотите удалить этого артиста?')) {
      try {
        await apiService.delete(`/artists/${id}`);
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
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
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

        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Имя</TableCell>
                <TableCell>Описание</TableCell>
                <TableCell>Изображение</TableCell>
                <TableCell align="right">Действия</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {artists.map((artist) => (
                <TableRow key={artist.id}>
                  <TableCell>{artist.name}</TableCell>
                  <TableCell>{artist.description}</TableCell>
                  <TableCell>
                    {artist.imageUrl && (
                      <img
                        src={artist.imageUrl}
                        alt={artist.name}
                        style={{ width: 50, height: 50, objectFit: 'cover', borderRadius: 4 }}
                      />
                    )}
                  </TableCell>
                  <TableCell align="right">
                    <IconButton onClick={() => handleOpen(artist)} color="primary">
                      <EditIcon />
                    </IconButton>
                    <IconButton onClick={() => handleDelete(artist.id)} color="error">
                      <DeleteIcon />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>

      <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
        <DialogTitle>
          {editingArtist ? 'Редактировать артиста' : 'Добавить артиста'}
        </DialogTitle>
        <form onSubmit={handleSubmit}>
          <DialogContent>
            <Stack spacing={2}>
              <TextField
                autoFocus
                label="Имя артиста"
                fullWidth
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                required
              />
              <TextField
                label="Описание"
                fullWidth
                multiline
                rows={3}
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              />
              <TextField
                label="URL изображения"
                fullWidth
                value={formData.imageUrl}
                onChange={(e) => setFormData({ ...formData, imageUrl: e.target.value })}
                placeholder="https://example.com/image.jpg"
              />
            </Stack>
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

export default ArtistCreate; 