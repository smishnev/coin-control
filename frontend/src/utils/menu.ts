export interface MenuItem {
  key: string;
  label: string;
  path: string;
}

export const menuItems: MenuItem[] = [
  {
    key: 'item1',
    label: 'content_item1',
    path: '/'
  },
  {
    key: 'item2', 
    label: 'content_item2',
    path: '/item2'
  },
  {
    key: 'bybit',
    label: 'bybit',
    path: '/bybit'
  },
  {
    key: 'userProfile',
    label: 'userProfile',
    path: '/profile'
  }
];

export const getMenuItemByPath = (path: string): MenuItem | undefined => {
  return menuItems.find(item => item.path === path || path.startsWith(item.path + '/'));
};

export const getActiveMenuKey = (pathname: string): string => {
  // For nested routes like /bybit/BTC, we want to keep 'bybit' active
  if (pathname.startsWith('/bybit')) {
    return 'bybit';
  }
  
  const menuItem = menuItems.find(item => item.path === pathname);
  return menuItem?.key || 'item1';
};
