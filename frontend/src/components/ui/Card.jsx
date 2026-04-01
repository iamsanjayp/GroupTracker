import './Card.css';

export default function Card({ children, className = '', hover = false, padding = true, ...props }) {
  return (
    <div className={`card ${hover ? 'card-hover' : ''} ${padding ? '' : 'card-no-pad'} ${className}`} {...props}>
      {children}
    </div>
  );
}

export function CardHeader({ children, className = '' }) {
  return <div className={`card-header ${className}`}>{children}</div>;
}

export function CardBody({ children, className = '' }) {
  return <div className={`card-body ${className}`}>{children}</div>;
}
