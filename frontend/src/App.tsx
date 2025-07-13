import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import AuthScreen from './components/AuthScreen';
import { useAuth } from './contexts/AuthContext';
import MainLayout from './layouts/MainLayout';
import BybitForm from './pages/bybit';
import UserProfile from './pages/user-profile';

const App: React.FC = () => {
  const [activeMenu, setActiveMenu] = useState('item1');
  const { t } = useTranslation();
  const { isAuthenticated, isLoading, user, logout } = useAuth();

  const menuContent: Record<string, React.ReactNode> = {
    item1: <div className="text-xl font-semibold">{t('content_item1')}</div>,
    item2: <div className="text-xl font-semibold">{t('content_item2')}</div>,
    bybit: <BybitForm />,
    userProfile: <UserProfile />,
  };

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">{t('Loading...')}</p>
        </div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return <AuthScreen />;
  }

  return (
    <MainLayout
      activeMenu={activeMenu}
      setActiveMenu={setActiveMenu}
      user={user}
      onLogout={logout}
    >
      {menuContent[activeMenu]}
    </MainLayout>
  );
};

export default App;