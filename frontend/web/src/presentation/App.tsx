import { BrowserRouter, Routes, Route, Navigate, Outlet, useLocation } from 'react-router-dom';
import { ThemeProvider } from '@mui/material/styles';
import { LocalizationProvider } from '@mui/x-date-pickers';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import { AuthProvider } from './contexts/AuthContext';
import { PlayerProvider } from './contexts/PlayerContext';
import ErrorBoundary from './components/ErrorBoundary';
import Layout from './components/Layout';
import ProtectedRoute from './components/ProtectedRoute';
import AdminRoute from './components/AdminRoute';
import Home from './pages/Home';
import Search from './pages/Search';
import Profile from './pages/Profile';
import Me from './pages/Me';
import GenreList from './pages/Admin/GenreList';
import LicenseList from './pages/Admin/LicenseList';
import ArtistList from './pages/Admin/ArtistList';
import AlbumList from './pages/Admin/AlbumList';
import { theme } from './theme';

const ProtectedRoutes = () => {
  return (
    <ProtectedRoute>
      <Routes>
        <Route path="/profile" element={<Profile />} />
        <Route path="/me" element={<Me />} />
        
        <Route path="/admin" element={<AdminRoute><Outlet /></AdminRoute>}>
          <Route index element={<Navigate to="genres" replace />} />
          <Route path="genres" element={<GenreList />} />
          <Route path="licenses" element={<LicenseList />} />
          <Route path="artists" element={<ArtistList />} />
          <Route path="albums" element={<AlbumList />} />
        </Route>
      </Routes>
    </ProtectedRoute>
  );
};

const AuthProviderWithLocation = ({ children }: { children: React.ReactNode }) => {
  const location = useLocation();
  return <AuthProvider initialPath={location.pathname}>{children}</AuthProvider>;
};

const App = () => {
  return (
    <ErrorBoundary>
      <ThemeProvider theme={theme}>
        <LocalizationProvider dateAdapter={AdapterDateFns}>
          <BrowserRouter future={{ v7_startTransition: true, v7_relativeSplatPath: true }}>
            <AuthProviderWithLocation>
              <PlayerProvider>
                <Routes>
                  <Route path="/login" element={<Navigate to="/auth/login" replace />} />
                  <Route path="/register" element={<Navigate to="/auth/register" replace />} />
                  
                  <Route element={<Layout />}>
                    {/* Public routes */}
                    <Route path="/" element={<Home />} />
                    <Route path="/search" element={<Search />} />
                    <Route path="/artists" element={<ArtistList />} />
                    <Route path="/artists/:id" element={<ArtistList />} />
                    <Route path="/albums" element={<AlbumList />} />
                    <Route path="/albums/:id" element={<AlbumList />} />
                    <Route path="/playlists" element={<AlbumList />} />
                    <Route path="/playlists/:id" element={<AlbumList />} />
                    <Route path="/tracks/:id" element={<AlbumList />} />
                    
                    {/* Protected routes */}
                    <Route path="/profile/*" element={<ProtectedRoutes />} />
                    <Route path="/me/*" element={<ProtectedRoutes />} />
                    <Route path="/admin/*" element={<ProtectedRoutes />} />
                  </Route>
                </Routes>
              </PlayerProvider>
            </AuthProviderWithLocation>
          </BrowserRouter>
        </LocalizationProvider>
      </ThemeProvider>
    </ErrorBoundary>
  );
};

export default App; 