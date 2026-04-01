import { useState, useEffect } from 'react';
import api from '../lib/api';
import useAuthStore from '../store/authStore';
import Card from '../components/ui/Card';
import Button from '../components/ui/Button';
import Badge from '../components/ui/Badge';
import Modal from '../components/ui/Modal';
import Input, { TextArea, Select } from '../components/ui/Input';
import ProgressBar from '../components/ui/ProgressBar';
import './ProjectsPage.css';

export default function ProjectsPage() {
  const { isAdmin } = useAuthStore();
  const [projects, setProjects] = useState([]);
  const [selectedProject, setSelectedProject] = useState(null);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showTaskModal, setShowTaskModal] = useState(false);
  const [showMemberModal, setShowMemberModal] = useState(false);
  const [loading, setLoading] = useState(true);
  const [newProject, setNewProject] = useState({ name: '', description: '' });
  const [newTask, setNewTask] = useState({ title: '', description: '', priority: 'medium', due_date: '' });
  const [newMember, setNewMember] = useState({ user_id: '', share_percentage: 0 });
  const [teamMembers, setTeamMembers] = useState([]);

  useEffect(() => { fetchProjects(); fetchTeamMembers(); }, []);

  const fetchProjects = async () => {
    setLoading(true);
    try {
      const res = await api.get('/projects');
      setProjects(res.data || []);
    } catch (err) { console.error(err); }
    setLoading(false);
  };

  const fetchTeamMembers = async () => {
    try {
      const res = await api.get('/teams/members');
      setTeamMembers(res.data || []);
    } catch {}
  };

  const fetchProjectDetail = async (id) => {
    try {
      const res = await api.get(`/projects/${id}`);
      setSelectedProject(res.data);
    } catch (err) { console.error(err); }
  };

  const handleCreateProject = async () => {
    if (!newProject.name.trim()) return;
    try {
      await api.post('/projects', newProject);
      setShowCreateModal(false);
      setNewProject({ name: '', description: '' });
      fetchProjects();
    } catch (err) {
      alert(err.response?.data?.error || 'Failed to create project');
    }
  };

  const handleAddTask = async () => {
    if (!newTask.title.trim()) return;
    try {
      await api.post(`/projects/${selectedProject.id}/tasks`, newTask);
      setShowTaskModal(false);
      setNewTask({ title: '', description: '', priority: 'medium', due_date: '' });
      fetchProjectDetail(selectedProject.id);
    } catch (err) { alert(err.response?.data?.error || 'Failed'); }
  };

  const handleUpdateTaskStatus = async (taskId, status) => {
    try {
      await api.put(`/projects/${selectedProject.id}/tasks/${taskId}`, { status });
      fetchProjectDetail(selectedProject.id);
    } catch (err) { console.error(err); }
  };

  const handleAddMember = async () => {
    if (!newMember.user_id) return;
    try {
      await api.post(`/projects/${selectedProject.id}/members`, {
        user_id: parseInt(newMember.user_id),
        share_percentage: parseFloat(newMember.share_percentage) || 0,
      });
      setShowMemberModal(false);
      setNewMember({ user_id: '', share_percentage: 0 });
      fetchProjectDetail(selectedProject.id);
    } catch (err) { alert(err.response?.data?.error || 'Failed'); }
  };

  const statusBadge = (status) => {
    const map = { active: 'success', completed: 'primary', on_hold: 'warning' };
    return map[status] || 'default';
  };

  const priorityBadge = (priority) => {
    const map = { high: 'error', medium: 'warning', low: 'default' };
    return map[priority] || 'default';
  };

  const taskStatusOptions = ['todo', 'in_progress', 'review', 'done'];

  if (loading) {
    return <div className="page-loading"><div className="spinner spinner-lg"></div></div>;
  }

  if (selectedProject) {
    const tasks = selectedProject.tasks || [];
    const doneTasks = tasks.filter(t => t.status === 'done').length;
    const members = selectedProject.members || [];

    return (
      <div className="animate-fade-in">
        <div className="page-header">
          <div>
            <button className="back-btn" onClick={() => setSelectedProject(null)}>← Back to Projects</button>
            <h1>{selectedProject.name}</h1>
            <p>{selectedProject.description || 'No description'}</p>
          </div>
          <Badge variant={statusBadge(selectedProject.status)} size="md">{selectedProject.status}</Badge>
        </div>

        {/* Progress */}
        <Card className="project-progress-card">
          <ProgressBar value={doneTasks} max={tasks.length || 1} label="Task Completion" variant="success" />
        </Card>

        {/* Members */}
        <div className="section-header">
          <h2>Team Members ({members.length})</h2>
          {isAdmin() && <Button size="sm" onClick={() => setShowMemberModal(true)}>+ Add Member</Button>}
        </div>
        <div className="members-grid">
          {members.map(m => (
            <Card key={m.user_id} className="member-chip">
              <div className="member-avatar">{m.name?.charAt(0)}</div>
              <div className="member-info">
                <span className="member-name">{m.name}</span>
                <span className="member-share">{m.share_percentage}% share</span>
              </div>
            </Card>
          ))}
        </div>

        {/* Tasks */}
        <div className="section-header" style={{ marginTop: '28px' }}>
          <h2>Tasks ({tasks.length})</h2>
          {isAdmin() && <Button size="sm" onClick={() => setShowTaskModal(true)}>+ Add Task</Button>}
        </div>

        {tasks.length === 0 ? (
          <Card><p className="text-muted" style={{ textAlign: 'center', padding: '20px' }}>No tasks yet</p></Card>
        ) : (
          <div className="tasks-list">
            {tasks.map(task => (
              <div key={task.id} className={`task-item task-${task.status}`}>
                <div className="task-left">
                  <input
                    type="checkbox"
                    className="task-checkbox"
                    checked={task.status === 'done'}
                    onChange={() => handleUpdateTaskStatus(task.id, task.status === 'done' ? 'todo' : 'done')}
                  />
                  <div>
                    <span className={`task-title ${task.status === 'done' ? 'task-done' : ''}`}>{task.title}</span>
                    {task.assignee_name && <span className="task-assignee">→ {task.assignee_name}</span>}
                  </div>
                </div>
                <div className="task-right">
                  <Badge variant={priorityBadge(task.priority)} size="sm">{task.priority}</Badge>
                  <select
                    className="task-status-select"
                    value={task.status}
                    onChange={(e) => handleUpdateTaskStatus(task.id, e.target.value)}
                  >
                    {taskStatusOptions.map(s => (
                      <option key={s} value={s}>{s.replace('_', ' ')}</option>
                    ))}
                  </select>
                </div>
              </div>
            ))}
          </div>
        )}

        {/* Add Task Modal */}
        <Modal isOpen={showTaskModal} onClose={() => setShowTaskModal(false)} title="Add Task">
          <div className="modal-form">
            <Input label="Title" value={newTask.title} onChange={e => setNewTask({...newTask, title: e.target.value})} placeholder="Task title" required />
            <TextArea label="Description" value={newTask.description} onChange={e => setNewTask({...newTask, description: e.target.value})} placeholder="Task description" />
            <Select label="Priority" value={newTask.priority} onChange={e => setNewTask({...newTask, priority: e.target.value})}>
              <option value="low">Low</option>
              <option value="medium">Medium</option>
              <option value="high">High</option>
            </Select>
            <Input label="Due Date" type="date" value={newTask.due_date} onChange={e => setNewTask({...newTask, due_date: e.target.value})} />
            <Button fullWidth onClick={handleAddTask}>Create Task</Button>
          </div>
        </Modal>

        {/* Add Member Modal */}
        <Modal isOpen={showMemberModal} onClose={() => setShowMemberModal(false)} title="Add Member">
          <div className="modal-form">
            <Select label="Member" value={newMember.user_id} onChange={e => setNewMember({...newMember, user_id: e.target.value})}>
              <option value="">Select member...</option>
              {teamMembers.map(m => <option key={m.id} value={m.id}>{m.name}</option>)}
            </Select>
            <Input label="Share %" type="number" value={newMember.share_percentage} onChange={e => setNewMember({...newMember, share_percentage: e.target.value})} min="0" max="100" step="1" />
            <Button fullWidth onClick={handleAddMember}>Add to Project</Button>
          </div>
        </Modal>
      </div>
    );
  }

  return (
    <div className="animate-fade-in">
      <div className="page-header">
        <div>
          <h1>Projects</h1>
          <p>Manage your team's projects and tasks</p>
        </div>
        {isAdmin() && <Button onClick={() => setShowCreateModal(true)}>+ New Project</Button>}
      </div>

      {projects.length === 0 ? (
        <Card>
          <div style={{ textAlign: 'center', padding: '40px' }}>
            <span style={{ fontSize: '3rem' }}>📁</span>
            <h3 style={{ marginTop: '12px' }}>No projects yet</h3>
            <p className="text-muted">Create your first project to get started</p>
          </div>
        </Card>
      ) : (
        <div className="project-grid">
          {projects.map(p => (
            <Card key={p.id} hover className="project-card" onClick={() => fetchProjectDetail(p.id)}>
              <div className="project-card-header">
                <h3>{p.name}</h3>
                <Badge variant={statusBadge(p.status)} size="sm">{p.status}</Badge>
              </div>
              <p className="project-desc">{p.description || 'No description'}</p>
              <div className="project-card-footer">
                <span className="project-meta">👥 {p.member_count || 0} members</span>
                <span className="project-meta">📋 {p.task_count || 0} tasks</span>
              </div>
            </Card>
          ))}
        </div>
      )}

      {/* Create Project Modal */}
      <Modal isOpen={showCreateModal} onClose={() => setShowCreateModal(false)} title="Create Project">
        <div className="modal-form">
          <Input label="Project Name" value={newProject.name} onChange={e => setNewProject({...newProject, name: e.target.value})} placeholder="Enter project name" required />
          <TextArea label="Description" value={newProject.description} onChange={e => setNewProject({...newProject, description: e.target.value})} placeholder="What is this project about?" />
          <Button fullWidth onClick={handleCreateProject}>Create Project</Button>
        </div>
      </Modal>
    </div>
  );
}
