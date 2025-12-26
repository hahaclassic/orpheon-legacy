import { Card, CardContent, CardMedia, Typography, Box, Chip } from "@mui/material";
import { useNavigate } from "react-router-dom";

interface ArtistCardProps {
  id: string;
  name: string;
  country: string;
  genre: string;
  albumCount?: number;
  coverImage?: string;
}

const ArtistCard = ({ id, name, country, genre, coverImage }: ArtistCardProps) => {
  const navigate = useNavigate();

  const handleArtistClick = () => {
    navigate(`/artists/${id}`);
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
      onClick={handleArtistClick}
    >
      <CardMedia
        component="img"
        sx={{
          height: 250,
          width: '100%',
          objectFit: 'cover',
          aspectRatio: '1/1'
        }}
        image={coverImage || "/default-cover.jpg"}
        alt={name}
      />
      <CardContent>
        <Typography gutterBottom variant="h6" component="div" noWrap>
          {name}
        </Typography>
        <Box sx={{ mt: 1, display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
          <Chip
            label={country}
            size="small"
            sx={{ mr: 0.5, mb: 0.5 }}
          />
          <Chip
            label={genre}
            size="small"
            sx={{ mr: 0.5, mb: 0.5 }}
          />
        </Box>
      </CardContent>
    </Card>
  );
};

export default ArtistCard; 