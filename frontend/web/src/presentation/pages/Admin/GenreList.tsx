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
} from '@mui/material';
import { Edit as EditIcon, Delete as DeleteIcon } from '@mui/icons-material';
import { api } from '../../../presentation/services/api';

interface Genre {
  id: string;
  title: string;
}

const GenreList = () => {
  const [genres, setGenres] = useState<Genre[]>([]);
  const [open, setOpen] = useState(false);
  const [editingGenre, setEditingGenre] = useState<Genre | null>(null);
  const [formData, setFormData] = useState({ title: '' });
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  const fetchGenres = async () => {
    try {
      setLoading(true);
      const response = await api.getGenres();
      console.log('Received genres data:', response);
      const genresData = Array.isArray(response) ? response : [];
      console.log('Processed genres data:', genresData);
      setGenres(genresData);
      setError(null);
    } catch (err) {
      setError('Ошибка при загрузке жанров');
      console.error('Error fetching genres:', err);
      setGenres([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchGenres();
  }, []);

  const handleOpen = (genre?: Genre) => {
    if (genre) {
      setEditingGenre(genre);
      setFormData({ title: genre.title });
    } else {
      setEditingGenre(null);
      setFormData({ title: '' });
    }
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
    setEditingGenre(null);
    setFormData({ title: '' });
    setError(null);
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const trimmedData = {
        title: formData.title,
      };

      console.log('Sending data:', trimmedData);

      if (editingGenre) {
        await api.updateGenre(editingGenre.id, trimmedData);
      } else {
        await api.createGenre(trimmedData);
      }
      handleClose();
      fetchGenres();
    } catch (err) {
      console.error('Error saving genre:', err);
      setError('Ошибка при сохранении жанра');
    }
  };

  const handleDelete = async (id: string) => {
    if (window.confirm('Вы уверены, что хотите удалить этот жанр?')) {
      try {
        await api.deleteGenre(id);
        fetchGenres();
      } catch (err) {
        setError('Ошибка при удалении жанра');
        console.error('Error deleting genre:', err);
      }
    }
  };

  return (
    <Box sx={{ p: 4, maxWidth: 1200, mx: 'auto' }}>
      <Paper sx={{ p: 4 }} elevation={3}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Typography variant="h5" fontWeight={700}>
            Управление жанрами
          </Typography>
          <Button variant="contained" color="primary" onClick={() => handleOpen()} sx={{ ml: 4 }}>
            Добавить жанр
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
                  <TableCell align="right">Действия</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {genres.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={3} align="center">
                      Нет доступных жанров
                    </TableCell>
                  </TableRow>
                ) : (
                  genres.map((genre) => (
                    <TableRow key={genre.id}>
                      <TableCell>{genre.id}</TableCell>
                      <TableCell>{genre.title}</TableCell>
                      <TableCell align="right">
                        <IconButton onClick={() => handleOpen(genre)} color="primary">
                          <EditIcon />
                        </IconButton>
                        <IconButton onClick={() => handleDelete(genre.id)} color="error">
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
          {editingGenre ? 'Редактировать жанр' : 'Добавить жанр'}
        </DialogTitle>
        <form onSubmit={handleSubmit}>
          <DialogContent>
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
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose}>Отмена</Button>
            <Button type="submit" variant="contained" color="primary">
              {editingGenre ? 'Сохранить' : 'Добавить'}
            </Button>
          </DialogActions>
        </form>
      </Dialog>
    </Box>
  );
};

export default GenreList; 