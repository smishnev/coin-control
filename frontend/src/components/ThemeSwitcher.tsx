import React, { useEffect, useState } from 'react';

const ThemeSwitcher: React.FC = () => {
  const [dark, setDark] = useState(() =>
    document.documentElement.classList.contains('dark')
  );

  useEffect(() => {
    if (dark) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }, [dark]);

  return (
    <button
      onClick={() => setDark((v) => !v)}
      className="w-10 h-10 flex items-center justify-center rounded-full bg-background transition"
      title="Theme Switcher"
    >
      {dark ? (
        <svg width="22" height="22" fill="none" viewBox="0 0 24 24">
          <path
            d="M21 12.79A9 9 0 0111.21 3a7 7 0 108.79 9.79z"
            fill="#5d87ff"
          />
        </svg>
      ) : (
        <svg width="22" height="22" fill="none" viewBox="0 0 24 24">
          <circle cx="12" cy="12" r="5" fill="#fbbf24" />
          <g stroke="#fbbf24" strokeWidth="2">
            <line x1="12" y1="1" x2="12" y2="3" />
            <line x1="12" y1="21" x2="12" y2="23" />
            <line x1="4.22" y1="4.22" x2="5.64" y2="5.64" />
            <line x1="18.36" y1="18.36" x2="19.78" y2="19.78" />
            <line x1="1" y1="12" x2="3" y2="12" />
            <line x1="21" y1="12" x2="23" y2="12" />
            <line x1="4.22" y1="19.78" x2="5.64" y2="18.36" />
            <line x1="18.36" y1="5.64" x2="19.78" y2="4.22" />
          </g>
        </svg>
      )}
    </button>
  );
};

export default ThemeSwitcher;