import React, { createContext, useContext, ReactNode } from 'react';
import type { Track } from '../types';
import usePlayer from '../hooks/usePlayer';

// Separate interfaces for better interface segregation
interface PlayerState {
  currentTrack: Track | null;
  isPlaying: boolean;
  volume: number;
  progress: number;
  duration: number;
  playlist: Track[];
  currentIndex: number;
}

interface PlayerControls {
  startPlayback: (track: Track, remainingTracks: Track[]) => void;
  togglePlay: () => void;
  setVolume: (volume: number) => void;
  setProgress: (progress: number) => void;
  playNext: () => void;
  playPrevious: () => void;
}

// Create separate contexts for state and controls
const PlayerStateContext = createContext<PlayerState | null>(null);
const PlayerControlsContext = createContext<PlayerControls | null>(null);

// Provider component that combines both contexts
export const PlayerProvider = ({ children }: { children: ReactNode }) => {
  const { state, controls } = usePlayer();

  return (
    <PlayerStateContext.Provider value={state}>
      <PlayerControlsContext.Provider value={controls}>
        {children}
      </PlayerControlsContext.Provider>
    </PlayerStateContext.Provider>
  );
};

// Custom hooks for accessing specific parts of the player
export const usePlayerState = () => {
  const context = useContext(PlayerStateContext);
  if (!context) {
    throw new Error('usePlayerState must be used within a PlayerProvider');
  }
  return context;
};

export const usePlayerControls = () => {
  const context = useContext(PlayerControlsContext);
  if (!context) {
    throw new Error('usePlayerControls must be used within a PlayerProvider');
  }
  return context;
};

// Main hook that combines both state and controls
export const usePlayerContext = () => {
  const state = usePlayerState();
  const controls = usePlayerControls();

  return {
    state,
    controls,
    // Convenience getters
    currentTrack: state.currentTrack,
    isPlaying: state.isPlaying,
    volume: state.volume,
    progress: state.progress,
    duration: state.duration,
    playlist: state.playlist,
    currentIndex: state.currentIndex,
  };
}; 