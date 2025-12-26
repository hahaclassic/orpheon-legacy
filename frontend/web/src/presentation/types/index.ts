export interface User {
  id: number;
  email: string;
  username: string;
  avatar?: string;
}

export interface Track {
  id: string;
  name: string;
  duration: number;
  track_number: number;
  coverUrl?: string;
  audioUrl: string;
  artists: Array<{
    id: string;
    name: string;
  }>;
  album_id: string;
  album: {
    id: string;
    title: string;
    label: string;
    license_id: string;
    release_date: string;
  };
  total_streams?: number;
  license?: {
    id: string;
    title: string;
    description: string;
    url: string;
  };
}

export interface Playlist {
  id: string;
  name: string;
  coverImage?: string;
  trackCount: number;
  tracks: Track[];
}

export interface ApiResponse<T> {
  data: T;
  error?: string;
} 