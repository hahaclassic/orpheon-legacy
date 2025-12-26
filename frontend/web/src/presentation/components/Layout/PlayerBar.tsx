import { Box, IconButton, Typography, Slider, Menu, MenuItem } from '@mui/material';
import {
  PlayArrow,
  Pause,
  SkipNext,
  SkipPrevious,
  VolumeUp,
  VolumeOff,
  MusicNote,
  Add as AddIcon,
  MoreVert,
} from '@mui/icons-material';
import { usePlayerContext } from '../../contexts/PlayerContext';
import { useCallback, useEffect, useState } from 'react';
import { formatDuration } from '../../utils/time';
import React from 'react';
import { useNavigate } from 'react-router-dom';
import Check from '@mui/icons-material/Check';
import { apiService } from '../../services/api';
import { useAuthContext } from '../../contexts/AuthContext';
import type { Playlist } from '../../types';

// Моковые данные для графика
const mockSegments = [
  { idx: 0, totalStreams: 50, range: [0, 10] },
  { idx: 1, totalStreams: 30, range: [10, 20] },
  { idx: 2, totalStreams: 50, range: [20, 30] },
  { idx: 3, totalStreams: 40, range: [30, 40] },
  { idx: 4, totalStreams: 80, range: [40, 50] },
  { idx: 5, totalStreams: 20, range: [50, 60] },
  { idx: 6, totalStreams: 60, range: [60, 70] },
  { idx: 7, totalStreams: 30, range: [70, 80] },
  { idx: 8, totalStreams: 50, range: [80, 90] },
  { idx: 9, totalStreams: 20, range: [90, 100] },
  { idx: 10, totalStreams: 20, range: [100, 110] },
  { idx: 11, totalStreams: 30, range: [110, 120] },
  { idx: 12, totalStreams: 80, range: [120, 130] },
  { idx: 13, totalStreams: 100, range: [130, 140] },
  { idx: 14, totalStreams: 60, range: [140, 150] },
  { idx: 15, totalStreams: 20, range: [150, 160] },
  { idx: 16, totalStreams: 10, range: [160, 170] },
  { idx: 17, totalStreams: 10, range: [170, 180] },
];

// Компонент для отображения графика статистики
const WaveformStats: React.FC<{ segments: any[], progress: number, duration: number, setProgress: (v: number) => void }> = ({ segments, progress, duration, setProgress }) => {
  const viewBoxWidth = 1200;
  const viewBoxHeight = 60;
  const maxStreams = Math.max(...segments.map(s => s.totalStreams), 1);
  const points = segments.map((seg, i) => {
    const x = (i / (segments.length - 1)) * viewBoxWidth;
    const y = viewBoxHeight - (seg.totalStreams / maxStreams) * (viewBoxHeight - 8);
    return { x, y };
  });
  // Catmull-Rom to Bezier for smooth curve
  function catmullRom2bezier(points: {x: number, y: number}[]) {
    if (points.length < 2) return '';
    let d = `M ${points[0].x},${points[0].y}`;
    for (let i = 0; i < points.length - 1; i++) {
      const p0 = points[i === 0 ? i : i - 1];
      const p1 = points[i];
      const p2 = points[i + 1];
      const p3 = points[i + 2 < points.length ? i + 2 : i + 1];

      const cp1x = p1.x + (p2.x - p0.x) / 6;
      const cp1y = p1.y + (p2.y - p0.y) / 6;
      const cp2x = p2.x - (p3.x - p1.x) / 6;
      const cp2y = p2.y - (p3.y - p1.y) / 6;

      d += ` C ${cp1x},${cp1y} ${cp2x},${cp2y} ${p2.x},${p2.y}`;
    }
    return d;
  }
  const areaPath = catmullRom2bezier(points) +
    ` L ${viewBoxWidth},${viewBoxHeight} L 0,${viewBoxHeight} Z`;
  const progressPercent = duration > 0 ? progress / duration : 0;

  // Tooltip state
  const [hoverX, setHoverX] = useState<number | null>(null);
  const [hoverTime, setHoverTime] = useState<number | null>(null);

  const handleMouseMove = (e: React.MouseEvent<SVGSVGElement, MouseEvent>) => {
    const rect = e.currentTarget.getBoundingClientRect();
    const x = e.clientX - rect.left;
    // Переводим x в координаты viewBox
    const svgX = (x / rect.width) * viewBoxWidth;
    // Переводим svgX в секунды
    const time = duration * (svgX / viewBoxWidth);
    setHoverX(svgX);
    setHoverTime(time);
  };
  const handleMouseLeave = () => {
    setHoverX(null);
    setHoverTime(null);
  };

  // Форматирование времени
  function formatTooltipTime(sec: number) {
    if (isNaN(sec) || sec < 0) return '0:00';
    const m = Math.floor(sec / 60);
    const s = Math.floor(sec % 60);
    return `${m}:${s.toString().padStart(2, '0')}`;
  }

  const handleClick = (e: React.MouseEvent<SVGSVGElement, MouseEvent>) => {
    const rect = e.currentTarget.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const svgX = (x / rect.width) * viewBoxWidth;
    const time = duration * (svgX / viewBoxWidth);
    setProgress(time);
  };

  return (
    <Box sx={{ width: '100%', height: '100%', position: 'relative' }}>
      <svg
        width="100%"
        height="100%"
        viewBox={`0 0 ${viewBoxWidth} ${viewBoxHeight}`}
        style={{ display: 'block' }}
        preserveAspectRatio="none"
        onMouseMove={handleMouseMove}
        onMouseLeave={handleMouseLeave}
        onClick={handleClick}
      >
        <path d={areaPath} fill="#e0baff" fillOpacity={0.5} stroke="none" />
        <path d={catmullRom2bezier(points)} fill="none" stroke="#e0baff" strokeWidth={2} opacity={0.5} />
        <rect
          x={progressPercent * viewBoxWidth - 1}
          y={0}
          width={2}
          height={viewBoxHeight}
          fill="#fff"
          style={{ pointerEvents: 'none' }}
        />
        {/* Линия и tooltip при наведении */}
        {hoverX !== null && hoverTime !== null && (
          <g>
            <rect x={hoverX - 0.5} y={0} width={1} height={viewBoxHeight} fill="#fff" opacity={0.7} />
          </g>
        )}
      </svg>
      {/* Tooltip над графиком */}
      {hoverX !== null && hoverTime !== null && (
        <>
          {/* Tooltip под курсором */}
          <Box
            sx={{
              position: 'absolute',
              left: `calc(${(hoverX / viewBoxWidth) * 100}% - 18px)` ,
              top: 0,
              pointerEvents: 'none',
              zIndex: 2,
              background: 'rgba(30,30,40,0.95)',
              color: '#fff',
              px: 1,
              py: 0.2,
              borderRadius: 1,
              fontSize: 12,
              textAlign: 'center',
              minWidth: 32,
              boxShadow: '0 2px 8px rgba(0,0,0,0.15)',
              transform: 'translateY(-120%)',
              whiteSpace: 'nowrap',
            }}
          >
            {formatTooltipTime(hoverTime)}
          </Box>
          {/* Левая и правая временные метки на графике */}
          <Box
            sx={{
              position: 'absolute',
              left: 0,
              top: '50%',
              transform: 'translateY(-50%)',
              pointerEvents: 'none',
              zIndex: 2,
              background: 'rgba(30,30,40,0.95)',
              color: '#fff',
              px: 1,
              py: 0.2,
              borderRadius: 1,
              fontSize: 12,
              minWidth: 32,
              boxShadow: '0 2px 8px rgba(0,0,0,0.15)',
              whiteSpace: 'nowrap',
            }}
          >
            {formatTooltipTime(progress)}
          </Box>
          <Box
            sx={{
              position: 'absolute',
              right: 0,
              top: '50%',
              transform: 'translateY(-50%)',
              pointerEvents: 'none',
              zIndex: 2,
              background: 'rgba(30,30,40,0.95)',
              color: '#fff',
              px: 1,
              py: 0.2,
              borderRadius: 1,
              fontSize: 12,
              textAlign: 'right',
              boxShadow: '0 2px 8px rgba(0,0,0,0.15)',
              whiteSpace: 'nowrap',
            }}
          >
            {formatTooltipTime(duration)}
          </Box>
        </>
      )}
    </Box>
  );
};

const PlayerBar = () => {
  const { state, controls } = usePlayerContext();
  const {
    currentTrack,
    isPlaying,
    volume,
    progress,
    duration,
  } = state;
  const {
    togglePlay,
    playNext,
    playPrevious,
    setVolume,
    setProgress,
  } = controls;

  const navigate = useNavigate();
  const { user, isAuthenticated } = useAuthContext();
  const [menuAnchorEl, setMenuAnchorEl] = useState<null | HTMLElement>(null);
  const [playlistMenuAnchorEl, setPlaylistMenuAnchorEl] = useState<null | HTMLElement>(null);
  const [playlists, setPlaylists] = useState<Playlist[]>([]);
  const [addingTrackId, setAddingTrackId] = useState<string | null>(null);
  const [addError, setAddError] = useState<string | null>(null);
  const [coverUrl, setCoverUrl] = useState<string | null>(null);
  const [listenedRanges, setListenedRanges] = useState<[number, number][]>([]);
  const [currentRange, setCurrentRange] = useState<[number, number] | null>(null);
  const [segments, setSegments] = useState<any[]>(mockSegments);

  // Загрузка обложки альбома при изменении трека
  useEffect(() => {
    const loadAlbumCover = async () => {
      if (!isAuthenticated || !currentTrack?.album?.id) {
        setCoverUrl(null);
        return;
      }

      try {
        const response = await apiService.get(`/albums/${currentTrack.album.id}/cover`, { responseType: 'blob' });
        const url = URL.createObjectURL(response);
        setCoverUrl(url);
        return () => URL.revokeObjectURL(url);
      } catch (err) {
        console.error('Failed to load album cover:', err);
        setCoverUrl(null);
      }
    };

    loadAlbumCover();
  }, [currentTrack?.album?.id, isAuthenticated]);

  // Получение сегментов для графика
  useEffect(() => {
    const fetchSegments = async () => {
      if (!isAuthenticated || !currentTrack?.id) {
        setSegments(mockSegments);
        return;
      }

      try {
        const data = await apiService.get(`/tracks/${currentTrack.id}/segments`);
        setSegments(
          Array.isArray(data) && data.length > 0
            ? data.map(seg => ({
                ...seg,
                totalStreams: seg.totalStreams ?? seg.total_streams,
              }))
            : mockSegments
        );
      } catch (err) {
        setSegments(mockSegments);
      }
    };

    fetchSegments();
  }, [currentTrack?.id, isAuthenticated]);

  const fetchPlaylists = async () => {
    if (!isAuthenticated) {
      setPlaylists([]);
      return;
    }

    try {
      const data = await apiService.get('/me/playlists');
      setPlaylists(data);
    } catch (err) {
      setPlaylists([]);
    }
  };

  // Отправка статистики на сервер
  const sendListeningStats = useCallback(async (lastRange: [number, number]) => {
    const allRanges = lastRange ? [...listenedRanges, lastRange] : listenedRanges;

    if (currentTrack && allRanges.length > 0) {
      try {
        await apiService.post(`/tracks/${currentTrack.id}/stats`, {
          track_id: currentTrack.id,
          ranges: allRanges.map(([start, end]) => ({
            start: Math.floor(start),
            end: Math.floor(end),
          })),
        });
      } catch (err) {
        console.error('Failed to send listening stats:', err);
      }
    }
  }, [currentTrack, listenedRanges]);

  // Инициализация range при начале проигрывания
  useEffect(() => {
    if (isPlaying && !currentRange) {
      setCurrentRange([0, 0]);
    }
  }, [isPlaying, currentRange]);

  // Обработка перемотки
  const handleSetProgress = (value: number) => {
    if (isPlaying && currentRange) {
      const [start, _] = currentRange;
      if (progress - start > 2) {
        setListenedRanges(prev => [...prev, [start, progress]]);
      }
      setCurrentRange([value, value]);
    }
    setProgress(value);
  };

  // Обработка завершения трека
  useEffect(() => {
    if (progress >= duration - 0.1 && duration > 0 && currentRange) {
      const [start, _] = currentRange;
      if (duration - start > 2) {
        sendListeningStats([start, duration]);
      }
      setListenedRanges([]);
      setCurrentRange(null);
      playNext();
    }
  }, [progress, duration, currentRange, setListenedRanges, sendListeningStats, playNext]);

  // Обработка переключения трека
  const handlePlayNext = useCallback(() => {
    if (currentRange) {
      const [start, _] = currentRange;
      if (progress - start > 2) {
        sendListeningStats([start, progress]);
      }
    }
    setListenedRanges([]);
    setCurrentRange(null);
    playNext();
  }, [currentRange, progress, sendListeningStats, playNext]);

  const handlePlayPrevious = useCallback(() => {
    if (currentRange) {
      const [start, _] = currentRange;
      if (progress - start > 2) {
        sendListeningStats([start, progress]);
      }
    }
    setListenedRanges([]);
    setCurrentRange(null);
    playPrevious();
  }, [currentRange, progress, sendListeningStats, playPrevious]);

  // Очистка при размонтировании
  useEffect(() => {
    return () => {
      if (currentRange) {
        const [start, _] = currentRange;
        if (progress - start > 2) {
          sendListeningStats([start, progress]);
        }
      }
      setListenedRanges([]);
      setCurrentRange(null);
    };
  }, []);

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    event.stopPropagation();
    setMenuAnchorEl(event.currentTarget);
  };
  const handleMenuClose = () => {
    setMenuAnchorEl(null);
  };

  const handleAddToPlaylistClick = async (e: React.MouseEvent<HTMLElement>) => {
    e.stopPropagation();
    if (!user) {
      navigate('/login', {
        state: {
          from: window.location.pathname,
          message: 'Чтобы добавить трек в плейлист, необходимо войти',
        },
      });
      return;
    }
    await fetchPlaylists();
    setPlaylistMenuAnchorEl(e.currentTarget);
  };

  const handlePlaylistMenuClose = () => {
    setPlaylistMenuAnchorEl(null);
    setAddError(null);
  };

  const handlePlaylistSelect = async (playlistId: string) => {
    if (!currentTrack) return;
    setAddingTrackId(currentTrack.id);
    setAddError(null);
    try {
      await apiService.post(`/playlists/${playlistId}/tracks`, { track_id: currentTrack.id });
      setAddError(null);
      setPlaylistMenuAnchorEl(null);
    } catch (err: any) {
      setAddError('Ошибка при добавлении трека в плейлист');
    } finally {
      setAddingTrackId(null);
    }
  };

  // Проверка, есть ли трек в плейлисте
  const isTrackInPlaylist = (playlist: Playlist, trackId: string) => {
    if (!playlist.tracks) return false;
    return playlist.tracks.some(track => track.id === trackId);
  };

  const toggleMute = useCallback(() => {
    setVolume(volume === 0 ? 1 : 0);
  }, [volume, setVolume]);

  useEffect(() => {
    const handleKeyPress = (event: KeyboardEvent) => {
      // Проверяем, не находится ли фокус в текстовом поле или textarea
      const activeElement = document.activeElement;
      const isInputElement = activeElement instanceof HTMLInputElement || 
                           activeElement instanceof HTMLTextAreaElement;
      
      if (event.code === 'Space' && !event.repeat && !isInputElement) {
        event.preventDefault();
        togglePlay();
      }
    };

    window.addEventListener('keydown', handleKeyPress);
    return () => window.removeEventListener('keydown', handleKeyPress);
  }, [togglePlay]);

  const handleArtistClick = (e: React.MouseEvent, artistId: string) => {
    e.stopPropagation();
    navigate(`/artists/${artistId}`);
  };
  const handleAlbumClick = () => {
    if (currentTrack?.album?.id) {
      navigate(`/albums/${currentTrack.album.id}`);
    }
    handleMenuClose();
  };

  return (
    <Box
      sx={{
        width: '100%',
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'center',
        px: 0,
        gap: 1,
        border: '1px solid rgba(255, 255, 255, 0.1)',
        borderRadius: '12px',
        boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)',
        background: '#181825',
      }}
    >
      <Box sx={{ width: '100%', height: '40%', minHeight: 40, maxHeight: 60, mb: 1 }}>
        <Box sx={{ width: '100%', height: '100%' }}>
          <WaveformStats segments={currentTrack ? segments : mockSegments} progress={progress} duration={duration} setProgress={handleSetProgress} />
        </Box>
        <Slider
          value={progress}
          max={duration}
          onChange={(_, value) => handleSetProgress(value as number)}
          disabled={!currentTrack}
          sx={{ width: '100%', mt: -2 }}
        />
      </Box>
      {/* Управление и информация о треке */}
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, width: '100%' }}>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, minWidth: 200 }}>
          {currentTrack ? (
            <>
              <Box
                component="img"
                src={coverUrl || ''}
                alt={currentTrack.name}
                sx={{
                  width: 56,
                  height: 56,
                  borderRadius: 1,
                  objectFit: 'cover',
                  ml: 1,
                }}
              />
              <Box>
                <Typography variant="subtitle1" noWrap>
                  {currentTrack.name}
                </Typography>
                <Typography variant="body2" color="text.secondary" noWrap>
                  {currentTrack.artists.map((artist, idx) => (
                    <React.Fragment key={artist.id}>
                      {idx > 0 && ', '}
                      <Typography
                        component="span"
                        sx={{
                          color: 'text.secondary',
                          cursor: 'pointer',
                          textDecoration: 'none',
                          '&:hover': { textDecoration: 'underline' },
                        }}
                        onClick={(e) => handleArtistClick(e, artist.id)}
                      >
                        {artist.name}
                      </Typography>
                    </React.Fragment>
                  ))}
                </Typography>
              </Box>
            </>
          ) : (
            <>
              <Box
                sx={{
                  width: 56,
                  height: 56,
                  borderRadius: 1,
                  bgcolor: 'action.hover',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  ml: 1,
                }}
              >
                <MusicNote sx={{ color: 'text.secondary' }} />
              </Box>
              <Box>
                <Typography variant="subtitle1" color="text.secondary">
                  Нет активного трека
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Выберите трек для воспроизведения
                </Typography>
              </Box>
            </>
          )}
        </Box>
        <Box sx={{ flex: 1, display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 1 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <IconButton onClick={handlePlayPrevious} disabled={!currentTrack}>
              <SkipPrevious />
            </IconButton>
            <IconButton onClick={togglePlay} disabled={!currentTrack}>
              {isPlaying ? <Pause /> : <PlayArrow />}
            </IconButton>
            <IconButton onClick={handlePlayNext} disabled={!currentTrack}>
              <SkipNext />
            </IconButton>
          </Box>
        </Box>
        {/* Кнопки + и 3 точки слева от микшера громкости */}
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <IconButton onClick={handleMenuOpen} disabled={!currentTrack}>
            <MoreVert />
          </IconButton>
          <Menu
            anchorEl={menuAnchorEl}
            open={Boolean(menuAnchorEl)}
            onClose={handleMenuClose}
            onClick={e => e.stopPropagation()}
          >
            <MenuItem onClick={handleAlbumClick}>Перейти к альбому</MenuItem>
          </Menu>
        </Box>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, minWidth: 200 }}>
          <IconButton onClick={toggleMute} disabled={!currentTrack}>
            {volume === 0 ? <VolumeOff /> : <VolumeUp />}
          </IconButton>
          <Slider
            value={volume}
            min={0}
            max={1}
            step={0.01}
            onChange={(_, value) => setVolume(value as number)}
            disabled={!currentTrack}
            sx={{ width: 100 }}
          />
        </Box>
      </Box>
    </Box>
  );
};

export default PlayerBar; 