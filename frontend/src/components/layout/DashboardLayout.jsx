import { useEffect } from 'react';
import { Outlet, Navigate, useLocation } from 'react-router-dom';
import useAuthStore from '../../store/authStore';
import Sidebar from './Sidebar';
import './DashboardLayout.css';

export default function DashboardLayout() {
  const { isPending, refreshUser } = useAuthStore();
  const location = useLocation();

  useEffect(() => {
    refreshUser();
  }, [refreshUser]);

  if (isPending() && location.pathname !== '/team') {
    return <Navigate to="/team" replace />;
  }

  return (
    <div className="dashboard-layout">
      <Sidebar />
      <main className="dashboard-main">
        <Outlet />
      </main>
    </div>
  );
}
