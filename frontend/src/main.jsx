import React from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import App from './App.jsx';
import { AuthProvider } from './context/AuthContext.jsx';
import { FileProvider } from './context/FileContext.jsx';
import './styles/tokens.css';
import './styles/global.css';
import './styles/layout.css';
import './styles/dashboard.css';
import './styles/animations.css';

createRoot(document.getElementById('root')).render(<React.StrictMode><BrowserRouter><AuthProvider><FileProvider><App /></FileProvider></AuthProvider></BrowserRouter></React.StrictMode>);
