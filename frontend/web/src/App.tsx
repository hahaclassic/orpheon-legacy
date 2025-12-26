import { BrowserRouter, Routes, Route, Navigate, Outlet } from 'react-router-dom';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { theme } from './presentation/styles/theme';
import Layout from './presentation/components/Layout';
import ProtectedRoute from './presentation/components/ProtectedRoute';
import AdminRoute from './presentation/components/AdminRoute';
import Home from './presentation/pages/Home';
import Library from './presentation/pages/Library';
import Search from './presentation/pages/Search';
import Login from './presentation/pages/Login';
import Register from './presentation/pages/Register';
import { Profile } from './presentation/pages/Profile';
import UserProfile from './presentation/pages/UserProfile';
import PlaylistPage from './presentation/pages/Playlist';
import ArtistPage from './presentation/pages/Artist';
import AlbumPage from './presentation/pages/Album';
import TrackPage from './presentation/pages/Track';
import { AdminPanel } from './presentation/pages/Admin/AdminPanel';
import GenreList from './presentation/pages/Admin/GenreList';
import LicenseList from './presentation/pages/Admin/LicenseList';
import ArtistList from './presentation/pages/Admin/ArtistList';
import AlbumList from './presentation/pages/Admin/AlbumList';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import { AuthProvider } from './presentation/contexts/AuthContext';
import { PlayerProvider } from './presentation/contexts/PlayerContext';

const App = () => {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <BrowserRouter>
        <AuthProvider>
          <PlayerProvider>
            <LocalizationProvider dateAdapter={AdapterDateFns}>
              <Routes>
                {/* Auth routes - без Layout */}
                <Route path="/login" element={<Login />} />
                <Route path="/register" element={<Register />} />

                {/* Public routes - с Layout, но без ProtectedRoute */}
                <Route element={<Layout><Outlet /></Layout>}>
                  <Route path="/" element={<Home />} />
                  <Route path="/search" element={<Search />} />
                  <Route path="/artists/:id" element={<ArtistPage />} />
                  <Route path="/albums/:id" element={<AlbumPage />} />
                  <Route path="/tracks/:id" element={<TrackPage />} />
                  <Route path="/playlists/:id" element={<PlaylistPage />} />
                </Route>

                {/* Protected routes - с Layout и ProtectedRoute */}
                <Route element={<Layout><ProtectedRoute><Outlet /></ProtectedRoute></Layout>}>
                  <Route path="/library" element={<Library />} />
                  <Route path="/me" element={<Profile />} />
                  <Route path="/users/:id" element={<UserProfile />} />
                  
                  {/* Admin routes */}
                  <Route path="/admin" element={<AdminRoute><Outlet /></AdminRoute>}>
                    <Route index element={<AdminPanel />} />
                    <Route path="genres" element={<GenreList />} />
                    <Route path="licenses" element={<LicenseList />} />
                    <Route path="artists" element={<ArtistList />} />
                    <Route path="albums" element={<AlbumList />} />
                  </Route>
                </Route>
              </Routes>
            </LocalizationProvider>
          </PlayerProvider>
        </AuthProvider>
      </BrowserRouter>
    </ThemeProvider>
  );
};

export default App;
