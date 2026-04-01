import { useState, useEffect } from 'react';
import api from '../lib/api';
import useAuthStore from '../store/authStore';
import Card from '../components/ui/Card';
import Button from '../components/ui/Button';
import Input from '../components/ui/Input';
import Modal from '../components/ui/Modal';
import './PointsPage.css';

export default function PointsPage() {
  const { isAdmin } = useAuthStore();
  const [points, setPoints] = useState(null);
  const [leaderboard, setLeaderboard] = useState([]);
  const [history, setHistory] = useState({ transactions: [], total_pages: 0 });
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);

  const [showModal, setShowModal] = useState(false);
  const [reason, setReason] = useState('');
  const [activityPoints, setActivityPoints] = useState('');
  const [rewardPoints, setRewardPoints] = useState('');
  const [saving, setSaving] = useState(false);

  const [showBulkModal, setShowBulkModal] = useState(false);
  const [bulkFile, setBulkFile] = useState(null);
  const [bulkUploading, setBulkUploading] = useState(false);

  useEffect(() => { fetchData(); }, [page]);

  const fetchData = async () => {
    setLoading(true);
    try {
      const [pointsRes, lbRes, histRes] = await Promise.all([
        api.get('/points/me'),
        api.get('/points/team').catch(() => null),
        api.get(`/points/history?page=${page}&limit=20`).catch(err => {
          console.error('History API err:', err);
          return { data: { transactions: [], totalPages: 0 } };
        }),
      ]);
      setPoints(pointsRes.data.points);
      
      if (lbRes) setLeaderboard(lbRes.data || []);
      if (histRes) setHistory(histRes.data || { transactions: [], total_pages: 0 });
    } catch (err) { console.error(err); }
    setLoading(false);
  };

  const handleManualAdd = async () => {
    if (!reason.trim()) { alert('Reason is required'); return; }
    
    setSaving(true);
    try {
      await api.post('/points/ps', {
        course_name: `[Manual] ${reason}`,
        level: 1, // dummy
        reward_points: parseFloat(rewardPoints) || 0,
        activity_points: parseFloat(activityPoints) || 0 // Reusing reward_points interface but adding ability to send activity points? Wait, the backend uses reward_points for both! Let's check the backend.
      });
      setShowModal(false);
      setReason('');
      setActivityPoints('');
      setRewardPoints('');
      fetchData();
    } catch (err) {
      alert(err.response?.data?.error || 'Failed to add manual points');
    }
    setSaving(false);
  };

  const handleBulkUpload = async () => {
    if (!bulkFile) {
        alert("Please select an Excel file (.xlsx)");
        return;
    }
    setBulkUploading(true);
    try {
        const XLSX = await import('xlsx');
        const buffer = await bulkFile.arrayBuffer();
        const workbook = XLSX.read(buffer, { type: 'array' });
        const firstSheetName = workbook.SheetNames[0];
        const rows = XLSX.utils.sheet_to_json(workbook.Sheets[firstSheetName]);
        
        const records = rows.map(r => ({
            email: String(r.Email || r.email || '').trim(),
            roll_no: String(r['Roll No'] || r['RollNo'] || r.roll_no || '').trim(),
            reason: String(r.Reason || r.reason || 'Bulk Upload').trim(),
            activity_points: parseFloat(r['Activity Points'] || r.ActivityPoints || r.activity_points) || 0,
            reward_points: parseFloat(r['Reward Points'] || r.RewardPoints || r.reward_points) || 0
        })).filter(r => r.email && r.roll_no && (r.activity_points > 0 || r.reward_points > 0));

        if (records.length === 0) {
            alert("No valid rows found. Ensure columns 'Email', 'Roll No', and at least one points column exist with > 0 points.");
            setBulkUploading(false);
            return;
        }

        const res = await api.post('/points/bulk', { records });
        alert(`Successfully imported ${res.data.success_count} records.\nFailed: ${res.data.failed_rows.length}\n${res.data.failed_rows.join(', ')}`);
        setShowBulkModal(false);
        setBulkFile(null);
        fetchData();
    } catch (err) {
        console.error(err);
        alert(err.response?.data?.error || "Failed to process bulk upload. Verify 'xlsx' package is installed.");
    }
    setBulkUploading(false);
  };

  if (loading) {
    return <div className="page-loading"><div className="spinner spinner-lg"></div></div>;
  }

  const totalActivity = points?.total_activity || 0;
  const totalReward = points?.total_reward || 0;
  const totalPoints = totalActivity + totalReward;

  return (
    <div className="animate-fade-in">
      <div className="page-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: '10px' }}>
        <div>
          <h1>Points & Progress</h1>
          <p>Track your activity points, rewards, and team standings</p>
        </div>
        <div style={{ display: 'flex', gap: '10px' }}>
          {isAdmin() && (
            <Button onClick={() => setShowBulkModal(true)} variant="secondary">Bulk Import</Button>
          )}
          <Button onClick={() => setShowModal(true)} variant="primary">Add Manual Points</Button>
        </div>
      </div>

      {/* Points Overview */}
      <div className="stats-grid">
        <div className="stat-card">
          <div className="stat-card-header">
            <span className="stat-card-label">Total Points</span>
            <div className="stat-card-icon purple">🏆</div>
          </div>
          <div className="stat-card-value">{totalPoints.toFixed(1)}</div>
          <div className="stat-card-sub">Activity + Reward combined</div>
        </div>
        <div className="stat-card">
          <div className="stat-card-header">
            <span className="stat-card-label">Activity Points</span>
            <div className="stat-card-icon green">📊</div>
          </div>
          <div className="stat-card-value">{totalActivity.toFixed(1)}</div>
          <div className="stat-card-sub">From daily activity logs</div>
        </div>
        <div className="stat-card">
          <div className="stat-card-header">
            <span className="stat-card-label">Reward Points</span>
            <div className="stat-card-icon orange">⭐</div>
          </div>
          <div className="stat-card-value">{totalReward.toFixed(1)}</div>
          <div className="stat-card-sub">From PS Slot completions</div>
        </div>
      </div>

      <div className="points-tip">
        <span className="points-tip-icon">💡</span>
        <span>Use the <strong>Add Manual Points</strong> button if you need to log activity or reward points for reasons outside the daily logger.</span>
      </div>

      {/* Transaction History */}
      <div className="section-header" style={{ marginTop: '30px' }}>
        <h2>📝 Transaction History</h2>
      </div>
      <Card>
        {(history.transactions || []).length === 0 ? (
          <p className="text-muted" style={{ textAlign: 'center', padding: '20px' }}>No point transactions yet.</p>
        ) : (
          <>
            <div className="table-responsive">
              <table className="leaderboard-table">
                <thead>
                  <tr>
                    <th>Date</th>
                    <th>Reason</th>
                    <th>Source</th>
                    <th>Activity Pts</th>
                    <th>Reward Pts</th>
                  </tr>
                </thead>
                <tbody>
                  {(history.transactions || []).map((tx) => (
                    <tr key={tx.id}>
                      <td>{tx.date}</td>
                      <td className="lb-name">{tx.reason?.replace('_', ' ')}</td>
                      <td>
                        <span className={`badge badge-${tx.source === 'manual' ? 'warning' : 'primary'}`}>
                          {tx.source}
                        </span>
                      </td>
                      <td style={{ color: tx.activity_points > 0 ? 'var(--success)' : 'inherit' }}>
                        {tx.activity_points > 0 ? `+${tx.activity_points}` : '-'}
                      </td>
                      <td style={{ color: tx.reward_points > 0 ? 'var(--warning)' : 'inherit' }}>
                        {tx.reward_points > 0 ? `+${tx.reward_points}` : '-'}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {/* Pagination Controls */}
            {history.total_pages > 1 && (
              <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', gap: '15px', marginTop: '20px' }}>
                <Button 
                  size="sm" 
                  variant="secondary" 
                  disabled={page === 1} 
                  onClick={() => setPage(page - 1)}
                >
                  Previous
                </Button>
                <span className="text-sm">Page {page} of {history.total_pages}</span>
                <Button 
                  size="sm" 
                  variant="secondary" 
                  disabled={page === history.total_pages} 
                  onClick={() => setPage(page + 1)}
                >
                  Next
                </Button>
              </div>
            )}
          </>
        )}
      </Card>

      {/* Leaderboard */}
      <div className="section-header" style={{ marginTop: '30px' }}>
        <h2>🏆 Team Leaderboard</h2>
      </div>
      <Card>
        {leaderboard.length === 0 ? (
          <p className="text-muted" style={{ textAlign: 'center', padding: '20px' }}>No data yet</p>
        ) : (
          <table className="leaderboard-table">
            <thead>
              <tr>
                <th>Rank</th>
                <th>Name</th>
                <th>Activity</th>
                <th>Reward</th>
                <th>Total</th>
              </tr>
            </thead>
            <tbody>
              {leaderboard.map((entry, i) => (
                <tr key={entry.user_id} className={i < 3 ? 'top-3' : ''}>
                  <td>
                    <span className="rank">
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
        )}
      </Card>

      <Modal isOpen={showModal} onClose={() => setShowModal(false)} title="Add Manual Points" size="sm">
        <div className="modal-form">
          <Input 
            label="Reason / Activity Name" 
            value={reason} 
            onChange={e => setReason(e.target.value)} 
            placeholder="E.g., Special event participation" 
            required 
          />
          <Input 
            label="Activity Points" 
            type="number" 
            step="0.1" 
            min="0"
            value={activityPoints} 
            onChange={e => setActivityPoints(e.target.value)} 
          />
          <Input 
            label="Reward Points" 
            type="number" 
            step="0.1"
            min="0"
            value={rewardPoints} 
            onChange={e => setRewardPoints(e.target.value)} 
          />
          <div style={{ display: 'flex', gap: '10px', marginTop: '10px' }}>
            <Button fullWidth variant="ghost" onClick={() => setShowModal(false)}>Cancel</Button>
            <Button fullWidth onClick={handleManualAdd} loading={saving}>Add Points</Button>
          </div>
        </div>
      </Modal>

      <Modal isOpen={showBulkModal} onClose={() => setShowBulkModal(false)} title="Bulk Upload Points" size="md">
        <div className="modal-form">
          <p className="text-sm text-muted" style={{ marginBottom: '1rem' }}>
            Upload an Excel (.xlsx) file mapping points. Required columns: <strong>Email</strong>, <strong>Roll No</strong>. Optimal points columns: <strong>Activity Points</strong>, <strong>Reward Points</strong>, <strong>Reason</strong>.
          </p>
          <div style={{ border: '2px dashed var(--border)', padding: '2rem', textAlign: 'center', borderRadius: '8px' }}>
            <input 
              type="file" 
              accept=".xlsx, .xls"
              onChange={(e) => setBulkFile(e.target.files[0])}
              style={{ display: 'block', margin: '0 auto' }}
            />
            {bulkFile && <p style={{ marginTop: '1rem', color: 'var(--primary)' }}>Selected: {bulkFile.name}</p>}
          </div>
          <div style={{ display: 'flex', gap: '10px', marginTop: '1rem' }}>
            <Button fullWidth variant="ghost" onClick={() => setShowBulkModal(false)}>Cancel</Button>
            <Button fullWidth onClick={handleBulkUpload} loading={bulkUploading} disabled={!bulkFile}>Import Data</Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}
