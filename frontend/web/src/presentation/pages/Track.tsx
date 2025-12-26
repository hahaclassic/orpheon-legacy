import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Box, Container, Typography, Grid, Card, CardContent, CardMedia, Chip, Link } from '@mui/material';
import { Person, Album, AccessTime } from '@mui/icons-material';

interface Track {
  id: string;
  title: string;
  duration: string;
  releaseDate: string;
  genres: string[];
  artist: {
    id: string;
    name: string;
  };
  album: {
    id: string;
    title: string;
    imageUrl: string;
  };
  description: string;
}

const TrackPage = () => {
  const { id } = useParams<{ id: string }>();
  const [track, setTrack] = useState<Track | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // TODO: Fetch track data from API
    // This is a placeholder for the actual API call
    const fetchTrackData = async () => {
      try {
        // const response = await api.getTrack(id);
        // setTrack(response.data);
        setLoading(false);
      } catch (error) {
        console.error('Error fetching track data:', error);
        setLoading(false);
      }
    };

    fetchTrackData();
  }, [id]);

  if (loading) {
    return <Typography>Loading...</Typography>;
  }

  if (!track) {
    return <Typography>Track not found</Typography>;
  }

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Grid container spacing={4}>
        {/* Track Header */}
        <Grid item xs={12} md={4}>
          <Card>
            <CardMedia
              component="img"
              image={track.album.imageUrl}
              alt={track.title}
              sx={{ aspectRatio: '1/1' }}
            />
          </Card>
        </Grid>
        <Grid item xs={12} md={8}>
          <Box sx={{ mb: 2 }}>
            <Typography variant="h3" component="h1" gutterBottom>
              {track.title}
            </Typography>
            
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2 }}>
              <Link
                href={`/artists/${track.artist.id}`}
                sx={{ display: 'flex', alignItems: 'center', gap: 1, textDecoration: 'none', color: 'inherit' }}
              >
                <Person /> {track.artist.name}
              </Link>
              <Link
                href={`/albums/${track.album.id}`}
                sx={{ display: 'flex', alignItems: 'center', gap: 1, textDecoration: 'none', color: 'inherit' }}
              >
                <Album /> {track.album.title}
              </Link>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <AccessTime /> {track.duration}
              </Box>
            </Box>

            <Typography variant="body1" color="text.secondary" sx={{ mt: 2 }}>
              {track.description}
            </Typography>

            <Box sx={{ mt: 2, display: 'flex', gap: 1, flexWrap: 'wrap' }}>
              {track.genres.map((genre) => (
                <Chip key={genre} label={genre} size="small" />
              ))}
            </Box>

            <Typography variant="body2" color="text.secondary" sx={{ mt: 2 }}>
              Released: {new Date(track.releaseDate).toLocaleDateString()}
            </Typography>
          </Box>
        </Grid>
      </Grid>
    </Container>
  );
};

export default TrackPage; 