import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import {
  Box,
  Container,
  Typography,
  Grid,
  Card,
  CardContent,
  CardMedia,
  CircularProgress,
  Alert,
  Chip,
  Link,
} from "@mui/material";
import { apiService } from "../../../presentation/services/api";

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

interface Album {
  id: string;
  title: string;
  label: string;
  license_id: string;
  release_date: string;
  artists: Artist[];
  genres: Genre[];
  coverUrl?: string;
}

const Home = () => {
  const [albums, setAlbums] = useState<Album[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  const fetchAlbums = async () => {
    try {
      const response = await apiService.get('/albums');
      const albumsData = Array.isArray(response) ? response : [];
      
      // Получаем обложки для каждого альбома
      const albumsWithCovers = await Promise.all(
        albumsData.map(async (album) => {
          try {
            const coverResponse = await apiService.get(`/albums/${album.id}/cover`, {
              responseType: 'blob'
            });
            const coverUrl = URL.createObjectURL(coverResponse);
            return { ...album, coverUrl };
          } catch (err) {
            console.error(`Error fetching cover for album ${album.id}:`, err);
            return { ...album, coverUrl: "/default-cover.jpg" };
          }
        })
      );
      
      setAlbums(albumsWithCovers);
    } catch (err) {
      console.error("Error fetching albums:", err);
      setError("Ошибка при загрузке альбомов");
    }
  };

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      setError(null);
      try {
        await fetchAlbums();
      } catch (err) {
        console.error("Error fetching data:", err);
        setError("Ошибка при загрузке данных");
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  const handleAlbumClick = (albumId: string) => {
    navigate(`/albums/${albumId}`);
  };

  const handleArtistClick = (e: React.MouseEvent, artistId: string) => {
    e.stopPropagation(); // Предотвращаем переход на страницу альбома
    navigate(`/artists/${artistId}`);
  };

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      {loading ? (
        <Box sx={{ display: "flex", justifyContent: "center", p: 3 }}>
          <CircularProgress />
        </Box>
      ) : (
        <>
          <Typography variant="h4" component="h1" gutterBottom>
            Все альбомы
          </Typography>
          <Grid container spacing={3}>
            {albums.length === 0 ? (
              <Grid item xs={12}>
                <Typography variant="h6" textAlign="center" color="text.secondary">
                  Нет доступных альбомов
                </Typography>
              </Grid>
            ) : (
              albums.map((album) => (
                <Grid item xs={12} sm={6} md={3} key={album.id}>
                  <Card
                    sx={{
                      height: "100%",
                      display: "flex",
                      flexDirection: "column",
                      cursor: "pointer",
                      "&:hover": {
                        transform: "scale(1.02)",
                        transition: "transform 0.2s ease-in-out",
                      },
                    }}
                  >
                    <CardMedia
                      component="img"
                      sx={{
                        height: 250,
                        width: '100%',
                        objectFit: 'cover',
                        aspectRatio: '1/1'
                      }}
                      image={album.coverUrl || "/default-cover.jpg"}
                      alt={album.title}
                      onClick={() => handleAlbumClick(album.id)}
                    />
                    <CardContent>
                      <Typography gutterBottom variant="h6" component="div" noWrap>
                        {album.title}
                      </Typography>
                      <Typography variant="body2" color="text.secondary" noWrap>
                        {album.artists.map((artist, index) => (
                          <span key={artist.id}>
                            <Link
                              component="button"
                              variant="body2"
                              onClick={(e) => handleArtistClick(e, artist.id)}
                              sx={{ 
                                color: 'inherit',
                                textDecoration: 'none',
                                '&:hover': { textDecoration: 'underline' }
                              }}
                            >
                              {artist.name}
                            </Link>
                            {index < album.artists.length - 1 ? ", " : ""}
                          </span>
                        ))}
                      </Typography>
                      <Box sx={{ mt: 1, display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                        {album.genres.map((genre) => (
                          <Chip
                            key={genre.id}
                            label={genre.title}
                            size="small"
                            sx={{ mr: 0.5, mb: 0.5 }}
                          />
                        ))}
                      </Box>
                      <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                        {new Date(album.release_date).getFullYear()}
                      </Typography>
                    </CardContent>
                  </Card>
                </Grid>
              ))
            )}
          </Grid>
        </>
      )}
    </Container>
  );
};

export default Home; 