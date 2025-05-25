import React, { useState } from 'react';
import MainLayout from './layouts/MainLayout';
import { useTranslation } from 'react-i18next';

const App: React.FC = () => {
  const [activeMenu, setActiveMenu] = useState('item1');
  const { t } = useTranslation();

  const menuContent: Record<string, React.ReactNode> = {
    item1: <div className="text-xl font-semibold">{t('content_item1')}</div>,
    item2: <div className="text-xl font-semibold">{t('content_item2')}</div>,
  };

  return (
    <MainLayout
      activeMenu={activeMenu}
      setActiveMenu={setActiveMenu}
    >
      {menuContent[activeMenu]}
    </MainLayout>
  );
};

export default App;