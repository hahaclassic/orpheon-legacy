import { Card, CardContent, CardMedia, Typography, Box, Link, Chip, CircularProgress } from "@mui/material";
import { useNavigate } from "react-router-dom";
import ImageIcon from '@mui/icons-material/Image';
import { useCoverImage } from "../../hooks/useCoverImage";

interface Artist {
  id: string;
  name: string;
}

interface Genre {
  id: string;
  title: string;
}

interface AlbumCardProps {
  id: string;
  title: string;
  label?: string;
  release_date?: string;
  artists?: Artist[] | null;
  genres?: Genre[] | null;
}

const AlbumCard = ({ 
  id, 
  title, 
  label = '', 
  release_date = '',
  artists = [],
  genres = []
}: AlbumCardProps) => {
  const navigate = useNavigate();
  const { coverUrl, loading, error } = useCoverImage('album', id);

  const handleAlbumClick = () => {
    navigate(`/albums/${id}`);
  };

  const handleArtistClick = (e: React.MouseEvent, artistId: string) => {
    e.stopPropagation(); // Предотвращаем переход на страницу альбома
    navigate(`/artists/${artistId}`);
  };

  return (
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
      onClick={handleAlbumClick}
    >
      {loading ? (
        <Box
          sx={{
            height: 250,
            width: '100%',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            bgcolor: 'primary.dark',
            borderRadius: 2,
          }}
        >
          <CircularProgress color="inherit" />
        </Box>
      ) : (
        <CardMedia
          component="img"
          sx={{
            height: 250,
            width: '100%',
            objectFit: 'cover',
            aspectRatio: '1/1'
          }}
          image={coverUrl || "/default-cover.jpg"}
          alt={title}
        />
      )}
      <CardContent>
        <Typography gutterBottom variant="h6" component="div" noWrap>
          {title}
        </Typography>
        {artists && artists.length > 0 && (
          <Typography variant="body2" color="text.secondary" noWrap>
            {artists.map((artist, index) => (
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
                {index < artists.length - 1 ? ", " : ""}
              </span>
            ))}
          </Typography>
        )}
        {genres && genres.length > 0 && (
          <Box sx={{ mt: 1, display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
            {genres.map((genre) => (
              <Chip
                key={genre.id}
                label={genre.title}
                size="small"
                sx={{ mr: 0.5, mb: 0.5 }}
              />
            ))}
          </Box>
        )}
        {label && (
          <Typography variant="body2" color="text.secondary" noWrap>
            {label}
          </Typography>
        )}
        {release_date && (
          <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
            {new Date(release_date).getFullYear()}
          </Typography>
        )}
      </CardContent>
    </Card>
  );
};

export default AlbumCard; 