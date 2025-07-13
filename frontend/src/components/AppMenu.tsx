import React from 'react';
import { useTranslation } from 'react-i18next';

const menuItems = [
  { label: 'item1', key: 'item1' },
  { label: 'item2', key: 'item2' },
  { label: 'bybit', key: 'bybit' },
  { label: 'userProfile', key: 'userProfile' },
];

interface AppMenuProps {
  activeMenu: string;
  setActiveMenu: (key: string) => void;
}

const AppMenu: React.FC<AppMenuProps> = ({ activeMenu, setActiveMenu }) => {
  const { t } = useTranslation();

  return (
    <aside className="h-screen w-56 bg-menu border-r border-gray-200 dark:border-gray-700 flex flex-col py-6 px-2">
      <nav className="flex flex-col gap-1">
        {menuItems.map((item) => (
          <button
            key={item.key}
            onClick={() => setActiveMenu(item.key)}
            className={`flex items-center px-4 py-2 rounded-lg text-base font-medium transition
              ${
                activeMenu === item.key
                  ? 'bg-menu-active text-primary'
                  : 'text-foreground hover:bg-menu-active'
              }`}
          >
            {t(item.label)}
          </button>
        ))}
      </nav>
    </aside>
  );
};

export default AppMenu;