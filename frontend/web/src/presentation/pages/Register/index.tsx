import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  TextField,
  Button,
  Typography,
  Link,
  Paper,
  styled,
} from '@mui/material';
import { useAuthContext } from '../../contexts/AuthContext';
import ErrorBoundary from '../../components/ErrorBoundary';

const RegisterContainer = styled(Box)({
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  minHeight: '100vh',
  padding: '24px',
  width: '100vw',
  boxSizing: 'border-box',
});

const RegisterForm = styled(Paper)({
  padding: '32px',
  width: '100%',
  maxWidth: '400px',
  display: 'flex',
  flexDirection: 'column',
  gap: '24px',
});

const Register = () => {
  const navigate = useNavigate();
  const { register } = useAuthContext();
  const [formData, setFormData] = useState({
    login: '',
    password: '',
    confirmPassword: '',
  });
  const [error, setError] = useState('');

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (formData.password !== formData.confirmPassword) {
      setError('Passwords do not match');
      return;
    }

    try {
      await register(formData.login, formData.password);
      navigate('/');
    } catch (err) {
      setError('Registration failed');
    }
  };

  return (
    <ErrorBoundary>
      <RegisterContainer>
        <RegisterForm elevation={3}>
          <Typography variant="h4" align="center" gutterBottom>
            Create Account
          </Typography>
          <Typography variant="body1" align="center" color="text.secondary" paragraph>
            Sign up to start using Orpheon
          </Typography>

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
            <TextField
              fullWidth
              label="Confirm Password"
              type="password"
              name="confirmPassword"
              value={formData.confirmPassword}
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
              type="submit"
              fullWidth
              variant="contained"
              size="large"
              sx={{ mt: 3 }}
            >
              Sign Up
            </Button>
          </form>

          <Box sx={{ textAlign: 'center', mt: 2 }}>
            <Typography variant="body2" color="text.secondary">
              Already have an account?{' '}
              <Link href="/login" underline="hover">
                Sign in
              </Link>
            </Typography>
          </Box>
        </RegisterForm>
      </RegisterContainer>
    </ErrorBoundary>
  );
};

export default Register; 