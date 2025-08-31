import React from 'react';
import { useTranslation } from 'react-i18next';
import AuthScreen from './components/AuthScreen';
import { useAuth } from './contexts/AuthContext';
import AppRouter from './components/AppRouter';

const App: React.FC = () => {
  const { t } = useTranslation();
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center ">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 mx-auto"></div>
          <p className="mt-4 ">{t('Loading...')}</p>
        </div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return <AuthScreen />;
  }

  return <AppRouter />;
};

export default App;