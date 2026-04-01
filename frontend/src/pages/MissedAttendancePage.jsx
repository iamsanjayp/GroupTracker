import { useState, useEffect } from 'react';
import api from '../lib/api';
import useAuthStore from '../store/authStore';
import Card from '../components/ui/Card';
import Button from '../components/ui/Button';
import Input, { Select } from '../components/ui/Input';

export default function MissedAttendancePage() {
  const { user, isAdmin } = useAuthStore();
  const [date, setDate] = useState(() => new Date().toISOString().split('T')[0]);
  const [hourSlot, setHourSlot] = useState(1);
  const [loading, setLoading] = useState(false);
  const [exports, setExports] = useState([]);

  useEffect(() => {
    if (isAdmin()) {
      fetchExports();
    }
  }, [isAdmin]);

  const fetchExports = async () => {
    setLoading(true);
    try {
      const res = await api.get('/attendance/missed/exports');
      setExports(res.data || []);
    } catch (err) {
      console.error(err);
    }
    setLoading(false);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      await api.post('/attendance/missed', {
        date,
        hour_slot: parseInt(hourSlot)
      });
      alert('Missed OTP logged successfully. It will be reported to the college for manual entry.');
      if (isAdmin()) {
        fetchExports();
      }
    } catch (err) {
      alert(err.response?.data?.error || 'Failed to log missed OTP');
    }
    setLoading(false);
  };

  const handleExportXLSX = async () => {
    if (exports.length === 0) return;

    try {
      const XLSX = await import('xlsx');
      
      const wsData = exports.map(exp => ({
        'Date': exp.Date,
        'Roll No': exp['Roll No'],
        'Name': exp.Name,
        'Mail Id': exp['Mail Id'],
        'Hour': exp.Hour
      }));

      const ws = XLSX.utils.json_to_sheet(wsData);
      
      // Auto-fit columns roughly
      const colWidths = [
        { wch: 12 }, // Date
        { wch: 15 }, // Roll No
        { wch: 25 }, // Name
        { wch: 30 }, // Mail Id
        { wch: 10 }  // Hour
      ];
      ws['!cols'] = colWidths;

      const wb = XLSX.utils.book_new();
      XLSX.utils.book_append_sheet(wb, ws, "Missed Attendance");
      
      const fileName = `Missed_OTP_Export_${new Date().toISOString().split('T')[0]}.xlsx`;
      XLSX.writeFile(wb, fileName);
    } catch (err) {
      console.error(err);
      alert('Failed to generate Excel file. Tell your admin to run `npm install xlsx` in the frontend directory.');
    }
  };

  return (
    <div className="animate-fade-in">
      <div className="page-header">
        <div>
          <h1>Missed OTP Portal</h1>
          <p>Report missed group hours for manual attendance marking</p>
        </div>
      </div>

      <div style={{ display: 'grid', gap: '2rem', gridTemplateColumns: 'minmax(300px, 1fr) 2fr', alignItems: 'start' }}>
        <Card>
          <h3>Log Missed Entry</h3>
          <p className="text-sm text-muted" style={{ marginBottom: '1rem', marginTop: '0.5rem' }}>
            If you failed to enter your OTP for a specific group hour, report it here.
          </p>
          <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
            <Input 
              label="Date" 
              type="date" 
              value={date} 
              onChange={e => setDate(e.target.value)} 
              required 
            />
            <Select 
              label="Hour Slot" 
              value={hourSlot} 
              onChange={e => setHourSlot(e.target.value)}
              required
            >
              {[1,2,3,4,5,6,7].map(h => (
                <option key={h} value={h}>Hour {h}</option>
              ))}
            </Select>
            <Button type="submit" fullWidth loading={loading} style={{ marginTop: '0.5rem' }}>
              Submit Report
            </Button>
          </form>
        </Card>

        {/* Admin Export View */}
        {isAdmin() && (
          <Card>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
              <div>
                <h3>Export Data for Institute</h3>
                <p className="text-sm text-muted" style={{ marginTop: '0.5rem' }}>
                  A consolidated list of all missed OTP reports from your team members.
                </p>
              </div>
              <Button variant="primary" onClick={handleExportXLSX} disabled={exports.length === 0}>
                Download XLSX
              </Button>
            </div>

            {loading ? (
              <div style={{ padding: '2rem', textAlign: 'center' }}><div className="spinner"></div></div>
            ) : exports.length === 0 ? (
              <div style={{ padding: '2rem', textAlign: 'center', color: 'var(--text-muted)' }}>
                No missed OTP reports found for this team.
              </div>
            ) : (
              <div className="table-responsive" style={{ overflowX: 'auto' }}>
                <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left', marginTop: '1rem' }}>
                  <thead>
                    <tr style={{ borderBottom: '1px solid var(--border)' }}>
                      <th style={{ padding: '0.75rem' }}>Date</th>
                      <th style={{ padding: '0.75rem' }}>Roll No</th>
                      <th style={{ padding: '0.75rem' }}>Name</th>
                      <th style={{ padding: '0.75rem' }}>Mail Id</th>
                      <th style={{ padding: '0.75rem' }}>Hour</th>
                    </tr>
                  </thead>
                  <tbody>
                    {exports.map((exp, idx) => (
                      <tr key={idx} style={{ borderBottom: '1px solid var(--border)' }}>
                        <td style={{ padding: '0.75rem' }}>{exp.Date}</td>
                        <td style={{ padding: '0.75rem' }}>{exp['Roll No']}</td>
                        <td style={{ padding: '0.75rem' }}>{exp.Name}</td>
                        <td style={{ padding: '0.75rem' }}>{exp['Mail Id']}</td>
                        <td style={{ padding: '0.75rem' }}>{exp.Hour}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </Card>
        )}
      </div>
    </div>
  );
}
