import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import useAuthStore from '../store/authStore';
import Button from '../components/ui/Button';
import Input from '../components/ui/Input';
import './LoginPage.css';

export default function LoginPage() {
  const [isRegister, setIsRegister] = useState(false);
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [rollNo, setRollNo] = useState('');
  const [error, setError] = useState('');
  const { login, register, loading } = useAuthStore();
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    let result;
    if (isRegister) {
      if (!name.trim()) { setError('Name is required'); return; }
      if (!rollNo.trim()) { setError('Roll Number is required'); return; }
      result = await register(name, email, password, rollNo.trim().toUpperCase());
    } else {
      result = await login(email, password);
    }

    if (result.success) {
      navigate('/dashboard');
    } else {
      setError(result.error);
    }
  };



  return (
    <div className="login-page">
      {/* Background decoration */}
      <div className="login-bg">
        <div className="login-bg-shape shape-1"></div>
        <div className="login-bg-shape shape-2"></div>
        <div className="login-bg-shape shape-3"></div>
      </div>

      <div className="login-container">
        {/* Left side - Branding */}
        <div className="login-hero">
          <div className="login-hero-content">
            <div className="login-hero-logo">
              <span>⚡</span>
            </div>
            <h1>GroupTracker</h1>
            <p>Team Productivity & Project-Based Learning Platform</p>
            <div className="login-features">
              <div className="login-feature">
                <span className="feature-icon">📊</span>
                <span>Track daily activities</span>
              </div>
              <div className="login-feature">
                <span className="feature-icon">📁</span>
                <span>Manage team projects</span>
              </div>
              <div className="login-feature">
                <span className="feature-icon">⭐</span>
                <span>Earn activity points</span>
              </div>
              <div className="login-feature">
                <span className="feature-icon">👥</span>
                <span>Collaborate with your team</span>
              </div>
            </div>
          </div>
        </div>

        {/* Right side - Form */}
        <div className="login-form-container">
          <div className="login-form-wrapper">
            <div className="login-form-header">
              <h2>{isRegister ? 'Create Account' : 'Welcome Back'}</h2>
              <p>{isRegister ? 'Join your team on GroupTracker' : 'Sign in to your account'}</p>
            </div>

            {error && (
              <div className="login-error">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <circle cx="12" cy="12" r="10" /><line x1="15" y1="9" x2="9" y2="15" /><line x1="9" y1="9" x2="15" y2="15" />
                </svg>
                {error}
              </div>
            )}

            <form onSubmit={handleSubmit} className="login-form">
              {isRegister && (
                <>
                  <Input
                    label="Full Name"
                    type="text"
                    placeholder="Enter your name"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    required
                  />
                  <Input
                    label="Roll Number"
                    type="text"
                    placeholder="Ex: 21CSR123"
                    value={rollNo}
                    onChange={(e) => setRollNo(e.target.value)}
                    required
                  />
                </>
              )}
              <Input
                label="Email"
                type="email"
                placeholder="Enter your email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
              <Input
                label="Password"
                type="password"
                placeholder={isRegister ? 'Create a password (min 6 chars)' : 'Enter your password'}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
              <Button type="submit" fullWidth loading={loading} size="lg">
                {isRegister ? 'Create Account' : 'Sign In'}
              </Button>
            </form>



            <p className="login-toggle">
              {isRegister ? 'Already have an account?' : "Don't have an account?"}
              <button onClick={() => { setIsRegister(!isRegister); setError(''); }}>
                {isRegister ? 'Sign In' : 'Create Account'}
              </button>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
