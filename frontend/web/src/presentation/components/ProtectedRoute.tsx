import { Navigate, useLocation } from 'react-router-dom';
import { useAuthContext } from '../contexts/AuthContext';
import type { ReactNode } from 'react';
import { Box, CircularProgress } from '@mui/material';

interface ProtectedRouteProps {
  children: ReactNode;
}

// List of public routes that don't require authentication
const PUBLIC_ROUTES = [
  '/',                    // Home
  '/search',             // Search
  '/login',              // Login
  '/register',           // Register
  '/artists',            // Artists list
  '/artists/',           // Artist profile (with ID)
  '/albums',             // Albums list
  '/albums/',            // Album details (with ID)
  '/playlists',          // Public playlists list
  '/playlists/',         // Public playlist details (with ID)
  '/tracks/',            // Track details (with ID)
];

const ProtectedRoute = ({ children }: ProtectedRouteProps) => {
  const { isAuthenticated, isLoading } = useAuthContext();
  const location = useLocation();

  console.log('[ProtectedRoute] Checking authentication:', {
    isAuthenticated,
    isLoading,
    currentPath: location.pathname,
  });

  // Check if the current path matches any public route pattern
  const isPublicRoute = PUBLIC_ROUTES.some(route => {
    // For routes with trailing slash, check if the path starts with the route
    if (route.endsWith('/')) {
      return location.pathname.startsWith(route);
    }
    // For exact routes, check for exact match
    return location.pathname === route;
  });

  // Allow access to public routes without authentication
  if (isPublicRoute) {
    console.log('[ProtectedRoute] Public route, allowing access');
    return <>{children}</>;
  }

  // Если это защищенный маршрут и идет загрузка, показываем загрузчик
  // но не делаем редирект
  if (isLoading) {
    console.log('[ProtectedRoute] Loading authentication state...');
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        <CircularProgress />
      </Box>
    );
  }

  // Только после завершения загрузки проверяем аутентификацию
  if (!isAuthenticated) {
    console.log('[ProtectedRoute] User not authenticated, redirecting to login');
    // Сохраняем текущий путь для возврата после логина
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  console.log('[ProtectedRoute] User authenticated, rendering protected content');
  return <>{children}</>;
};

export default ProtectedRoute; 