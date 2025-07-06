import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import './styles/tailwind.css';
import './styles/theme.css';
import './i18n';
import { AuthProvider } from './contexts/AuthContext';

const container = document.getElementById('root');

const root = ReactDOM.createRoot(container!);

root.render(
    <React.StrictMode>
        <AuthProvider>
            <App />
        </AuthProvider>
    </React.StrictMode>
);