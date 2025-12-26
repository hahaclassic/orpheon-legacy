import React, { useState } from 'react';
import { Box, Button, TextField, Typography, Stepper, Step, StepLabel, Paper, MenuItem, Select, FormControl, InputLabel } from '@mui/material';
import { useForm, Controller } from 'react-hook-form';
import { useNavigate } from 'react-router-dom';

interface AlbumMeta {
  title: string;
  description: string;
  releaseDate: string;
  artistId: string;
  genreId: string;
  licenseId: string;
}

interface TrackMeta {
  title: string;
  duration: number;
  audioFile: File;
}

const steps = ['Метаданные альбома', 'Добавление треков'];

export const AlbumUpload: React.FC = () => {
  const [activeStep, setActiveStep] = useState(0);
  const [albumId, setAlbumId] = useState<string | null>(null);
  const [tracks, setTracks] = useState<TrackMeta[]>([]);
  const navigate = useNavigate();

  const { control: albumControl, handleSubmit: handleAlbumSubmit } = useForm<AlbumMeta>();
  const { control: trackControl, handleSubmit: handleTrackSubmit, reset: resetTrackForm } = useForm<TrackMeta>();

  const handleNext = () => {
    setActiveStep((prevStep) => prevStep + 1);
  };

  const handleBack = () => {
    setActiveStep((prevStep) => prevStep - 1);
  };

  const onSubmitAlbum = async (data: AlbumMeta) => {
    try {
      const response = await fetch('/albums', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        throw new Error('Failed to create album');
      }

      const { id } = await response.json();
      setAlbumId(id);
      handleNext();
    } catch (error) {
      console.error('Error creating album:', error);
    }
  };

  const onSubmitTrack = async (data: TrackMeta) => {
    if (!albumId) return;

    const formData = new FormData();
    formData.append('audio', data.audioFile);
    formData.append('title', data.title);
    formData.append('duration', data.duration.toString());
    formData.append('albumId', albumId);

    try {
      const response = await fetch('/tracks', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
        body: formData,
      });

      if (!response.ok) {
        throw new Error('Failed to upload track');
      }

      setTracks([...tracks, data]);
      resetTrackForm();
    } catch (error) {
      console.error('Error uploading track:', error);
    }
  };

  const renderAlbumForm = () => (
    <form onSubmit={handleAlbumSubmit(onSubmitAlbum)}>
      <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
        <Controller
          name="title"
          control={albumControl}
          defaultValue=""
          render={({ field }) => (
            <TextField
              {...field}
              label="Название альбома"
              required
              fullWidth
            />
          )}
        />
        <Controller
          name="description"
          control={albumControl}
          defaultValue=""
          render={({ field }) => (
            <TextField
              {...field}
              label="Описание"
              multiline
              rows={4}
              fullWidth
            />
          )}
        />
        <Controller
          name="releaseDate"
          control={albumControl}
          defaultValue=""
          render={({ field }) => (
            <TextField
              {...field}
              label="Дата выпуска"
              type="date"
              required
              fullWidth
              InputLabelProps={{ shrink: true }}
            />
          )}
        />
        <Controller
          name="artistId"
          control={albumControl}
          defaultValue=""
          render={({ field }) => (
            <TextField
              {...field}
              label="ID артиста"
              required
              fullWidth
            />
          )}
        />
        <Controller
          name="genreId"
          control={albumControl}
          defaultValue=""
          render={({ field }) => (
            <FormControl fullWidth required>
              <InputLabel>Жанр</InputLabel>
              <Select
                {...field}
                label="Жанр"
              >
                {/* Здесь будет список жанров */}
              </Select>
            </FormControl>
          )}
        />
        <Controller
          name="licenseId"
          control={albumControl}
          defaultValue=""
          render={({ field }) => (
            <FormControl fullWidth required>
              <InputLabel>Лицензия</InputLabel>
              <Select
                {...field}
                label="Лицензия"
              >
                {/* Здесь будет список лицензий */}
              </Select>
            </FormControl>
          )}
        />
        <Button type="submit" variant="contained" color="primary">
          Далее
        </Button>
      </Box>
    </form>
  );

  const renderTrackForm = () => (
    <form onSubmit={handleTrackSubmit(onSubmitTrack)}>
      <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
        <Controller
          name="title"
          control={trackControl}
          defaultValue=""
          render={({ field }) => (
            <TextField
              {...field}
              label="Название трека"
              required
              fullWidth
            />
          )}
        />
        <Controller
          name="duration"
          control={trackControl}
          defaultValue={0}
          render={({ field }) => (
            <TextField
              {...field}
              label="Длительность (в секундах)"
              type="number"
              required
              fullWidth
            />
          )}
        />
        <Controller
          name="audioFile"
          control={trackControl}
          render={({ field: { value, onChange, ...field } }) => (
            <Button
              variant="outlined"
              component="label"
              fullWidth
            >
              Загрузить аудиофайл
              <input
                type="file"
                accept="audio/*"
                hidden
                onChange={(e) => {
                  const file = e.target.files?.[0];
                  if (file) {
                    onChange(file);
                  }
                }}
                {...field}
              />
            </Button>
          )}
        />
        <Box sx={{ display: 'flex', gap: 2 }}>
          <Button onClick={handleBack} variant="outlined">
            Назад
          </Button>
          <Button type="submit" variant="contained" color="primary">
            Добавить трек
          </Button>
        </Box>
      </Box>
    </form>
  );

  return (
    <Box sx={{ maxWidth: 600, mx: 'auto', p: 3 }}>
      <Paper sx={{ p: 3 }}>
        <Typography variant="h4" gutterBottom>
          Загрузка альбома
        </Typography>
        <Stepper activeStep={activeStep} sx={{ mb: 4 }}>
          {steps.map((label) => (
            <Step key={label}>
              <StepLabel>{label}</StepLabel>
            </Step>
          ))}
        </Stepper>
        {activeStep === 0 ? renderAlbumForm() : renderTrackForm()}
        {tracks.length > 0 && (
          <Box sx={{ mt: 4 }}>
            <Typography variant="h6">Добавленные треки:</Typography>
            {tracks.map((track, index) => (
              <Typography key={index}>{track.title}</Typography>
            ))}
          </Box>
        )}
      </Paper>
    </Box>
  );
}; 