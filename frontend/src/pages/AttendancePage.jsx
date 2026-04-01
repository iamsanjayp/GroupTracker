import { useState, useEffect } from 'react';
import { Navigate } from 'react-router-dom';
import api from '../lib/api';
import useAuthStore from '../store/authStore';
import Card from '../components/ui/Card';
import Button from '../components/ui/Button';
import Input, { Select } from '../components/ui/Input';
import './AttendancePage.css';

const STATUS_OPTIONS = ['Present', 'Absent', 'PS Slot', 'Event', 'OnDuty', 'Class'];

export default function AttendancePage() {
  const { isAdmin } = useAuthStore();
  const [date, setDate] = useState(() => new Date().toISOString().split('T')[0]);
  const [session, setSession] = useState('morning');
  const [members, setMembers] = useState([]);
  const [attendanceData, setAttendanceData] = useState({});
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    if (isAdmin()) {
      fetchData();
    }
  }, [date, session]);

  const fetchData = async () => {
    setLoading(true);
    try {
      const [membersRes, attRes] = await Promise.all([
        api.get('/teams/members'),
        api.get(`/attendance?date=${date}&session=${session}`)
      ]);
      
      const activeMembers = (membersRes.data || []).filter(m => m.join_status !== 'pending');
      setMembers(activeMembers);

      // Initialize map
      const initialMap = {};
      activeMembers.forEach(m => {
        initialMap[m.id] = {};
        const hours = session === 'morning' ? [1,2,3,4] : [5,6,7];
        hours.forEach(h => {
          initialMap[m.id][h] = 'Present'; // default
        });
      });

      // Override with saved data
      (attRes.data || []).forEach(record => {
        if (initialMap[record.user_id]) {
          initialMap[record.user_id][record.hour_slot] = record.status;
        }
      });

      setAttendanceData(initialMap);
    } catch (err) {
      console.error(err);
    }
    setLoading(false);
  };

  const handleStatusChange = (userId, hour, status) => {
    setAttendanceData(prev => ({
      ...prev,
      [userId]: {
        ...prev[userId],
        [hour]: status
      }
    }));
  };

  const handleSave = async () => {
    setSaving(true);
    try {
      const records = [];
      Object.keys(attendanceData).forEach(userId => {
        Object.keys(attendanceData[userId]).forEach(hour => {
          records.push({
            user_id: parseInt(userId),
            hour_slot: parseInt(hour),
            status: attendanceData[userId][hour]
          });
        });
      });

      await api.post('/attendance', {
        date,
        session,
        records
      });

      alert('Attendance saved successfully!');
    } catch (err) {
      alert(err.response?.data?.error || 'Failed to save attendance');
    }
    setSaving(false);
  };

  if (!isAdmin()) {
    return <Navigate to="/dashboard" replace />;
  }

  const hours = session === 'morning' ? [1, 2, 3, 4] : [5, 6, 7];

  return (
    <div className="animate-fade-in attendance-page">
      <div className="page-header">
        <div>
          <h1>Attendance</h1>
          <p>Mark group hour attendance for your team</p>
        </div>
      </div>

      <Card className="attendance-controls" style={{ marginBottom: '20px', padding: '20px' }}>
        <div style={{ display: 'flex', gap: '20px', alignItems: 'flex-end', flexWrap: 'wrap' }}>
          <div style={{ flex: 1, minWidth: '200px' }}>
            <Input 
              type="date" 
              label="Date" 
              value={date} 
              onChange={e => setDate(e.target.value)} 
            />
          </div>
          <div style={{ flex: 1, minWidth: '200px' }}>
            <Select 
              label="Session" 
              value={session} 
              onChange={e => setSession(e.target.value)}
            >
              <option value="morning">Morning (Hours 1 - 4)</option>
              <option value="afternoon">Afternoon (Hours 5 - 7)</option>
            </Select>
          </div>
          <div style={{ flex: 'none' }}>
            <Button onClick={handleSave} loading={saving} variant="primary">
              Save Attendance
            </Button>
          </div>
        </div>
      </Card>

      {loading ? (
        <div className="page-loading"><div className="spinner spinner-md"></div></div>
      ) : (
        <div className="attendance-table-container">
          <table className="attendance-table">
            <thead>
              <tr>
                <th className="sticky-col">Member</th>
                {hours.map(h => (
                  <th key={h} className="hour-col">Hour {h}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {members.map(member => (
                <tr key={member.id}>
                  <td className="sticky-col">
                    <div className="member-info">
                      <strong>{member.name}</strong>
                      <span className="text-xs text-muted">{member.roll_no || 'No Roll No'}</span>
                    </div>
                  </td>
                  {hours.map(h => (
                    <td key={h}>
                      <select 
                        className={`status-select status-${attendanceData[member.id]?.[h]?.toLowerCase().replace(' ', '-')}`}
                        value={attendanceData[member.id]?.[h] || 'Present'}
                        onChange={(e) => handleStatusChange(member.id, h, e.target.value)}
                      >
                        {STATUS_OPTIONS.map(opt => (
                          <option key={opt} value={opt}>{opt}</option>
                        ))}
                      </select>
                    </td>
                  ))}
                </tr>
              ))}
              {members.length === 0 && (
                <tr>
                  <td colSpan={hours.length + 1} style={{ textAlign: 'center', padding: '2rem' }}>
                    No active members found in team.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
