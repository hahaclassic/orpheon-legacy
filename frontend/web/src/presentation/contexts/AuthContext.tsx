import { createContext, useContext } from 'react';
import { useAuth } from '../hooks/useAuth';
import { useLocation } from 'react-router-dom';

// Separate interfaces for better interface segregation
interface User {
  id: string;
  name: string;
  registration_date: string;
  birth_date: string;
  access_lvl: number;
}

interface AuthState {
  isAuthenticated: boolean;
  isAdmin: boolean;
  isLoading: boolean;
  user: User | null;
}

interface AuthActions {
  login: (login: string, password: string) => Promise<User>;
  register: (login: string, password: string) => Promise<User>;
  logout: () => Promise<void>;
  changePassword: (oldPassword: string, newPassword: string) => Promise<void>;
}

interface AuthContextType {
  isAuthenticated: boolean;
  isAdmin: boolean;
  user: {
    id: string;
    name: string;
    registration_date: string;
    birth_date: string;
    access_lvl: number;
  } | null;
  isLoading: boolean;
  login: (login: string, password: string) => Promise<any>;
  register: (login: string, password: string) => Promise<any>;
  logout: () => Promise<void>;
  changePassword: (oldPassword: string, newPassword: string) => Promise<any>;
  updateUser: (user: {
    id: string;
    name: string;
    registration_date: string;
    birth_date: string;
    access_lvl: number;
  }) => Promise<any>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuthContext = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuthContext must be used within an AuthProvider');
  }
  return context;
};

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const location = useLocation();
  const auth = useAuth(location.pathname);
  return <AuthContext.Provider value={auth}>{children}</AuthContext.Provider>;
}; 