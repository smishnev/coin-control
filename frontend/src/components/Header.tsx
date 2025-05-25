import React from 'react';
import LanguageSwitcher from './LanguageSwitcher';
import ThemeSwitcher from './ThemeSwitcher';
import Avatar from './Avatar';

const Header: React.FC = () => {
  return (
    <header className="flex justify-between items-center h-16 px-6 bg-menu shadow-sm rounded-b-lg border-b border-gray-200 dark:border-gray-700">
      <div className="flex gap-3 ms-auto">
        <LanguageSwitcher />
        <ThemeSwitcher />
      <Avatar />
      </div>
    </header>
  );
};

export default Header;