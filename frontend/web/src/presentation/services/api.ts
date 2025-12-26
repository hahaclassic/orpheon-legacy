import axios from 'axios';
import type { AxiosResponse, AxiosInstance, AxiosRequestConfig } from 'axios';

const API_URL = import.meta.env.VITE_API_URL;

interface User {
  id: string;
  name: string;
  registration_date: string;
  birth_date: string;
  access_lvl: number;
}

interface Genre {
  id: string;
  title: string;
}

interface License {
  id: string;
  title: string;
  description: string;
}

interface Artist {
  id: string;
  name: string;
  description: string;
  country: string;
}

interface Playlist {
  id: string;
  title: string;
  // Add any other necessary properties for a playlist
}

// Создаем тип для apiService, который возвращает данные напрямую
type ApiService = {
  get: <T = any>(url: string, config?: AxiosRequestConfig) => Promise<T>;
  post: <T = any>(url: string, data?: any, config?: AxiosRequestConfig) => Promise<T>;
  put: <T = any>(url: string, data?: any, config?: AxiosRequestConfig) => Promise<T>;
  delete: <T = any>(url: string, config?: AxiosRequestConfig) => Promise<T>;
  interceptors: {
    request: {
      use: (fulfilled: (config: AxiosRequestConfig) => AxiosRequestConfig, rejected?: (error: any) => any) => number;
    };
    response: {
      use: (fulfilled: (response: AxiosResponse) => any, rejected?: (error: any) => any) => number;
    };
  };
  (config: AxiosRequestConfig): Promise<any>;
};

const apiService = axios.create({
  baseURL: API_URL,
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
  },
}) as unknown as ApiService;

// Request interceptor
apiService.interceptors.request.use(
  (config) => {
    console.log('[API] Request:', {
      url: config.url,
      method: config.method,
      headers: config.headers,
      withCredentials: config.withCredentials,
    });
    return config;
  },
  (error) => {
    console.error('[API] Request Error:', error);
    return Promise.reject(error);
  }
);

// Response interceptor
apiService.interceptors.response.use(
  (response) => {
    console.log('[API] Response:', {
      url: response.config.url,
      status: response.status,
      data: response.data,
      headers: response.headers,
    });
    return response.data;
  },
  async (error) => {
    console.error('[API] Error:', {
      url: error.config?.url,
      status: error.response?.status,
      data: error.response?.data,
      headers: error.response?.headers,
    });

    if (error.response?.status === 401) {
      // Если получили 401, значит бэкенд уже попытался обновить токен и не смог
      console.error('[API] Authentication failed:', error);
      localStorage.removeItem('authState');
      window.location.href = '/login';
    }

    return Promise.reject(error);
  }
);

// API methods
export const api = {
  register: (login: string, password: string): Promise<User> =>
    apiService.post('/auth/register', { login, password }),

  login: (login: string, password: string): Promise<User> =>
    apiService.post('/auth/login', { login, password }),

  logout: (): Promise<void> => apiService.post('/auth/logout'),

  changePassword: (oldPassword: string, newPassword: string): Promise<void> =>
    apiService.post('/auth/password/update', { old: oldPassword, new: newPassword }),

  getMe: (): Promise<User> => apiService.get('/me'),

  updateUser: (user: User): Promise<User> => apiService.put('/me', user),

  // User endpoints
  getUser: (id: string): Promise<User> => apiService.get(`/users/${id}`),
  getUserPlaylists: (id: string): Promise<Playlist[]> => apiService.get(`/users/${id}/playlists`),
  getUserFavorites: (id: string): Promise<Playlist[]> => apiService.get(`/users/${id}/favorites`),

  // Tracks
  getTracks: () => apiService.get('/tracks'),
  getTrack: (id: string) => apiService.get(`/tracks/${id}`),
  createTrack: (data: FormData) =>
    apiService.post('/tracks', data, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    }),
  updateTrack: (id: string, data: FormData) =>
    apiService.put(`/tracks/${id}`, data, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    }),
  deleteTrack: (id: string) => apiService.delete(`/tracks/${id}`),

  // Albums
  getAlbums: () => apiService.get('/albums'),
  getAlbum: (id: string) => apiService.get(`/albums/${id}`),
  createAlbum: (data: FormData) =>
    apiService.post('/albums', data, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    }),
  updateAlbum: (id: string, data: FormData) =>
    apiService.put(`/albums/${id}`, data, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    }),
  deleteAlbum: (id: string) => apiService.delete(`/albums/${id}`),

  // Artists
  getArtists: () => apiService.get<Artist[]>('/artists'),
  getArtist: (id: string) => apiService.get(`/artists/${id}`),
  createArtist: (data: Omit<Artist, 'id'>) => apiService.post<Artist>('/artists', data),
  updateArtist: (id: string, data: Omit<Artist, 'id'>) => apiService.put<Artist>(`/artists/${id}`, data),
  deleteArtist: (id: string) => apiService.delete(`/artists/${id}`),
  getArtistTracks: (id: string) => apiService.get(`/artists/${id}/tracks`),
  getArtistAlbums: (id: string) => apiService.get(`/artists/${id}/albums`),

  // Genres
  getGenres: (): Promise<Genre[]> => apiService.get('/genres'),
  getGenre: (id: string): Promise<Genre> => apiService.get(`/genres/${id}`),
  createGenre: (data: { title: string }): Promise<Genre> =>
    apiService.post('/genres', data),
  updateGenre: (id: string, data: { title: string }): Promise<Genre> =>
    apiService.put(`/genres/${id}`, data),
  deleteGenre: (id: string): Promise<void> => apiService.delete(`/genres/${id}`),

  // Licenses
  getLicenses: (): Promise<License[]> => apiService.get('/licenses'),
  getLicense: (id: string): Promise<License> => apiService.get(`/licenses/${id}`),
  createLicense: (data: { title: string; description: string }): Promise<License> =>
    apiService.post('/licenses', data),
  updateLicense: (id: string, data: { title: string; description: string }): Promise<License> =>
    apiService.put(`/licenses/${id}`, data),
  deleteLicense: (id: string): Promise<void> => apiService.delete(`/licenses/${id}`),
};

export { apiService }; 