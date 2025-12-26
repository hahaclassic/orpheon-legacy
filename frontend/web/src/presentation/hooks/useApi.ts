import { useState, useCallback } from 'react';
import { apiService } from '../services/api';
import type { Playlist, Track, ApiResponse, User } from '../types';

export const useApi = () => {
  const [data, setData] = useState<unknown>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const execute = useCallback(async (promise: Promise<unknown>) => {
    try {
      setLoading(true);
      setError(null);
      const result = await promise;
      setData(result);
      return result;
    } catch (err) {
      setError(err instanceof Error ? err : new Error('An error occurred'));
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const reset = useCallback(() => {
    setData(null);
    setError(null);
  }, []);

  const getPlaylists = useCallback(async (): Promise<Playlist[]> => {
    const response = await execute(apiService.get('/playlists'));
    return (response as ApiResponse<Playlist[]>).data;
  }, [execute]);

  const getLikedTracks = useCallback(async (): Promise<Track[]> => {
    const response = await execute(apiService.get('/tracks/liked'));
    return (response as ApiResponse<Track[]>).data;
  }, [execute]);

  const createPlaylist = useCallback(async (name: string): Promise<Playlist> => {
    const response = await execute(apiService.post('/playlists', { name }));
    return (response as ApiResponse<Playlist>).data;
  }, [execute]);

  const updatePlaylist = useCallback(async (id: number, data: Partial<Playlist>): Promise<Playlist> => {
    const response = await execute(apiService.put(`/playlists/${id}`, data));
    return (response as ApiResponse<Playlist>).data;
  }, [execute]);

  const deletePlaylist = useCallback(async (id: number): Promise<void> => {
    await execute(apiService.delete(`/playlists/${id}`));
  }, [execute]);

  const likeTrack = useCallback(async (trackId: number): Promise<void> => {
    await execute(apiService.post(`/tracks/${trackId}/like`));
  }, [execute]);

  const unlikeTrack = useCallback(async (trackId: number): Promise<void> => {
    await execute(apiService.delete(`/tracks/${trackId}/like`));
  }, [execute]);

  const updateProfile = useCallback(async (data: { username: string; email: string }): Promise<User> => {
    const response = await execute(apiService.put('/users/profile', data));
    return (response as ApiResponse<User>).data;
  }, [execute]);

  return {
    execute,
    reset,
    data,
    loading,
    error,
    getPlaylists,
    getLikedTracks,
    createPlaylist,
    updatePlaylist,
    deletePlaylist,
    likeTrack,
    unlikeTrack,
    updateProfile,
  };
}; 