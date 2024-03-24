import React, { useEffect } from 'react';
import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom'
import Home from './pages/home';
import Login from './pages/login';
import Auth from './pages/auth';
import './assets/index.css';
import "./assets/output.css";
import { AuthProvider, useAuth } from './context/user-auth';
import {  ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

const AppRoutes: React.FC = () => {
  const { isAuthenticated } = useAuth();
  
  return (
    <Routes>
      <Route path="/" element={isAuthenticated ? <Home /> : <Navigate to="/login" />} />
      <Route path="/login" element={<Login />} />
      <Route path="/auth/twitter/callback" element={<Auth />} />
    </Routes>
  );
};

const App: React.FC = () => {
  return (
    <Router>
      <AuthProvider>
        <ToastContainer />
        <AppRoutes />
      </AuthProvider>
    </Router>
  );
};

export default App;