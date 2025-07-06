import React, { createContext, ReactNode, useContext, useEffect, useState } from 'react';
import * as AppAPI from '../../wailsjs/go/main/App';

interface User {
  id: string;
  nickname: string;
  user_id: string;
}

interface AuthContextType {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (nickname: string, password: string) => Promise<void>;
  logout: () => void;
  register: (nickname: string, password: string, firstName: string, lastName: string) => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(localStorage.getItem('token'));
  const [isLoading, setIsLoading] = useState(true);

  const isAuthenticated = !!token && !!user;

  // Check token
  useEffect(() => {
    const validateToken = async () => {
      const storedToken = localStorage.getItem('token');
      if (storedToken) {
        try {
          const claims = await AppAPI.ValidateToken(storedToken);
          // If token is valid, get user information
          const authData = await AppAPI.GetAuthByID(claims.user_id);
          setUser({
            id: authData.id.toString(),
            nickname: authData.nickname,
            user_id: authData.user_id.toString(),
          });
          setToken(storedToken);
        } catch (error) {
          // Token is invalid, delete it
          localStorage.removeItem('token');
          setToken(null);
          setUser(null);
        }
      }
      setIsLoading(false);
    };

    validateToken();
  }, []);

  const login = async (nickname: string, password: string) => {
    try {
      const response = await AppAPI.Login({ nickname, password });
      if (response.auth) {
        setUser({
          id: response.auth.id.toString(),
          nickname: response.auth.nickname,
          user_id: response.auth.user_id.toString(),
        });
        setToken(response.token);
        localStorage.setItem('token', response.token);
      }
    } catch (error) {
      throw error;
    }
  };

  const logout = () => {
    setUser(null);
    setToken(null);
    localStorage.removeItem('token');
  };

  const register = async (nickname: string, password: string, firstName: string, lastName: string) => {
    try {
      // Create user and auth record in one transaction
      const authData = await AppAPI.CreateUserWithAuth({
        nickname,
        password,
        user_id: '', // Will be filled in transaction
      }, firstName, lastName);

      // Automatically login after registration
      await login(nickname, password);
    } catch (error) {
      throw error;
    }
  };

  const value: AuthContextType = {
    user,
    token,
    isAuthenticated,
    isLoading,
    login,
    logout,
    register,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
}; 