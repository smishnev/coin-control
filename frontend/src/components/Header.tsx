import React from 'react';
import LanguageSwitcher from './LanguageSwitcher';
import ThemeSwitcher from './ThemeSwitcher';
import Avatar from './Avatar';

interface User {
  id: string;
  nickname: string;
  user_id: string;
}

interface HeaderProps {
  user: User | null;
  onLogout: () => void;
}

const Header: React.FC<HeaderProps> = ({ user, onLogout }) => {
  return (
    <header className="flex justify-between items-center h-16 px-6 bg-menu shadow-sm rounded-b-lg border-b border-border">
      <div className="flex items-center gap-4">
        {user && (
          <div className="flex items-center gap-2">
            <span className="text-sm text-muted-foreground">
              Welcome, {user.nickname}
            </span>
            <button
              onClick={onLogout}
              className="text-sm text-destructive hover:text-destructive-foreground px-2 py-1 rounded hover:bg-destructive/10 hover:bg-destructive/20"
            >
              Logout
            </button>
          </div>
        )}
      </div>
      <div className="flex gap-3 ms-auto">
        <LanguageSwitcher />
        <ThemeSwitcher />
        <Avatar />
      </div>
    </header>
  );
};

export default Header;