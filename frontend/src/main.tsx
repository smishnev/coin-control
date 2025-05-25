import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import './styles/tailwind.css';
import './styles/theme.css';
import './i18n';

const container = document.getElementById('root');

const root = ReactDOM.createRoot(container!);

root.render(
    <React.StrictMode>
        <App />
    </React.StrictMode>
);