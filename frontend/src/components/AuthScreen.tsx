import React, { useState } from 'react';
import ForgotPassword from './ForgotPassword';
import Login from './Login';
import Register from './Register';

type AuthMode = 'login' | 'register' | 'forgot-password';

const AuthScreen: React.FC = () => {
  const [mode, setMode] = useState<AuthMode>('login');

  const switchToRegister = () => setMode('register');
  const switchToLogin = () => setMode('login');
  const switchToForgotPassword = () => setMode('forgot-password');

  return (
    <div>
      {mode === 'login' && (
        <Login 
          onSwitchToRegister={switchToRegister} 
          onSwitchToForgotPassword={switchToForgotPassword}
        />
      )}
      {mode === 'register' && (
        <Register onSwitchToLogin={switchToLogin} />
      )}
      {mode === 'forgot-password' && (
        <ForgotPassword onSwitchToLogin={switchToLogin} />
      )}
    </div>
  );
};

export default AuthScreen; 