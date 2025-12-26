import { useState } from 'react';
import { Box, Card, CardContent, CardMedia, IconButton, Typography } from '@mui/material';
import { useNavigate } from 'react-router-dom';
import ImageIcon from '@mui/icons-material/Image';
import FavoriteIcon from '@mui/icons-material/Favorite';
import FavoriteBorderIcon from '@mui/icons-material/FavoriteBorder';
import { apiService } from '../../../presentation/services/api';

interface Owner {
  id: string;
  name: string;
}

interface PlaylistCardProps {
  id: string;
  name: string;
  coverUrl?: string;
  isFavorite?: boolean;
  rating?: number;
  owner?: Owner;
  ownerId?: string;
  ownerName?: string;
  onFavoriteChange?: (isFavorite: boolean) => void;
}

const PlaylistCard = ({
  id,
  name,
  coverUrl,
  isFavorite = false,
  rating = 0,
  owner,
  ownerId,
  ownerName,
  onFavoriteChange,
}: PlaylistCardProps) => {
  const navigate = useNavigate();
  const [updatingFavorite, setUpdatingFavorite] = useState(false);
  const [localIsFavorite, setLocalIsFavorite] = useState(isFavorite);
  const [localRating, setLocalRating] = useState(rating);

  const handleCardClick = () => {
    navigate(`/playlists/${id}`);
  };

  const handleOwnerClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    const ownerIdToNavigate = owner?.id || ownerId;
    if (ownerIdToNavigate) {
      navigate(`/users/${ownerIdToNavigate}`);
    }
  };

  const handleFavoriteClick = async (e: React.MouseEvent) => {
    e.stopPropagation();
    if (updatingFavorite) return;

    try {
      setUpdatingFavorite(true);
      if (localIsFavorite) {
        await apiService.delete(`/me/favorites/${id}`);
        setLocalIsFavorite(false);
        setLocalRating(prev => prev - 1);
      } else {
        await apiService.post(`/me/favorites/${id}`);
        setLocalIsFavorite(true);
        setLocalRating(prev => prev + 1);
      }
      onFavoriteChange?.(!localIsFavorite);
    } catch (err) {
      console.error("Error updating favorite status:", err);
    } finally {
      setUpdatingFavorite(false);
    }
  };

  const displayOwnerName = owner?.name || ownerName;

  return (
    <Card
      sx={{
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
        cursor: 'pointer',
        '&:hover': {
          transform: 'scale(1.02)',
          transition: 'transform 0.2s ease-in-out',
        },
      }}
      onClick={handleCardClick}
    >
      {coverUrl ? (
        <CardMedia
          component="img"
          sx={{
            height: 250,
            width: '100%',
            objectFit: 'cover',
            aspectRatio: '1/1'
          }}
          image={coverUrl}
          alt={name}
        />
      ) : (
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
          <ImageIcon sx={{ fontSize: 64, color: 'primary.contrastText', opacity: 0.3 }} />
        </Box>
      )}
      <CardContent>
        <Typography gutterBottom variant="h6" component="div" noWrap>
          {name}
        </Typography>
        {displayOwnerName && (
          <Typography
            variant="body2"
            color="text.secondary"
            sx={{
              cursor: 'pointer',
              '&:hover': {
                textDecoration: 'underline',
              },
            }}
            onClick={handleOwnerClick}
          >
            {displayOwnerName}
          </Typography>
        )}
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5, mt: 1 }}>
          <IconButton 
            onClick={handleFavoriteClick}
            size="small"
            disabled={updatingFavorite}
            sx={{ 
              color: 'white',
              '&:hover': { color: 'white' }
            }}
          >
            {localIsFavorite ? (
              <FavoriteIcon sx={{ fontSize: 20 }} />
            ) : (
              <FavoriteBorderIcon sx={{ fontSize: 20 }} />
            )}
          </IconButton>
          <Typography variant="body2" color="white">
            {localRating}
        </Typography>
        </Box>
      </CardContent>
    </Card>
  );
};

export default PlaylistCard; 