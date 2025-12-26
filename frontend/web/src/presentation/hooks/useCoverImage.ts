import { useState, useEffect } from 'react';
import { apiService } from '../services/api';

export const useCoverImage = (type: 'album' | 'playlist', id: string | number) => {
  const [coverUrl, setCoverUrl] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchCover = async () => {
      try {
        setLoading(true);
        setError(null);
        const response = await apiService.get(`/${type}s/${id}/cover`, {
          responseType: 'blob'
        });
        const imageUrl = URL.createObjectURL(response);
        setCoverUrl(imageUrl);
      } catch (err) {
        console.error(`Error fetching ${type} cover:`, err);
        setError('Failed to load cover image');
      } finally {
        setLoading(false);
      }
    };

    if (id) {
      fetchCover();
    }

    return () => {
      if (coverUrl) {
        URL.revokeObjectURL(coverUrl);
      }
    };
  }, [type, id]);

  return { coverUrl, loading, error };
}; 