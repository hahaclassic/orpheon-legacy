import { useState, useRef, useCallback, useEffect } from 'react';
import type { Track } from '../types';
import api from '../../core/infrastructure/services/api';

interface PlayerState {
  currentTrack: Track | null;
  isPlaying: boolean;
  volume: number;
  progress: number;
  duration: number;
  queue: Track[];
}

interface QueueManager {
  addToQueue: (tracks: Track[]) => void;
  clearQueue: () => void;
  getNextTrack: () => Track | null;
  getPreviousTrack: () => Track | null;
  getQueue: () => Track[];
  setContextTracks: (tracks: Track[]) => void;
  getContextTracks: () => Track[];
  setCurrentIndex: (index: number) => void;
  getCurrentIndex: () => number;
}

interface AudioPlayer {
  play: () => Promise<void>;
  pause: () => void;
  setVolume: (volume: number) => void;
  setProgress: (progress: number) => void;
  addEventListener: (event: string, callback: () => void) => void;
  removeEventListener: (event: string, callback: () => void) => void;
  getCurrentTime: () => number;
  getDuration: () => number;
  setSrc: (src: string) => void;
}

class TrackQueue implements QueueManager {
  private queue: Track[] = [];
  private history: Track[] = [];
  private contextTracks: Track[] = [];
  private currentIndex: number = -1;

  addToQueue(tracks: Track[]): void {
    this.queue.push(...tracks);
  }

  clearQueue(): void {
    this.queue = [];
    this.history = [];
    this.currentIndex = -1;
  }

  setContextTracks(tracks: Track[]): void {
    this.contextTracks = tracks;
  }

  getContextTracks(): Track[] {
    return [...this.contextTracks];
  }

  setCurrentIndex(index: number): void {
    this.currentIndex = index;
  }

  getCurrentIndex(): number {
    return this.currentIndex;
  }

  getNextTrack(): Track | null {
    if (this.currentIndex < this.contextTracks.length - 1) {
      this.currentIndex++;
      return this.contextTracks[this.currentIndex];
    }
    return null;
  }

  getPreviousTrack(): Track | null {
    if (this.currentIndex > 0) {
      this.currentIndex--;
      return this.contextTracks[this.currentIndex];
    }
    return null;
  }

  getQueue(): Track[] {
    return [...this.queue];
  }
}

// Single Responsibility Principle: отдельный класс для работы с аудио
class HTMLAudioPlayer implements AudioPlayer {
  private audio: HTMLAudioElement;

  constructor() {
    this.audio = new Audio();
  }

  async play(): Promise<void> {
    await this.audio.play();
  }

  pause(): void {
    this.audio.pause();
  }

  setVolume(volume: number): void {
    this.audio.volume = volume;
  }

  setProgress(progress: number): void {
    this.audio.currentTime = progress;
  }

  addEventListener(event: string, callback: () => void): void {
    this.audio.addEventListener(event, callback);
  }

  removeEventListener(event: string, callback: () => void): void {
    this.audio.removeEventListener(event, callback);
  }

  getCurrentTime(): number {
    return this.audio.currentTime;
  }

  getDuration(): number {
    return this.audio.duration || 0;
  }

  setSrc(src: string): void {
    this.audio.src = src;
  }
}

const initialState: PlayerState = {
  currentTrack: null,
  isPlaying: false,
  volume: 1,
  progress: 0,
  duration: 0,
  queue: [],
};

// Функция для сохранения состояния в localStorage
const saveStateToStorage = (state: PlayerState) => {
  try {
    localStorage.setItem('playerState', JSON.stringify({
      ...state,
      currentTrack: state.currentTrack ? {
        id: state.currentTrack.id,
        name: state.currentTrack.name,
        duration: state.currentTrack.duration,
        artists: state.currentTrack.artists,
        coverUrl: state.currentTrack.coverUrl,
      } : null,
      queue: state.queue.map(track => ({
        id: track.id,
        name: track.name,
        duration: track.duration,
        artists: track.artists,
        coverUrl: track.coverUrl,
      })),
    }));
  } catch (error) {
    console.error('Failed to save player state:', error);
  }
};

// Функция для загрузки состояния из localStorage
const loadStateFromStorage = (): PlayerState => {
  try {
    const savedState = localStorage.getItem('playerState');
    if (savedState) {
      return JSON.parse(savedState);
    }
  } catch (error) {
    console.error('Failed to load player state:', error);
  }
  return initialState;
};

export const usePlayer = () => {
  const [state, setState] = useState<PlayerState>(loadStateFromStorage);
  const playerRef = useRef<AudioPlayer>(new HTMLAudioPlayer());
  const queueManagerRef = useRef<QueueManager>(new TrackQueue());

  const updateState = useCallback((updates: Partial<PlayerState>) => {
    setState(prev => {
      const newState = { ...prev, ...updates };
      saveStateToStorage(newState);
      return newState;
    });
  }, []);

  // Восстанавливаем состояние при инициализации
  useEffect(() => {
    if (state.currentTrack) {
      const audioUrl = `${api.defaults.baseURL}/tracks/${state.currentTrack.id}/audio`;
      playerRef.current.setSrc(audioUrl);
      playerRef.current.setVolume(state.volume);
      playerRef.current.setProgress(state.progress);
      if (state.isPlaying) {
        playerRef.current.play();
      }
    }
  }, []);

  const setTrack = useCallback((track: Track) => {
    const audioUrl = `${api.defaults.baseURL}/tracks/${track.id}/audio`;
    playerRef.current.setSrc(audioUrl);
    updateState({
      currentTrack: track,
      isPlaying: true,
      progress: 0,
      duration: 0,
    });
    playerRef.current.play();
  }, [updateState]);

  const togglePlay = useCallback(() => {
    if (!state.currentTrack) return;

    if (state.isPlaying) {
      playerRef.current.pause();
      updateState({ isPlaying: false });
    } else {
      playerRef.current.play();
      updateState({ isPlaying: true });
    }
  }, [state.currentTrack, state.isPlaying, updateState]);

  const setVolume = useCallback((volume: number) => {
    playerRef.current.setVolume(volume);
    updateState({ volume });
  }, [updateState]);

  const setProgress = useCallback((progress: number) => {
    playerRef.current.setProgress(progress);
    updateState({ progress });
  }, [updateState]);

  const playNext = useCallback(() => {
    const nextTrack = queueManagerRef.current.getNextTrack();
    if (nextTrack) {
      setTrack(nextTrack);
    }
  }, [setTrack]);

  const playPrevious = useCallback(() => {
    const previousTrack = queueManagerRef.current.getPreviousTrack();
    if (previousTrack) {
      setTrack(previousTrack);
    }
  }, [setTrack]);

  const startPlayback = useCallback((track: Track, contextTracks: Track[]) => {
    // Сохраняем все треки из контекста
    queueManagerRef.current.setContextTracks(contextTracks);
    
    // Находим индекс выбранного трека
    const trackIndex = contextTracks.findIndex(t => t.id === track.id);
    if (trackIndex === -1) return;

    // Устанавливаем текущий индекс
    queueManagerRef.current.setCurrentIndex(trackIndex);
    
    // Устанавливаем выбранный трек как текущий
    setTrack(track);
  }, [setTrack]);

  useEffect(() => {
    const handleTimeUpdate = () => {
      const currentTime = playerRef.current.getCurrentTime();
      updateState({ progress: currentTime });
    };

    const handleLoadedMetadata = () => {
      const duration = playerRef.current.getDuration();
      updateState({ duration });
    };

    const handleEnded = () => {
      playNext();
    };

    playerRef.current.addEventListener('timeupdate', handleTimeUpdate);
    playerRef.current.addEventListener('loadedmetadata', handleLoadedMetadata);
    playerRef.current.addEventListener('ended', handleEnded);

    return () => {
      playerRef.current.removeEventListener('timeupdate', handleTimeUpdate);
      playerRef.current.removeEventListener('loadedmetadata', handleLoadedMetadata);
      playerRef.current.removeEventListener('ended', handleEnded);
    };
  }, [updateState, playNext]);

  return {
    state,
    controls: {
      startPlayback,
      togglePlay,
      setVolume,
      setProgress,
      playNext,
      playPrevious,
    },
  };
};

export default usePlayer; 