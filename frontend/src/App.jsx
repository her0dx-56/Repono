import { Navigate, Route, Routes } from 'react-router-dom';
import { useAuth } from './hooks/useAuth.js';
import LandingPage from './pages/LandingPage.jsx';
import AuthPage from './pages/AuthPage.jsx';
import DashboardPage from './pages/DashboardPage.jsx';
import SharedFilePage from './pages/SharedFilePage.jsx';
export default function App(){const {user}=useAuth();return <Routes><Route path="/" element={user?<Navigate to="/app" replace/>:<LandingPage/>}/><Route path="/home" element={<LandingPage/>}/><Route path="/login" element={user?<Navigate to="/app" replace/>:<AuthPage mode="login"/>}/><Route path="/signup" element={user?<Navigate to="/app" replace/>:<AuthPage mode="signup"/>}/><Route path="/app/*" element={user?<DashboardPage/>:<Navigate to="/login" replace/>}/><Route path="/share/:id" element={<SharedFilePage/>}/><Route path="*" element={<Navigate to="/" replace/>}/></Routes>}
