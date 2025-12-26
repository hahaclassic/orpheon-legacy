import api from './api';

export interface Genre {
  id: string;
  title: string;
}

export interface Artist {
  id: string;
  name: string;
  description: string;
  country: string;
}

export interface Album {
  id: string;
  title: string;
  label: string;
  license_id: string;
  release_date: string;
  artists: Artist[];
  genres: Genre[];
  imageUrl?: string;
}

export interface Track {
  id: string;
  title: string;
  duration: string;
  trackNumber: number;
}

export const albumService = {
  getAlbum: async (id: string): Promise<Album> => {
    const response = await api.get(`/albums/${id}`);
    return response.data;
  },

  getAlbums: async (): Promise<Album[]> => {
    const response = await api.get('/albums');
    return response.data;
  },
}; 