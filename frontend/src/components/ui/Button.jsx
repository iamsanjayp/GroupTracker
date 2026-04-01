import './Button.css';

export default function Button({ 
  children, variant = 'primary', size = 'md', loading = false, 
  fullWidth = false, icon, ...props 
}) {
  return (
    <button
      className={`btn btn-${variant} btn-${size} ${fullWidth ? 'btn-full' : ''} ${loading ? 'btn-loading' : ''}`}
      disabled={loading || props.disabled}
      {...props}
    >
      {loading && <span className="spinner" />}
      {icon && !loading && <span className="btn-icon">{icon}</span>}
      {children}
    </button>
  );
}
