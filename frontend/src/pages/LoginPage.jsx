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

  const handleGoogleLogin = async () => {
    try {
      const res = await fetch('/api/auth/google');
      const data = await res.json();
      window.location.href = data.url;
    } catch {
      setError('Failed to initiate Google login');
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

            <div className="login-divider">
              <span>or</span>
            </div>

            <Button 
              variant="secondary" 
              fullWidth 
              size="lg" 
              onClick={handleGoogleLogin}
              icon={
                <svg width="18" height="18" viewBox="0 0 24 24">
                  <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 0 1-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z"/>
                  <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                  <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                  <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
                </svg>
              }
            >
              Continue with Google
            </Button>

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
