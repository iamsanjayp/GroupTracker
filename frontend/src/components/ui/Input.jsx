import './Input.css';

export default function Input({ label, error, icon, ...props }) {
  return (
    <div className={`input-group ${error ? 'input-error' : ''}`}>
      {label && <label className="input-label">{label}</label>}
      <div className="input-wrapper">
        {icon && <span className="input-icon">{icon}</span>}
        <input className={`input ${icon ? 'input-with-icon' : ''}`} {...props} />
      </div>
      {error && <span className="input-error-msg">{error}</span>}
    </div>
  );
}

export function Select({ label, error, children, ...props }) {
  return (
    <div className={`input-group ${error ? 'input-error' : ''}`}>
      {label && <label className="input-label">{label}</label>}
      <select className="input select" {...props}>
        {children}
      </select>
      {error && <span className="input-error-msg">{error}</span>}
    </div>
  );
}

export function TextArea({ label, error, ...props }) {
  return (
    <div className={`input-group ${error ? 'input-error' : ''}`}>
      {label && <label className="input-label">{label}</label>}
      <textarea className="input textarea" {...props} />
      {error && <span className="input-error-msg">{error}</span>}
    </div>
  );
}
