import { 
  Box, 
  List, 
  ListItem, 
  ListItemIcon, 
  ListItemText, 
  styled, 
  IconButton, 
  Menu, 
  MenuItem, 
  Typography,
  Avatar,
  Tooltip,
  Divider
} from '@mui/material';
import { 
  Home, 
  LibraryMusic, 
  Search, 
  Logout, 
  Person, 
  AdminPanelSettings,
  ChevronLeft,
  ChevronRight
} from '@mui/icons-material';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuthContext } from '../../contexts/AuthContext';
import { useState } from 'react';

const SidebarContainer = styled(Box, {
  shouldForwardProp: (prop) => prop !== 'isCollapsed'
})<{ isCollapsed: boolean }>(({ theme, isCollapsed }) => ({
  width: isCollapsed ? 64 : 200,
  height: '100%',
  borderRight: '1px solid rgba(0,0,0,0)',
  padding: '24px 0',
  backgroundColor: 'rgba(20, 18, 30, 0.95)',
  display: 'flex',
  flexDirection: 'column',
  transition: theme.transitions.create('width', {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.enteringScreen,
  }),
  overflow: 'hidden',
  position: 'relative',
}));

const Logo = styled(Box)({
  padding: '0 24px',
  marginBottom: '24px',
  display: 'flex',
  alignItems: 'center',
  gap: '12px',
  '& svg': {
    width: '32px',
    height: '32px',
    display: 'block',
  },
  '& span': {
    fontSize: '24px',
    fontWeight: 700,
    color: 'primary.main',
    lineHeight: 1,
  },
});

const StyledListItem = styled(ListItem)(({ theme }) => ({
  padding: '12px 24px',
  '&:hover': {
    backgroundColor: 'rgba(255, 255, 255, 0.1)',
  },
  '&.Mui-selected': {
    backgroundColor: 'rgba(255, 255, 255, 0.1)',
    '&:hover': {
      backgroundColor: 'rgba(255, 255, 255, 0.15)',
    },
  },
}));

const navigationItems = [
  { text: 'Home', icon: <Home />, path: '/' },
  { text: 'Search', icon: <Search />, path: '/search' },
  { text: 'Library', icon: <LibraryMusic />, path: '/library' },
];

const ProfileSection = styled(Box)({
  display: 'flex',
  alignItems: 'center',
  gap: 12,
  padding: '12px 24px',
  minHeight: 64,
  cursor: 'pointer',
  background: 'none',
  boxSizing: 'border-box',
  marginBottom: '-10px',
});

interface SidebarProps {
  isCollapsed: boolean;
  onToggleCollapse: () => void;
}

const Sidebar = ({ isCollapsed, onToggleCollapse }: SidebarProps) => {
  const navigate = useNavigate();
  const location = useLocation();
  const { logout, isAdmin, user } = useAuthContext();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  const handleProfileClick = (event: React.MouseEvent<HTMLElement>) => {
    if (!user) {
      navigate('/login', { 
        state: { 
          from: location.pathname,
          message: 'Чтобы получить доступ к профилю, необходимо войти'
        }
      });
      return;
    }
    setAnchorEl(event.currentTarget);
  };

  const handleProfileClose = () => {
    setAnchorEl(null);
  };

  const handleLogout = async () => {
    try {
      await logout();
      navigate('/');
    } catch (error) {
      console.error('Logout failed:', error);
    }
  };

  const handleProfileNavigate = () => {
    navigate('/me');
    handleProfileClose();
  };

  return (
    <SidebarContainer isCollapsed={isCollapsed}>
      <Box sx={{ 
        display: 'flex', 
        alignItems: 'center', 
        justifyContent: 'space-between', 
        px: 2, 
        mb: 2,
        position: 'relative',
        zIndex: 1
      }}>
        <Logo
          sx={{ 
            cursor: 'pointer',
            opacity: 1,
            transition: 'opacity 0.2s',
            overflow: 'hidden',
            whiteSpace: 'nowrap',
            '& span': {
              opacity: isCollapsed ? 0 : 1,
              transition: 'opacity 0.2s',
              display: isCollapsed ? 'none' : 'inline',
            }
          }}
          onClick={() => navigate('/')}
        >
          <span>Orpheon</span>
        </Logo>
        <IconButton 
          onClick={onToggleCollapse} 
          size="small"
          sx={{
            position: 'absolute',
            right: 8,
            top: '50%',
            transform: 'translateY(-50%)',
            zIndex: 2,
            backgroundColor: 'background.paper',
            '&:hover': {
              backgroundColor: 'rgba(255, 255, 255, 0.1)',
            }
          }}
        >
          {isCollapsed ? <ChevronRight /> : <ChevronLeft />}
        </IconButton>
      </Box>

      <List>
        {navigationItems.map((item) => (
          <Tooltip 
            key={item.text} 
            title={isCollapsed ? item.text : ''} 
            placement="right"
          >
            <StyledListItem
              onClick={() => navigate(item.path)}
              sx={{
                cursor: 'pointer',
                backgroundColor: location.pathname === item.path ? 'rgba(255, 255, 255, 0.1)' : 'transparent',
                px: isCollapsed ? 2 : 3,
                justifyContent: isCollapsed ? 'center' : 'flex-start',
              }}
            >
              <ListItemIcon sx={{ 
                color: 'inherit',
                minWidth: isCollapsed ? 'auto' : 40,
              }}>
                {item.icon}
              </ListItemIcon>
              {!isCollapsed && <ListItemText primary={item.text} />}
            </StyledListItem>
          </Tooltip>
        ))}
        {isAdmin && (
          <Tooltip 
            title={isCollapsed ? 'Admin' : ''} 
            placement="right"
          >
            <StyledListItem
              onClick={() => navigate('/admin')}
              sx={{
                cursor: 'pointer',
                backgroundColor: location.pathname === '/admin' ? 'rgba(255, 255, 255, 0.1)' : 'transparent',
                px: isCollapsed ? 2 : 3,
                justifyContent: isCollapsed ? 'center' : 'flex-start',
              }}
            >
              <ListItemIcon sx={{ 
                color: 'inherit',
                minWidth: isCollapsed ? 'auto' : 40,
              }}>
                <AdminPanelSettings />
              </ListItemIcon>
              {!isCollapsed && <ListItemText primary="Admin" />}
            </StyledListItem>
          </Tooltip>
        )}
      </List>

      <Box sx={{ flexGrow: 1 }} />

      <ProfileSection onClick={handleProfileClick}>
        {user ? (
          <>
            <Avatar 
              sx={{ width: 32, height: 32 }}
            >
              {user.name ? user.name.charAt(0).toUpperCase() : '?'}
            </Avatar>
            {!isCollapsed && (
              <Box sx={{ overflow: 'hidden' }}>
                <Typography variant="subtitle2" noWrap>
                  {user.name}
                </Typography>
              </Box>
            )}
          </>
        ) : (
          <>
            <Avatar sx={{ width: 32, height: 32 }}>
              <Person />
            </Avatar>
            {!isCollapsed && (
              <Typography variant="subtitle2">
                Войти
              </Typography>
            )}
          </>
        )}
      </ProfileSection>

      {user && (
        <Menu
          anchorEl={anchorEl}
          open={Boolean(anchorEl)}
          onClose={handleProfileClose}
          onClick={handleProfileClose}
          PaperProps={{
            elevation: 0,
            sx: {
              overflow: 'visible',
              filter: 'drop-shadow(0px 2px 8px rgba(0,0,0,0.32))',
              mt: 1.5,
              '& .MuiAvatar-root': {
                width: 32,
                height: 32,
                ml: -0.5,
                mr: 1,
              },
            },
          }}
          transformOrigin={{ horizontal: 'right', vertical: 'top' }}
          anchorOrigin={{ horizontal: 'right', vertical: 'bottom' }}
        >
          <MenuItem onClick={handleProfileNavigate}>
            <ListItemIcon>
              <Person fontSize="small" />
            </ListItemIcon>
            Профиль
          </MenuItem>
          {isAdmin && (
            <MenuItem onClick={() => navigate('/admin')}>
              <ListItemIcon>
                <AdminPanelSettings fontSize="small" />
              </ListItemIcon>
              Админ панель
            </MenuItem>
          )}
          <Divider />
          <MenuItem onClick={handleLogout}>
            <ListItemIcon>
              <Logout fontSize="small" />
            </ListItemIcon>
            Выйти
          </MenuItem>
        </Menu>
      )}
    </SidebarContainer>
  );
};

export default Sidebar; 