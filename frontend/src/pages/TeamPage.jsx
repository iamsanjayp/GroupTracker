import { useState, useEffect } from 'react';
import api from '../lib/api';
import useAuthStore from '../store/authStore';
import Card from '../components/ui/Card';
import Button from '../components/ui/Button';
import Badge from '../components/ui/Badge';
import Input from '../components/ui/Input';
import Modal from '../components/ui/Modal';
import { Select } from '../components/ui/Input';
import './TeamPage.css';

export default function TeamPage() {
  const { user, isAdmin, isCaptainVC, isPending, hasTeam, refreshUser } = useAuthStore();
  const [team, setTeam] = useState(null);
  const [members, setMembers] = useState([]);
  const [pendingMembers, setPendingMembers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showJoinModal, setShowJoinModal] = useState(false);
  const [showRoleModal, setShowRoleModal] = useState(false);
  const [teamName, setTeamName] = useState('');
  const [inviteCode, setInviteCode] = useState('');
  const [selectedMember, setSelectedMember] = useState(null);
  const [newRole, setNewRole] = useState('');
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    if (hasTeam()) {
      fetchTeamData();
    } else {
      setLoading(false);
    }
  }, []);

  const fetchTeamData = async () => {
    setLoading(true);
    try {
      const requests = [
        api.get('/teams/me'),
        api.get('/teams/members')
      ];
      
      if (isCaptainVC()) {
        requests.push(api.get('/teams/pending'));
      }

      const responses = await Promise.all(requests);
      setTeam(responses[0].data.team);
      setMembers(responses[1].data || []);
      
      if (isCaptainVC() && responses[2]) {
        setPendingMembers(responses[2].data || []);
      }
    } catch (err) { console.error(err); }
    setLoading(false);
  };

  const handleCreateTeam = async () => {
    if (!teamName.trim()) return;
    try {
      const res = await api.post('/teams', { name: teamName });
      // Save new tokens with updated team_id
      localStorage.setItem('access_token', res.data.access_token);
      localStorage.setItem('refresh_token', res.data.refresh_token);
      localStorage.setItem('user', JSON.stringify(res.data.user));
      setShowCreateModal(false);
      setTeamName('');
      window.location.reload();
    } catch (err) {
      alert(err.response?.data?.error || 'Failed to create team');
    }
  };

  const handleJoinTeam = async () => {
    if (!inviteCode.trim()) return;
    try {
      const res = await api.post('/teams/join', { invite_code: inviteCode });
      // Save new tokens with updated team_id
      localStorage.setItem('access_token', res.data.access_token);
      localStorage.setItem('refresh_token', res.data.refresh_token);
      localStorage.setItem('user', JSON.stringify(res.data.user));
      setShowJoinModal(false);
      setInviteCode('');
      window.location.reload();
    } catch (err) {
      alert(err.response?.data?.error || 'Invalid invite code');
    }
  };

  const handleUpdateRole = async () => {
    if (!selectedMember || !newRole) return;
    try {
      await api.put(`/teams/members/${selectedMember.id}/role`, { role: newRole });
      setShowRoleModal(false);
      fetchTeamData();
    } catch (err) {
      alert(err.response?.data?.error || 'Failed to update role');
    }
  };

  const handleRemoveMember = async (memberId) => {
    if (!confirm('Are you sure you want to remove this member?')) return;
    try {
      await api.delete(`/teams/members/${memberId}`);
      fetchTeamData();
    } catch (err) {
      alert(err.response?.data?.error || 'Failed to remove member');
    }
  };

  const handleApproveMember = async (memberId) => {
    try {
      await api.put(`/teams/members/${memberId}/approve`);
      fetchTeamData();
    } catch (err) {
      alert(err.response?.data?.error || 'Failed to approve member');
    }
  };

  const copyInviteCode = () => {
    navigator.clipboard.writeText(team?.invite_code || '');
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const roleBadgeVariant = (role) => {
    const map = { captain: 'primary', vice_captain: 'info', manager: 'warning', strategist: 'success', member: 'default' };
    return map[role] || 'default';
  };

  if (loading) {
    return <div className="page-loading"><div className="spinner spinner-lg"></div></div>;
  }

  if (isPending()) {
    return (
      <div className="animate-fade-in" style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', height: '60vh', textAlign: 'center' }}>
        <div style={{ fontSize: '4rem', marginBottom: '20px' }}>⏳</div>
        <h2>Waiting for Approval</h2>
        <p className="text-muted" style={{ maxWidth: '400px', marginTop: '10px' }}>
          You have successfully joined the team with code <strong>{user?.team_id}</strong>. 
          Please wait for the Captain or Vice Captain to approve your join request.
        </p>
        <Button variant="secondary" onClick={refreshUser} style={{ marginTop: '20px' }}>Refresh Status</Button>
      </div>
    );
  }

  if (!hasTeam()) {
    return (
      <div className="animate-fade-in">
        <div className="page-header">
          <h1>Team Setup</h1>
          <p>Create a new team or join an existing one</p>
        </div>

        <div className="team-setup-grid">
          <Card hover className="team-setup-card" onClick={() => setShowCreateModal(true)}>
            <span className="setup-icon">🚀</span>
            <h3>Create a Team</h3>
            <p>Start a new team and invite members</p>
            <Button variant="primary" size="sm">Create Team</Button>
          </Card>

          <Card hover className="team-setup-card" onClick={() => setShowJoinModal(true)}>
            <span className="setup-icon">🔗</span>
            <h3>Join a Team</h3>
            <p>Enter an invite code to join</p>
            <Button variant="secondary" size="sm">Join Team</Button>
          </Card>
        </div>

        <Modal isOpen={showCreateModal} onClose={() => setShowCreateModal(false)} title="Create Team" size="sm">
          <div className="modal-form">
            <Input label="Team Name" value={teamName} onChange={e => setTeamName(e.target.value)} placeholder="Enter team name" />
            <Button fullWidth onClick={handleCreateTeam}>Create</Button>
          </div>
        </Modal>

        <Modal isOpen={showJoinModal} onClose={() => setShowJoinModal(false)} title="Join Team" size="sm">
          <div className="modal-form">
            <Input label="Invite Code" value={inviteCode} onChange={e => setInviteCode(e.target.value)} placeholder="Paste invite code" />
            <Button fullWidth onClick={handleJoinTeam}>Join</Button>
          </div>
        </Modal>
      </div>
    );
  }

  return (
    <div className="animate-fade-in">
      <div className="page-header">
        <div>
          <h1>{team?.name}</h1>
          <p>Manage your team members and roles</p>
        </div>
      </div>

      {/* Invite Code Card */}
      <Card className="invite-card">
        <div className="invite-content">
          <div>
            <h4>Invite Code</h4>
            <p className="text-sm text-muted">Share this code with people to join your team</p>
          </div>
          <div className="invite-code-wrapper">
            <code className="invite-code">{team?.invite_code}</code>
            <Button size="sm" variant={copied ? 'success' : 'secondary'} onClick={copyInviteCode}>
              {copied ? '✓ Copied' : 'Copy'}
            </Button>
          </div>
        </div>
      </Card>

      {/* Pending Approvals */}
      {isCaptainVC() && pendingMembers.length > 0 && (
        <>
          <div className="section-header" style={{ marginTop: '30px' }}>
            <h2>Pending Approvals ({pendingMembers.length})</h2>
          </div>
          <div className="member-list pending-list">
            {pendingMembers.map(member => (
              <Card key={member.id} className="member-row" style={{ borderLeftColor: 'var(--warning)' }}>
                <div className="member-row-left">
                  <div className="member-row-avatar warning-bg">{member.name?.charAt(0)?.toUpperCase()}</div>
                  <div className="member-row-info">
                    <span className="member-row-name">{member.name}</span>
                    <span className="member-row-email text-muted">Roll No: {member.roll_no || 'N/A'}</span>
                  </div>
                </div>
                <div className="member-row-right">
                  <div className="member-actions">
                    <Button size="sm" variant="primary" onClick={() => handleApproveMember(member.id)}>
                      Approve
                    </Button>
                    <Button size="sm" variant="ghost" onClick={() => handleRemoveMember(member.id)} style={{ color: 'var(--error)' }}>
                      Reject
                    </Button>
                  </div>
                </div>
              </Card>
            ))}
          </div>
        </>
      )}

      {/* Members */}
      <div className="section-header" style={{ marginTop: isCaptainVC() && pendingMembers.length > 0 ? '30px' : '0' }}>
        <h2>Active Members ({members.filter(m => m.join_status !== 'pending').length})</h2>
      </div>

      <div className="member-list">
        {members.filter(m => m.join_status !== 'pending').map(member => (
          <Card key={member.id} className="member-row">
            <div className="member-row-left">
              <div className="member-row-avatar">{member.name?.charAt(0)?.toUpperCase()}</div>
              <div className="member-row-info">
                <span className="member-row-name">
                  {member.name}
                  {member.id === user?.id && <span className="you-badge">you</span>}
                </span>
                <span className="member-row-email">{member.email}</span>
              </div>
            </div>
            <div className="member-row-right">
              <Badge variant={roleBadgeVariant(member.role)} size="md">
                {member.role?.replace('_', ' ')}
              </Badge>
              {isAdmin() && member.id !== user?.id && (
                <div className="member-actions">
                  <Button size="sm" variant="ghost" onClick={() => { setSelectedMember(member); setNewRole(member.role); setShowRoleModal(true); }}>
                    Edit Role
                  </Button>
                  <Button size="sm" variant="ghost" onClick={() => handleRemoveMember(member.id)}>
                    ✕
                  </Button>
                </div>
              )}
            </div>
          </Card>
        ))}
      </div>

      {/* Role Modal */}
      <Modal isOpen={showRoleModal} onClose={() => setShowRoleModal(false)} title={`Change Role — ${selectedMember?.name}`} size="sm">
        <div className="modal-form">
          <Select label="Role" value={newRole} onChange={e => setNewRole(e.target.value)}>
            <option value="captain">Captain</option>
            <option value="vice_captain">Vice Captain</option>
            <option value="manager">Manager</option>
            <option value="strategist">Strategist</option>
            <option value="member">Member</option>
          </Select>
          <Button fullWidth onClick={handleUpdateRole}>Update Role</Button>
        </div>
      </Modal>
    </div>
  );
}
