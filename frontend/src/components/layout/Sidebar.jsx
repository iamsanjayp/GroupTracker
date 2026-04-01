import { NavLink, useNavigate } from 'react-router-dom';
import useAuthStore from '../../store/authStore';
import './Sidebar.css';

const menuItems = [
  { path: '/dashboard', label: 'Dashboard', icon: '📊' },
  { path: '/daily-log', label: 'Daily Log', icon: '🕐' },
  { path: '/projects', label: 'Projects', icon: '📁' },
  { path: '/team', label: 'Team', icon: '👥' },
  { path: '/skills', label: 'Skills', icon: '🎯' },
  { path: '/points', label: 'Points', icon: '⭐' },
];

export default function Sidebar() {
  const { user, isAdmin, isPending, logout } = useAuthStore();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const getMenuItems = () => {
    if (isPending()) {
      return [{ path: '/team', label: 'Team', icon: '👥' }];
    }

    const items = [...menuItems];
    
    // Add Missed OTP for all active members
    items.push({ path: '/missed-attendance', label: 'Missed OTP', icon: '📝' });

    // Add Attendance only for leaders
    if (isAdmin()) {
      items.push({ path: '/attendance', label: 'Attendance', icon: '📋' });
    }

    return items;
  };

  const visibleMenu = getMenuItems();

  return (
    <aside className="sidebar">
      {/* Brand */}
      <div className="sidebar-brand">
        <div className="sidebar-logo">
          <span className="logo-icon">⚡</span>
        </div>
        <div className="sidebar-brand-text">
          <h4>GroupTracker</h4>
          <span className="text-xs text-muted">Team Productivity</span>
        </div>
      </div>

      {/* Navigation */}
      <nav className="sidebar-nav">
        <div className="sidebar-section-label">Menu</div>
        {visibleMenu.map(item => (
          <NavLink
            key={item.path}
            to={item.path}
            className={({ isActive }) => `sidebar-link ${isActive ? 'active' : ''}`}
          >
            <span className="sidebar-link-icon">{item.icon}</span>
            <span className="sidebar-link-label">{item.label}</span>
          </NavLink>
        ))}
      </nav>

      {/* User section */}
      <div className="sidebar-footer">
        <div className="sidebar-user">
          <div className="sidebar-avatar">
            {user?.name?.charAt(0)?.toUpperCase() || '?'}
          </div>
          <div className="sidebar-user-info">
            <span className="sidebar-user-name">{user?.name || 'User'}</span>
            <span className="sidebar-user-role">{user?.role?.replace('_', ' ') || 'member'}</span>
          </div>
        </div>
        <button className="sidebar-logout" onClick={handleLogout} title="Logout">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
            <polyline points="16 17 21 12 16 7" />
            <line x1="21" y1="12" x2="9" y2="12" />
          </svg>
        </button>
      </div>
    </aside>
  );
}
