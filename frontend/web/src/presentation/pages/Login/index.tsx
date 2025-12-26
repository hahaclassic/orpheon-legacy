import { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import {
  Box,
  Typography,
  Link,
  Paper,
  TextField,
  Button,
  styled,
} from '@mui/material';
import { useAuthContext } from '../../contexts/AuthContext';
import ErrorBoundary from '../../components/ErrorBoundary';

const LoginContainer = styled(Box)({
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  minHeight: '100vh',
  padding: '24px',
  width: '100vw',
  boxSizing: 'border-box',
});

const LoginForm = styled(Paper)({
  padding: '32px',
  width: '100%',
  maxWidth: '400px',
  display: 'flex',
  flexDirection: 'column',
  gap: '24px',
});

const Login = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { login } = useAuthContext();
  const [formData, setFormData] = useState({
    login: '',
    password: '',
  });
  const [error, setError] = useState('');

  const message = (location.state as any)?.message;

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    try {
      await login(formData.login, formData.password);
      const from = (location.state as any)?.from || '/';
      navigate(from, { replace: true });
    } catch (err) {
      setError('Invalid login or password');
    }
  };

  return (
    <ErrorBoundary>
      <LoginContainer>
        <LoginForm elevation={3}>
          <Typography variant="h4" align="center" gutterBottom>
            Welcome Back
          </Typography>
          {message ? (
            <Typography variant="body1" align="center" color="text.secondary" paragraph>
              {message}
            </Typography>
          ) : (
            <Typography variant="body1" align="center" color="text.secondary" paragraph>
              Sign in to continue to Orpheon
            </Typography>
          )}

          <form onSubmit={handleSubmit}>
            <TextField
              fullWidth
              label="Login"
              name="login"
              value={formData.login}
              onChange={handleChange}
              required
              margin="normal"
            />
            <TextField
              fullWidth
              label="Password"
              type="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              required
              margin="normal"
            />
            {error && (
              <Typography color="error" align="center" sx={{ mt: 2 }}>
                {error}
              </Typography>
            )}
            <Button
              fullWidth
              variant="contained"
              size="large"
              sx={{ mt: 3 }}
              type="submit"
            >
              Sign In
            </Button>
          </form>

          <Box sx={{ textAlign: 'center', mt: 2 }}>
            <Typography variant="body2" color="text.secondary">
              Don't have an account?{' '}
              <Link href="/register" underline="hover">
                Sign up
              </Link>
            </Typography>
          </Box>
        </LoginForm>
      </LoginContainer>
    </ErrorBoundary>
  );
};

export default Login; 