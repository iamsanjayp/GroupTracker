import './ProgressBar.css';

export default function ProgressBar({ value = 0, max = 100, label, showValue = true, variant = 'primary' }) {
  const percent = Math.min(Math.round((value / max) * 100), 100);
  
  return (
    <div className="progress-container">
      {(label || showValue) && (
        <div className="progress-header">
          {label && <span className="progress-label">{label}</span>}
          {showValue && <span className="progress-value">{percent}%</span>}
        </div>
      )}
      <div className="progress-track">
        <div 
          className={`progress-fill progress-${variant}`} 
          style={{ width: `${percent}%` }}
        />
      </div>
    </div>
  );
}
