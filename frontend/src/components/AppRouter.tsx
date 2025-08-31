import React, { Suspense, useState } from 'react';
import { HashRouter, Routes, Route, useLocation } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useAuth } from '../contexts/AuthContext';
import MainLayout from '../layouts/MainLayout';
import { getActiveMenuKey } from '../utils/menu';

// Lazy load components
const BybitPage = React.lazy(() => import('../pages/bybit'));
const CoinDetailPage = React.lazy(() => import('../pages/coin-detail'));
const UserProfilePage = React.lazy(() => import('../pages/user-profile'));

// Loading component
const PageLoading: React.FC<{ message?: string }> = ({ message = 'Loading...' }) => (
  <div className="p-4 flex items-center justify-center">
    <div className="text-center">
      <div className="animate-spin rounded-full h-8 w-8 border-b-2 mx-auto mb-2"></div>
      <p>{message}</p>
    </div>
  </div>
);

// Router content component
const RouterContent: React.FC = () => {
  const location = useLocation();
  const { user, logout } = useAuth();
  const { t } = useTranslation();
  const [activeMenu, setActiveMenu] = useState(getActiveMenuKey(location.pathname));

  // Update active menu when location changes
  React.useEffect(() => {
    setActiveMenu(getActiveMenuKey(location.pathname));
  }, [location.pathname]);

  return (
    <MainLayout
      activeMenu={activeMenu}
      user={user}
      onLogout={logout}
    >
      <Suspense fallback={<PageLoading />}>
        <Routes>
          <Route 
            path="/" 
            element={
              <div className="text-xl font-semibold">{t('content_item1')}</div>
            } 
          />
          <Route 
            path="/item2" 
            element={
              <div className="text-xl font-semibold">{t('content_item2')}</div>
            } 
          />
          <Route 
            path="/bybit" 
            element={
              <Suspense fallback={<PageLoading message="Loading Bybit…" />}>
                <BybitPage />
              </Suspense>
            } 
          />
          <Route 
            path="/bybit/:coinId" 
            element={
              <Suspense fallback={<PageLoading message="Loading Coin Details…" />}>
                <CoinDetailPage />
              </Suspense>
            } 
          />
          <Route 
            path="/profile" 
            element={
              <Suspense fallback={<PageLoading message="Loading Profile…" />}>
                <UserProfilePage />
              </Suspense>
            } 
          />
          {/* Fallback route */}
          <Route 
            path="*" 
            element={
              <div className="text-center py-8">
                <h2 className="text-2xl font-bold mb-4">Page Not Found</h2>
                <p>The page you're looking for doesn't exist.</p>
              </div>
            } 
          />
        </Routes>
      </Suspense>
    </MainLayout>
  );
};

const AppRouter: React.FC = () => {
  return (
    <HashRouter>
      <RouterContent />
    </HashRouter>
  );
};

export default AppRouter;
