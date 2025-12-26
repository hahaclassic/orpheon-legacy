import { useState } from 'react';
import { 
  Box, 
  Typography, 
  IconButton,
  Menu,
  MenuItem,
  Tooltip,
  Link,
} from '@mui/material';
import { PlayArrow, Pause, Add as AddIcon, MoreVert, Headphones } from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import { usePlayerContext } from '../contexts/PlayerContext';
import { useAuthContext } from '../contexts/AuthContext';
import type { Track } from '../types';
import React from 'react';

interface TrackItemProps {
  track: Track;
  tracks: Track[];
  index: number;
  onAddToPlaylist?: (event: React.MouseEvent<HTMLElement>, track: Track) => void;
  showTrackNumber?: boolean;
  showAlbumLink?: boolean;
  onTrackClick?: (trackId: string) => void;
}

const formatDuration = (seconds: number) => {
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = Math.floor(seconds % 60);
  return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
};

const formatNumber = (num: number) => {
  if (num >= 1000000) {
    return `${(num / 1000000).toFixed(1)}M`;
  }
  if (num >= 1000) {
    return `${(num / 1000).toFixed(1)}K`;
  }
  return num.toString();
};

const TrackCover = ({ track, index, currentTrackId, isPlaying }: { 
  track: Track; 
  index: number; 
  currentTrackId?: string; 
  isPlaying?: boolean;
}) => (
  <Box
    sx={{
      width: 56,
      height: 56,
      position: 'relative',
      mr: 2,
      flexShrink: 0,
      '&:hover .play-icon': {
        opacity: 1,
      },
    }}
  >
    {track.coverUrl ? (
      <Box
        component="img"
        src={track.coverUrl}
        alt={track.name}
        sx={{
          width: '100%',
          height: '100%',
          borderRadius: 1,
          objectFit: 'cover',
        }}
        onError={(e) => {
          console.error('Error loading track cover:', track.coverUrl);
          e.currentTarget.style.display = 'none';
        }}
      />
    ) : (
      <Box
        sx={{
          width: '100%',
          height: '100%',
          borderRadius: 1,
          bgcolor: 'action.hover',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
        }}
      >
        <Typography variant="body2" color="text.secondary">
          {track.track_number || index + 1}
        </Typography>
      </Box>
    )}
    <IconButton
      size="small"
      className="play-icon"
      sx={{
        position: 'absolute',
        top: '50%',
        left: '50%',
        transform: 'translate(-50%, -50%)',
        opacity: currentTrackId === track.id ? 1 : 0,
        transition: 'opacity 0.2s',
        bgcolor: 'rgba(0, 0, 0, 0.5)',
        '&:hover': {
          bgcolor: 'rgba(0, 0, 0, 0.7)',
        },
      }}
    >
      {currentTrackId === track.id && isPlaying ? <Pause /> : <PlayArrow />}
    </IconButton>
  </Box>
);

const TrackInfo = ({ track }: { track: Track }) => (
  <Box sx={{ flex: 1, minWidth: 0 }}>
    <Typography
      variant="subtitle1"
      sx={{
        fontWeight: 500,
        whiteSpace: 'nowrap',
        overflow: 'hidden',
        textOverflow: 'ellipsis',
      }}
    >
      {track.name}
    </Typography>
    <Typography
      variant="body2"
      sx={{
        color: 'text.secondary',
        whiteSpace: 'nowrap',
        overflow: 'hidden',
        textOverflow: 'ellipsis',
      }}
    >
      {formatDuration(track.duration)} • {track.artists.map((artist, index) => (
        <React.Fragment key={artist.id}>
          {index > 0 && ', '}
          <Link
            href={`/artists/${artist.id}`}
            onClick={(e) => e.stopPropagation()}
            sx={{
              color: 'text.secondary',
              textDecoration: 'none',
              '&:hover': {
                textDecoration: 'underline',
              },
            }}
          >
            {artist.name}
          </Link>
        </React.Fragment>
      ))}
      {track.license && (
        <>
          {' • '}
          <Tooltip title={track.license.description}>
            <Link
              href={track.license.url}
              target="_blank"
              rel="noopener noreferrer"
              onClick={(e) => e.stopPropagation()}
              sx={{
                color: 'text.secondary',
                textDecoration: 'none',
                '&:hover': {
                  textDecoration: 'underline',
                },
              }}
            >
              {track.license.title}
            </Link>
          </Tooltip>
        </>
      )}
    </Typography>
  </Box>
);

const TrackStats = ({ totalStreams }: { totalStreams?: number }) => {
  if (totalStreams === undefined) return null;
  
  return (
    <Box
      sx={{
        display: 'flex',
        alignItems: 'center',
        gap: 0.5,
        ml: 2,
        mr: 1,
        color: 'text.secondary',
      }}
    >
      <Headphones sx={{ fontSize: 16 }} />
      <Typography variant="body2" color="text.secondary">
        {formatNumber(totalStreams)}
      </Typography>
    </Box>
  );
};

const TrackItem = ({
  track,
  tracks,
  index,
  onAddToPlaylist,
  showTrackNumber = true,
  showAlbumLink = true,
  onTrackClick,
}: TrackItemProps) => {
  const [menuAnchorEl, setMenuAnchorEl] = useState<null | HTMLElement>(null);
  const { state, controls } = usePlayerContext();
  const { currentTrack, isPlaying } = state;
  const { startPlayback, togglePlay } = controls;
  const navigate = useNavigate();
  const { user } = useAuthContext();

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    event.stopPropagation();
    setMenuAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setMenuAnchorEl(null);
  };

  const handleTrackClick = () => {
    if (onTrackClick) {
      onTrackClick(track.id);
    } else if (currentTrack?.id === track.id) {
      togglePlay();
    } else {
      startPlayback(track, tracks);
    }
  };

  const handleAddToPlaylistClick = (e: React.MouseEvent<HTMLElement>) => {
    e.stopPropagation();
    if (!user) {
      navigate('/login', { 
        state: { 
          from: window.location.pathname,
          message: 'Чтобы добавить трек в плейлист, необходимо войти'
        }
      });
      return;
    }
    if (onAddToPlaylist) {
      onAddToPlaylist(e, track);
    }
  };

  return (
    <Box
      sx={{
        display: 'flex',
        alignItems: 'center',
        p: 1,
        cursor: 'pointer',
        '&:hover': {
          backgroundColor: 'action.hover',
        },
      }}
      onClick={handleTrackClick}
    >
      <TrackCover 
        track={track} 
        index={index} 
        currentTrackId={currentTrack?.id} 
        isPlaying={isPlaying} 
      />
      <TrackInfo track={track} />
      <TrackStats totalStreams={track.total_streams} />
      
      {onAddToPlaylist && (
        <IconButton
          size="small"
          onClick={handleAddToPlaylistClick}
        >
          <AddIcon />
        </IconButton>
      )}
      
      <IconButton size="small" onClick={handleMenuOpen}>
        <MoreVert />
      </IconButton>
      
      <Menu
        anchorEl={menuAnchorEl}
        open={Boolean(menuAnchorEl)}
        onClose={handleMenuClose}
        onClick={(e) => e.stopPropagation()}
      >
        <MenuItem onClick={() => {
          handleMenuClose();
          navigate(`/albums/${track.album.id}`);
        }}>
          Перейти к альбому
        </MenuItem>
      </Menu>
    </Box>
  );
};

export default TrackItem; 