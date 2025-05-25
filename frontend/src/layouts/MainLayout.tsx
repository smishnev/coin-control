import React from 'react';
import Header from '../components/Header';
import AppMenu from '../components/AppMenu';

interface MainLayoutProps {
  children: React.ReactNode;
  activeMenu: string;
  setActiveMenu: (key: string) => void;
}

const MainLayout: React.FC<MainLayoutProps> = ({ children, activeMenu, setActiveMenu }) => (
  <div className="flex h-screen bg-background text-foreground">
    <AppMenu activeMenu={activeMenu} setActiveMenu={setActiveMenu} />
    <div className="flex-1 flex flex-col">
      <Header />
      <main className="flex-1 p-6">{children}</main>
    </div>
  </div>
);

export default MainLayout;