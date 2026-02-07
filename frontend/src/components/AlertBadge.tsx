interface AlertBadgeProps {
  level: 'LOW' | 'MEDIUM' | 'HIGH';
}

export function AlertBadge({ level }: AlertBadgeProps) {
  const colors = {
    LOW: 'bg-green-500/20 text-green-500 border-green-500/50',
    MEDIUM: 'bg-amber-500/20 text-amber-500 border-amber-500/50',
    HIGH: 'bg-red-500/20 text-red-500 border-red-500/50',
  };

  return (
    <span className={`inline-flex items-center px-3 py-1 rounded-full text-xs font-medium border ${colors[level]}`}>
      {level}
    </span>
  );
}
