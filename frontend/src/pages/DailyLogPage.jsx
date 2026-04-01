import { useState, useEffect } from 'react';
import api from '../lib/api';
import Button from '../components/ui/Button';
import { Select } from '../components/ui/Input';
import Badge from '../components/ui/Badge';
import './DailyLogPage.css';

const ACTIVITY_TYPES = [
  { value: '', label: 'Select type...' },
  { value: 'project_work', label: '💻 Project Work' },
  { value: 'ps_slot', label: '🎯 PS Slot' },
  { value: 'self_study', label: '📚 Self Study' },
  { value: 'event', label: '🎉 Event' },
  { value: 'class_participation', label: '🏫 Class Participation' },
];

const DEFAULT_POINTS = {
  project_work: 1.0,
  ps_slot: 1.0,
  self_study: 0.75,
  event: 1.5,
  class_participation: 0.5,
};

const HOUR_LABELS = [
  'Hour 1 — Morning Start',
  'Hour 2 — Focus Block',
  'Hour 3 — Deep Work',
  'Hour 4 — Midday',
  'Hour 5 — Afternoon',
  'Hour 6 — Late Focus',
  'Hour 7 — Wrap Up',
];

function getTodayStr() {
  return new Date().toISOString().split('T')[0];
}

export default function DailyLogPage() {
  const [date, setDate] = useState(getTodayStr());
  const [hours, setHours] = useState(
    Array.from({ length: 7 }, (_, i) => ({
      hour_slot: i + 1,
      activity_type: '',
      description: '',
      activity_points: 0,
      reward_points: 0,
      project_id: null,
    }))
  );
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);
  const [loading, setLoading] = useState(true);
  const [totalPoints, setTotalPoints] = useState({ activity: 0, reward: 0 });

  useEffect(() => {
    fetchDayLog();
  }, [date]);

  const fetchDayLog = async () => {
    setLoading(true);
    setSaved(false);
    try {
      const res = await api.get(`/activities?date=${date}`);
      const existing = res.data.activities || [];

      const merged = Array.from({ length: 7 }, (_, i) => {
        const found = existing.find(a => a.hour_slot === i + 1);
        if (found) {
          return {
            hour_slot: i + 1,
            activity_type: found.activity_type,
            description: found.description,
            activity_points: found.activity_points,
            reward_points: found.reward_points,
            project_id: found.project_id,
          };
        }
        return {
          hour_slot: i + 1,
          activity_type: '',
          description: '',
          activity_points: 0,
          reward_points: 0,
          project_id: null,
        };
      });

      setHours(merged);
      setTotalPoints({
        activity: res.data.total_activity_points || 0,
        reward: res.data.total_reward_points || 0,
      });
    } catch (err) {
      console.error('Failed to fetch day log:', err);
    }
    setLoading(false);
  };

  const updateHour = (index, field, value) => {
    setHours(prev => {
      const updated = [...prev];
      const oldType = updated[index].activity_type;
      updated[index] = { ...updated[index], [field]: value };

      // Auto-suggest points ONLY when activity type first changes (from empty or different type)
      if (field === 'activity_type' && value && value !== oldType) {
        const autoPoints = DEFAULT_POINTS[value] ?? 0;
        updated[index].activity_points = autoPoints;
        // Clear reward points if switching away from PS Slot
        if (value !== 'ps_slot') {
          updated[index].reward_points = 0;
        }
      }

      // If type is cleared, reset everything
      if (field === 'activity_type' && !value) {
        updated[index].activity_points = 0;
        updated[index].reward_points = 0;
      }

      return updated;
    });
    setSaved(false);
  };

  const handleSave = async () => {
    const validActivities = hours.filter(h => h.activity_type && h.description.trim());
    if (validActivities.length === 0) {
      alert('Please fill in at least one activity with type and description.');
      return;
    }

    setSaving(true);
    try {
      await api.post('/activities/bulk', {
        date,
        activities: validActivities,
      });
      setSaved(true);
      fetchDayLog();
    } catch (err) {
      alert('Failed to save activities: ' + (err.response?.data?.error || 'Unknown error'));
    }
    setSaving(false);
  };

  const filledCount = hours.filter(h => h.activity_type && h.description.trim()).length;
  const calcTotal = hours.reduce((sum, h) => sum + (h.activity_points || 0), 0);
  const calcReward = hours.reduce((sum, h) => sum + (h.reward_points || 0), 0);

  return (
    <div className="animate-fade-in">
      <div className="page-header">
        <div>
          <h1>Daily Activity Log</h1>
          <p>Log your 7 hours of activity for the day</p>
        </div>
        <div className="log-header-actions">
          <input
            type="date"
            value={date}
            onChange={(e) => setDate(e.target.value)}
            className="date-picker"
          />
        </div>
      </div>

      {/* Summary Bar */}
      <div className="log-summary">
        <div className="log-summary-item">
          <span className="log-summary-icon">🕐</span>
          <div>
            <span className="log-summary-value">{filledCount}/7</span>
            <span className="log-summary-label">Hours Filled</span>
          </div>
        </div>
        <div className="log-summary-item">
          <span className="log-summary-icon">📊</span>
          <div>
            <span className="log-summary-value">{calcTotal.toFixed(1)}</span>
            <span className="log-summary-label">Activity Points</span>
          </div>
        </div>
        <div className="log-summary-item">
          <span className="log-summary-icon">⭐</span>
          <div>
            <span className="log-summary-value">{calcReward.toFixed(1)}</span>
            <span className="log-summary-label">Reward Points (PS)</span>
          </div>
        </div>
        <div className="log-summary-save">
          <Button onClick={handleSave} loading={saving} size="lg">
            {saved ? '✓ Saved!' : 'Save All'}
          </Button>
        </div>
      </div>

      {/* Hour Grid */}
      {loading ? (
        <div className="page-loading"><div className="spinner spinner-lg"></div></div>
      ) : (
        <div className="hour-grid">
          {hours.map((hour, index) => (
            <div
              key={index}
              className={`hour-row ${hour.activity_type ? 'hour-filled' : ''} ${saved && hour.activity_type ? 'hour-saved' : ''}`}
              style={{ animationDelay: `${index * 50}ms` }}
            >
              <div className="hour-number">
                <span className="hour-slot">{index + 1}</span>
                <span className="hour-label">{HOUR_LABELS[index]}</span>
              </div>

              <div className="hour-fields">
                <div className="hour-type">
                  <Select
                    value={hour.activity_type}
                    onChange={(e) => updateHour(index, 'activity_type', e.target.value)}
                  >
                    {ACTIVITY_TYPES.map(t => (
                      <option key={t.value} value={t.value}>{t.label}</option>
                    ))}
                  </Select>
                </div>

                <div className="hour-desc">
                  <input
                    type="text"
                    className="input"
                    placeholder="What did you work on?"
                    value={hour.description}
                    onChange={(e) => updateHour(index, 'description', e.target.value)}
                  />
                </div>

                <div className="hour-points">
                  <input
                    type="number"
                    className="input points-input"
                    placeholder="Act Pts"
                    title="Activity Points"
                    value={hour.activity_points}
                    onChange={(e) => updateHour(index, 'activity_points', parseFloat(e.target.value) || 0)}
                    step="0.25"
                    min="0"
                  />
                </div>

                {hour.activity_type === 'ps_slot' && (
                  <div className="hour-reward">
                    <input
                      type="number"
                      className="input points-input reward-input"
                      placeholder="Reward"
                      title="Reward Points (PS Slot only)"
                      value={hour.reward_points}
                      onChange={(e) => updateHour(index, 'reward_points', parseFloat(e.target.value) || 0)}
                      step="0.5"
                      min="0"
                    />
                  </div>
                )}
              </div>

              {hour.activity_type && (
                <div className="hour-badge">
                  <Badge variant={hour.description ? 'success' : 'warning'} size="sm">
                    {hour.description ? 'Ready' : 'Need desc'}
                  </Badge>
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
