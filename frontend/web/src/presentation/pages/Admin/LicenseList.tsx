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
import { apiService } from '../../../presentation/services/api';

interface License {
  id: string;
  title: string;
  description: string;
  url: string;
}

const LicenseList = () => {
  const [licenses, setLicenses] = useState<License[]>([]);
  const [open, setOpen] = useState(false);
  const [editingLicense, setEditingLicense] = useState<License | null>(null);
  const [formData, setFormData] = useState({ title: '', description: '', url: '' });
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  const fetchLicenses = async () => {
    try {
      setLoading(true);
      const response = await apiService.get('/licenses');
      console.log('Received licenses data:', response);
      const licensesData = Array.isArray(response) ? response : [];
      console.log('Processed licenses data:', licensesData);
      setLicenses(licensesData);
      setError(null);
    } catch (err) {
      setError('Ошибка при загрузке лицензий');
      console.error('Error fetching licenses:', err);
      setLicenses([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchLicenses();
  }, []);

  const handleOpen = (license?: License) => {
    if (license) {
      setEditingLicense(license);
      setFormData({ title: license.title, description: license.description, url: license.url });
    } else {
      setEditingLicense(null);
      setFormData({ title: '', description: '', url: '' });
    }
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
    setEditingLicense(null);
    setFormData({ title: '', description: '', url: '' });
    setError(null);
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const trimmedData = {
        title: formData.title.trim(),
        description: formData.description.trim(),
        url: formData.url.trim(),
      };

      console.log('Sending data:', trimmedData);

      if (editingLicense) {
        await apiService.put(`/licenses/${editingLicense.id}`, trimmedData);
      } else {
        await apiService.post('/licenses/', trimmedData);
      }
      handleClose();
      fetchLicenses();
    } catch (err) {
      console.error('Error saving license:', err);
      setError('Ошибка при сохранении лицензии');
    }
  };

  const handleDelete = async (id: string) => {
    if (window.confirm('Вы уверены, что хотите удалить эту лицензию?')) {
      try {
        await apiService.delete(`/licenses/${id}`);
        fetchLicenses();
      } catch (err) {
        setError('Ошибка при удалении лицензии');
        console.error('Error deleting license:', err);
      }
    }
  };

  return (
    <Box sx={{ p: 4, maxWidth: 1200, mx: 'auto' }}>
      <Paper sx={{ p: 4 }} elevation={3}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Typography variant="h5" fontWeight={700}>
            Управление лицензиями
          </Typography>
          <Button variant="contained" color="primary" onClick={() => handleOpen()} sx={{ ml: 4 }}>
            Добавить лицензию
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
                  <TableCell>Описание</TableCell>
                  <TableCell>URL лицензии</TableCell>
                  <TableCell align="right">Действия</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {licenses.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={5} align="center">
                      Нет доступных лицензий
                    </TableCell>
                  </TableRow>
                ) : (
                  licenses.map((license) => (
                    <TableRow key={license.id}>
                      <TableCell>{license.id}</TableCell>
                      <TableCell>{license.title}</TableCell>
                      <TableCell>{license.description}</TableCell>
                      <TableCell>
                        {license.url && (
                          <a href={license.url} target="_blank" rel="noopener noreferrer">
                            {license.url}
                          </a>
                        )}
                      </TableCell>
                      <TableCell align="right">
                        <IconButton onClick={() => handleOpen(license)} color="primary">
                          <EditIcon />
                        </IconButton>
                        <IconButton onClick={() => handleDelete(license.id)} color="error">
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
          {editingLicense ? 'Редактировать лицензию' : 'Добавить лицензию'}
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
            <TextField
              margin="dense"
              label="Описание"
              fullWidth
              multiline
              rows={3}
              value={formData.description}
              onChange={handleChange}
              name="description"
            />
            <TextField
              margin="dense"
              label="URL лицензии"
              fullWidth
              value={formData.url}
              onChange={handleChange}
              name="url"
              placeholder="https://example.com/license"
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose}>Отмена</Button>
            <Button type="submit" variant="contained" color="primary">
              {editingLicense ? 'Сохранить' : 'Добавить'}
            </Button>
          </DialogActions>
        </form>
      </Dialog>
    </Box>
  );
};

export default LicenseList; 