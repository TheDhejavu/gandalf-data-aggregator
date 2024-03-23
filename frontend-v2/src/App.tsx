import React from 'react';
import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom'
import Home from './pages/home';
import Login from './pages/login';
import Auth from './pages/auth';
import './assets/index.css';
import "./assets/output.css";


const AuthenticatedRoutes: React.FC = () => {
  return (
    <Routes>
      <Route path="/" element={localStorage.getItem('token') != "" ? <Home /> : <Navigate to="/login" />} />
      <Route path="/login" element={<Login />} />
      <Route path="/auth/twitter/callback" element={<Auth />} />
    </Routes>
  );
};

const App: React.FC = () => {
  return (
    <Router>
      <AuthenticatedRoutes />
    </Router>
  );
};

export default App;