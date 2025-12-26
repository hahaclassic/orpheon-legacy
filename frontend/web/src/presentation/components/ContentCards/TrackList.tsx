import { List, ListItem, ListItemText, ListItemAvatar, Avatar, Typography, Box, Divider } from "@mui/material";
import { useNavigate } from "react-router-dom";

interface Track {
  id: string;
  title: string;
  artistName: string;
  albumName: string;
  duration: number;
  coverImage?: string;
}

interface TrackListProps {
  tracks: Track[];
}

const formatDuration = (seconds: number) => {
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = seconds % 60;
  return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
};

const TrackList = ({ tracks }: TrackListProps) => {
  const navigate = useNavigate();

  const handleTrackClick = (trackId: string) => {
    navigate(`/tracks/${trackId}`);
  };

  return (
    <List sx={{ width: '100%', bgcolor: 'background.paper' }}>
      {tracks.map((track, index) => (
        <Box key={track.id}>
          <ListItem
            button
            onClick={() => handleTrackClick(track.id)}
            sx={{
              '&:hover': {
                backgroundColor: 'rgba(0, 0, 0, 0.04)',
              },
            }}
          >
            <ListItemAvatar>
              <Avatar
                variant="rounded"
                src={track.coverImage || "/default-cover.jpg"}
                alt={track.title}
                sx={{ width: 56, height: 56 }}
              />
            </ListItemAvatar>
            <ListItemText
              primary={track.title}
              secondary={
                <Box component="span" sx={{ display: 'flex', flexDirection: 'column' }}>
                  <Typography variant="body2" color="text.secondary">
                    {track.artistName}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    {track.albumName}
                  </Typography>
                </Box>
              }
            />
            <Typography variant="body2" color="text.secondary">
              {formatDuration(track.duration)}
            </Typography>
          </ListItem>
          {index < tracks.length - 1 && <Divider />}
        </Box>
      ))}
    </List>
  );
};

export default TrackList; 