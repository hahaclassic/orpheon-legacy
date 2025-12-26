import { Box, styled } from '@mui/material';
import Sidebar from './Sidebar';
import PlayerBar from './PlayerBar';
import { useState } from 'react';
import { useAuthContext } from '../../contexts/AuthContext';

const LayoutContainer = styled(Box)({
  display: 'flex',
  height: '100vh',
  overflow: 'hidden',
  backgroundColor: 'rgba(20, 18, 30, 0.95)',
});

const MainContent = styled(Box, {
  shouldForwardProp: (prop) => prop !== 'isSidebarCollapsed'
})<{ isSidebarCollapsed: boolean }>(({ theme, isSidebarCollapsed }) => ({
  flex: 1,
  display: 'flex',
  flexDirection: 'column',
  overflow: 'hidden',
  minWidth: 0,
  transition: theme.transitions.create('margin', {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.enteringScreen,
  }),
}));

const ContentArea = styled(Box)({
  flex: 1,
  overflow: 'auto',
  padding: '24px',
  backgroundColor: '#181825',
  display: 'flex',
  flexDirection: 'column',
  minHeight: 0,
  borderRadius: '12px',
  margin: '12px 12px 12px 0',
  boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)',
  border: '1px solid rgba(255, 255, 255, 0.1)',
});

const PlayerBarContainer = styled(Box)({
  height: 140,
  backgroundColor: '#181825',
  borderTop: 'none',
  borderColor: 'divider',
  margin: '4px 12px 16px 0',
  borderRadius: '12px',
  boxShadow: '0 -4px 6px -1px rgba(0, 0, 0, 0.1)',
});

const Layout = ({ children }: { children: React.ReactNode }) => {
  const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(false);
  const { isAuthenticated } = useAuthContext();

  return (
    <LayoutContainer>
      <Sidebar 
        isCollapsed={isSidebarCollapsed} 
        onToggleCollapse={() => setIsSidebarCollapsed(!isSidebarCollapsed)} 
      />
      <MainContent isSidebarCollapsed={isSidebarCollapsed}>
        <ContentArea>
          {children}
        </ContentArea>
        {isAuthenticated && (
          <PlayerBarContainer>
            <PlayerBar />
          </PlayerBarContainer>
        )}
      </MainContent>
    </LayoutContainer>
  );
};

export default Layout; 