import { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  Box,
  Typography,
} from '@mui/material';
import type { Playlist } from '../../types';

interface PlaylistDialogProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (name: string) => Promise<void>;
  playlist?: Playlist;
  title: string;
}

const PlaylistDialog = ({
  open,
  onClose,
  onSubmit,
  playlist,
  title,
}: PlaylistDialogProps) => {
  const [name, setName] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (playlist) {
      setName(playlist.name);
    } else {
      setName('');
    }
  }, [playlist]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (!name.trim()) {
      setError('Playlist name is required');
      return;
    }

    try {
      setLoading(true);
      await onSubmit(name);
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save playlist');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <form onSubmit={handleSubmit}>
        <DialogTitle>{title}</DialogTitle>
        <DialogContent>
          <Box sx={{ pt: 2 }}>
            <TextField
              autoFocus
              fullWidth
              label="Playlist Name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              error={!!error}
              helperText={error}
              disabled={loading}
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose} disabled={loading}>
            Cancel
          </Button>
          <Button
            type="submit"
            variant="contained"
            disabled={loading || !name.trim()}
          >
            {loading ? 'Saving...' : 'Save'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};

export default PlaylistDialog; 