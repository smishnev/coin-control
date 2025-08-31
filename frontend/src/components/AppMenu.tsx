import React from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { menuItems } from '../utils/menu';

interface AppMenuProps {
  activeMenu: string;
  setActiveMenu?: (key: string) => void;
}

const AppMenu: React.FC<AppMenuProps> = ({ activeMenu, setActiveMenu }) => {
  const { t } = useTranslation();
  const navigate = useNavigate();

  const handleMenuClick = (item: any) => {
    setActiveMenu?.(item.key);
    navigate(item.path);
  };

  return (
    <aside className="h-screen w-56 bg-menu border-r border-border flex flex-col py-6 px-2">
      <nav className="flex flex-col gap-1">
        {menuItems.map((item) => (
          <button
            key={item.key}
            onClick={() => handleMenuClick(item)}
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