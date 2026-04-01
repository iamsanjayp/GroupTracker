import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import useAuthStore from '../store/authStore';
import api from '../lib/api';
import Card from '../components/ui/Card';
import Button from '../components/ui/Button';
import ProgressBar from '../components/ui/ProgressBar';
import Badge from '../components/ui/Badge';
import './DashboardPage.css';

export default function DashboardPage() {
  const { user, isAdmin, hasTeam } = useAuthStore();
  const navigate = useNavigate();
  const [memberData, setMemberData] = useState(null);
  const [adminData, setAdminData] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!hasTeam()) return;

    const fetchData = async () => {
      setLoading(true);
      try {
        const memberRes = await api.get('/dashboard/member');
        setMemberData(memberRes.data);

        if (isAdmin()) {
          const adminRes = await api.get('/dashboard/admin');
          setAdminData(adminRes.data);
        }
      } catch (err) {
        console.error('Dashboard fetch error:', err);
      }
      setLoading(false);
    };

    fetchData();
  }, []);

  if (!hasTeam()) {
    return (
      <div className="no-team-container">
        <div className="no-team-card">
          <span className="no-team-icon">👥</span>
          <h2>Join or Create a Team</h2>
          <p>You need to be part of a team to use GroupTracker</p>
          <div className="no-team-actions">
            <Button onClick={() => navigate('/team')}>Go to Team Setup</Button>
          </div>
        </div>
      </div>
    );
  }

  if (loading) {
    return (
      <div className="page-loading">
        <div className="spinner spinner-lg"></div>
        <p>Loading dashboard...</p>
      </div>
    );
  }

  const roleBadgeVariant = (role) => {
    const map = { captain: 'primary', vice_captain: 'info', manager: 'warning', strategist: 'success', member: 'default' };
    return map[role] || 'default';
  };

  return (
    <div className="animate-fade-in">
      <div className="page-header">
        <h1>Dashboard</h1>
        <p>Welcome back, {user?.name}! Here's your overview for today.</p>
      </div>

      {/* Stats Cards */}
      <div className="stats-grid">
        <div className="stat-card">
          <div className="stat-card-header">
            <span className="stat-card-label">Hours Logged Today</span>
            <div className="stat-card-icon purple">🕐</div>
          </div>
          <div className="stat-card-value">{memberData?.today?.hours_logged || 0}/7</div>
          <div className="stat-card-sub">
            <ProgressBar value={memberData?.today?.hours_logged || 0} max={7} showValue={false} />
          </div>
        </div>

        <div className="stat-card">
          <div className="stat-card-header">
            <span className="stat-card-label">Today's Activity Points</span>
            <div className="stat-card-icon green">📊</div>
          </div>
          <div className="stat-card-value">{memberData?.today?.activity_points?.toFixed(1) || '0.0'}</div>
          <div className="stat-card-sub">+{memberData?.today?.reward_points?.toFixed(1) || '0.0'} reward pts</div>
        </div>

        <div className="stat-card">
          <div className="stat-card-header">
            <span className="stat-card-label">Total Points</span>
            <div className="stat-card-icon orange">⭐</div>
          </div>
          <div className="stat-card-value">{memberData?.total?.total_points?.toFixed(1) || '0.0'}</div>
          <div className="stat-card-sub">{memberData?.total?.activity_points?.toFixed(1)} activity + {memberData?.total?.reward_points?.toFixed(1)} reward</div>
        </div>

        <div className="stat-card">
          <div className="stat-card-header">
            <span className="stat-card-label">Days Logged This Month</span>
            <div className="stat-card-icon blue">📅</div>
          </div>
          <div className="stat-card-value">{memberData?.month?.logged_days || 0}</div>
          <div className="stat-card-sub">Keep the streak going!</div>
        </div>

        <div className="stat-card">
          <div className="stat-card-header">
            <span className="stat-card-label">Attendance</span>
            <div className="stat-card-icon blue">📋</div>
          </div>
          <div className="stat-card-value">
            {memberData?.attendance_percentage !== undefined ? `${memberData.attendance_percentage.toFixed(1)}%` : '100%'}
          </div>
          <div className="stat-card-sub">Group hour presence</div>
        </div>
      </div>

      {/* Quick Actions */}
      <div className="section-header">
        <h2>Quick Actions</h2>
      </div>
      <div className="quick-actions">
        <Card hover className="quick-action-card" onClick={() => navigate('/daily-log')}>
          <span className="qa-icon">🕐</span>
          <span className="qa-label">Log Today's Activity</span>
          <span className="qa-arrow">→</span>
        </Card>
        <Card hover className="quick-action-card" onClick={() => navigate('/projects')}>
          <span className="qa-icon">📁</span>
          <span className="qa-label">View Projects</span>
          <span className="qa-arrow">→</span>
        </Card>
        <Card hover className="quick-action-card" onClick={() => navigate('/points')}>
          <span className="qa-icon">⭐</span>
          <span className="qa-label">Points & Leaderboard</span>
          <span className="qa-arrow">→</span>
        </Card>
      </div>

      {/* Admin Section */}
      {isAdmin() && adminData && (
        <div className="admin-section">
          <div className="section-header">
            <h2>Team Overview</h2>
            <Badge variant="primary">Admin</Badge>
          </div>

          <div className="stats-grid">
            <div className="stat-card">
              <div className="stat-card-header">
                <span className="stat-card-label">Active Today</span>
                <div className="stat-card-icon green">👥</div>
              </div>
              <div className="stat-card-value">
                {adminData.today?.active_today}/{adminData.today?.total_members}
              </div>
              <div className="stat-card-sub">team members logged activity</div>
            </div>

            <div className="stat-card">
              <div className="stat-card-header">
                <span className="stat-card-label">Projects</span>
                <div className="stat-card-icon purple">📁</div>
              </div>
              <div className="stat-card-value">{adminData.projects?.total || 0}</div>
              <div className="stat-card-sub">
                {adminData.projects?.completed_tasks}/{adminData.projects?.total_tasks} tasks done
              </div>
            </div>
          </div>

          {/* Leaderboard */}
          {adminData.leaderboard?.length > 0 && (
            <Card>
              <h3 style={{ marginBottom: '16px' }}>🏆 Team Leaderboard</h3>
              <table className="leaderboard-table">
                <thead>
                  <tr>
                    <th>#</th>
                    <th>Name</th>
                    <th>Activity</th>
                    <th>Reward</th>
                    <th>Total</th>
                  </tr>
                </thead>
                <tbody>
                  {adminData.leaderboard.map((entry, i) => (
                    <tr key={entry.user_id} className={i < 3 ? 'top-3' : ''}>
                      <td>
                        <span className={`rank rank-${i + 1}`}>
                          {i === 0 ? '🥇' : i === 1 ? '🥈' : i === 2 ? '🥉' : i + 1}
                        </span>
                      </td>
                      <td className="lb-name">{entry.name}</td>
                      <td>{entry.total_activity?.toFixed(1)}</td>
                      <td>{entry.total_reward?.toFixed(1)}</td>
                      <td className="lb-total">{entry.total_points?.toFixed(1)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </Card>
          )}
        </div>
      )}
    </div>
  );
}
