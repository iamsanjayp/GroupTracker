import { useState, useEffect } from 'react';
import api from '../lib/api';
import useAuthStore from '../store/authStore';
import Card from '../components/ui/Card';
import Button from '../components/ui/Button';
import Badge from '../components/ui/Badge';
import Modal from '../components/ui/Modal';
import './SkillsPage.css';

const CATEGORY_LABELS = {
  primary: { label: 'Primary Skills', emoji: '💻', desc: 'Your core programming languages' },
  secondary: { label: 'Secondary Skills', emoji: '🛠️', desc: 'Your technical domains' },
  special: { label: 'Special Skills', emoji: '🌟', desc: 'Your soft skills & competencies' },
};

export default function SkillsPage() {
  const { user } = useAuthStore();
  const isCaptain = user?.role === 'captain';

  const [options, setOptions] = useState({});
  const [mySkills, setMySkills] = useState([]);
  const [teamSkills, setTeamSkills] = useState({});
  const [members, setMembers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [hasSkills, setHasSkills] = useState(false);

  // Selection state
  const [showModal, setShowModal] = useState(false);
  const [editingUserId, setEditingUserId] = useState(null);
  const [editingUserName, setEditingUserName] = useState('');
  const [selected, setSelected] = useState({ primary: [], secondary: [], special: [] });
  const [saving, setSaving] = useState(false);

  useEffect(() => { fetchData(); }, []);

  const fetchData = async () => {
    setLoading(true);
    try {
      const [optRes, myRes, teamRes, membersRes] = await Promise.all([
        api.get('/skills/options'),
        api.get('/skills/me'),
        api.get('/skills/team').catch(() => ({ data: {} })),
        api.get('/teams/members').catch(() => ({ data: [] })),
      ]);

      setOptions(optRes.data || {});
      const skills = myRes.data?.skills || [];
      setMySkills(skills);
      setHasSkills(skills.length >= 6);
      setTeamSkills(teamRes.data || {});
      setMembers(membersRes.data || []);
    } catch (err) {
      console.error(err);
    }
    setLoading(false);
  };

  const openSetSkills = (userId = null, userName = '') => {
    const targetId = userId || user?.id;
    setEditingUserId(targetId);
    setEditingUserName(userName || 'Your');

    // Pre-fill with existing selections
    const existingSkills = userId ? (teamSkills[userId] || []) : mySkills;
    const prefill = { primary: [], secondary: [], special: [] };
    existingSkills.forEach(s => {
      if (prefill[s.category]) {
        prefill[s.category].push(s.skill_name);
      }
    });
    setSelected(prefill);
    setShowModal(true);
  };

  const toggleSkill = (category, name) => {
    setSelected(prev => {
      const curr = [...prev[category]];
      const idx = curr.indexOf(name);
      if (idx >= 0) {
        curr.splice(idx, 1);
      } else if (curr.length < 2) {
        curr.push(name);
      }
      return { ...prev, [category]: curr };
    });
  };

  const handleSave = async () => {
    if (selected.primary.length !== 2 || selected.secondary.length !== 2 || selected.special.length !== 2) {
      alert('Please select exactly 2 skills in each category.');
      return;
    }

    setSaving(true);
    try {
      if (editingUserId && editingUserId !== user?.id) {
        await api.put(`/skills/member/${editingUserId}`, selected);
      } else {
        await api.post('/skills/me', selected);
      }
      setShowModal(false);
      fetchData();
    } catch (err) {
      alert(err.response?.data?.error || 'Failed to save skills');
    }
    setSaving(false);
  };

  if (loading) {
    return <div className="page-loading"><div className="spinner spinner-lg"></div></div>;
  }

  const groupedSkills = (skills) => {
    const groups = { primary: [], secondary: [], special: [] };
    (skills || []).forEach(s => {
      if (groups[s.category]) groups[s.category].push(s);
    });
    return groups;
  };

  const myGrouped = groupedSkills(mySkills);

  return (
    <div className="animate-fade-in skills-page">
      <div className="page-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div>
          <h1>Skills</h1>
          <p>Your 6 skill picks across 3 categories</p>
        </div>
        {(!hasSkills || isCaptain) && (
          <Button variant="primary" onClick={() => openSetSkills()}>
            {hasSkills ? 'Edit My Skills' : 'Set My Skills'}
          </Button>
        )}
      </div>

      {/* My Skills Display */}
      {hasSkills ? (
        <div className="skills-grid">
          {Object.entries(CATEGORY_LABELS).map(([cat, info]) => (
            <Card key={cat} className="skill-category-card">
              <div className="skill-cat-header">
                <span className="skill-cat-emoji">{info.emoji}</span>
                <div>
                  <h3>{info.label}</h3>
                  <p className="text-sm text-muted">{info.desc}</p>
                </div>
              </div>
              <div className="skill-tags">
                {(myGrouped[cat] || []).map(s => (
                  <span key={s.skill_name} className={`skill-tag skill-tag-${cat}`}>
                    {s.skill_name}
                  </span>
                ))}
                {(myGrouped[cat] || []).length === 0 && (
                  <span className="text-muted text-sm">Not set</span>
                )}
              </div>
            </Card>
          ))}
        </div>
      ) : (
        <Card style={{ textAlign: 'center', padding: '3rem' }}>
          <div style={{ fontSize: '3rem', marginBottom: '1rem' }}>🎯</div>
          <h3>Set Your Skills</h3>
          <p className="text-muted" style={{ marginTop: '0.5rem', maxWidth: '400px', margin: '0.5rem auto 1.5rem' }}>
            Choose 2 primary, 2 secondary, and 2 special skills. These determine your project allocation and class participation focus.
          </p>
          <Button variant="primary" onClick={() => openSetSkills()}>Choose Skills</Button>
        </Card>
      )}

      {/* Captain: Team Skills Overview */}
      {isCaptain && members.length > 0 && (
        <>
          <div className="section-header" style={{ marginTop: '2rem' }}>
            <h2>👥 Team Skills</h2>
          </div>
          <div className="team-skills-list">
            {members.filter(m => m.join_status !== 'pending').map(member => {
              const memberSkills = teamSkills[member.id] || [];
              const mGrouped = groupedSkills(memberSkills);
              return (
                <Card key={member.id} className="team-skill-row">
                  <div className="team-skill-header">
                    <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                      <div className="member-row-avatar">{member.name?.charAt(0)?.toUpperCase()}</div>
                      <div>
                        <strong>{member.name}</strong>
                        {member.id === user?.id && <span className="you-badge">you</span>}
                        <div className="text-xs text-muted">{member.roll_no || member.email}</div>
                      </div>
                    </div>
                    <Button size="sm" variant="ghost" onClick={() => openSetSkills(member.id, member.name)}>
                      {memberSkills.length >= 6 ? 'Edit' : 'Set Skills'}
                    </Button>
                  </div>
                  {memberSkills.length > 0 ? (
                    <div className="team-skill-tags">
                      {Object.entries(CATEGORY_LABELS).map(([cat, info]) => (
                        <div key={cat} className="team-skill-cat">
                          <span className="text-xs text-muted">{info.emoji} {info.label}:</span>
                          <div className="skill-tags-inline">
                            {(mGrouped[cat] || []).map(s => (
                              <Badge key={s.skill_name} variant={cat === 'primary' ? 'primary' : cat === 'secondary' ? 'info' : 'success'} size="sm">
                                {s.skill_name}
                              </Badge>
                            ))}
                            {(mGrouped[cat] || []).length === 0 && <span className="text-xs text-muted">—</span>}
                          </div>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <p className="text-sm text-muted" style={{ marginTop: '8px' }}>No skills selected yet</p>
                  )}
                </Card>
              );
            })}
          </div>
        </>
      )}

      {/* Skill Selection Modal */}
      <Modal isOpen={showModal} onClose={() => setShowModal(false)} title={`${editingUserName === 'Your' ? 'Set Your' : `Set ${editingUserName}'s`} Skills`} size="lg">
        <div className="skill-selection-modal">
          {Object.entries(CATEGORY_LABELS).map(([cat, info]) => (
            <div key={cat} className="skill-selection-group">
              <div className="skill-selection-header">
                <h4>{info.emoji} {info.label}</h4>
                <span className="text-sm text-muted">
                  {selected[cat]?.length || 0}/2 selected
                </span>
              </div>
              <p className="text-sm text-muted" style={{ marginBottom: '12px' }}>{info.desc}</p>
              <div className="skill-option-grid">
                {(options[cat] || []).map(name => {
                  const isSelected = selected[cat]?.includes(name);
                  const isDisabled = !isSelected && selected[cat]?.length >= 2;
                  return (
                    <button
                      key={name}
                      className={`skill-option ${isSelected ? 'selected' : ''} ${isDisabled ? 'disabled' : ''}`}
                      onClick={() => !isDisabled && toggleSkill(cat, name)}
                      disabled={isDisabled}
                    >
                      {isSelected && <span className="skill-check">✓</span>}
                      {name}
                    </button>
                  );
                })}
              </div>
            </div>
          ))}
          <div style={{ display: 'flex', gap: '10px', marginTop: '20px' }}>
            <Button fullWidth variant="ghost" onClick={() => setShowModal(false)}>Cancel</Button>
            <Button fullWidth onClick={handleSave} loading={saving}
              disabled={selected.primary?.length !== 2 || selected.secondary?.length !== 2 || selected.special?.length !== 2}>
              Save Skills
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}
