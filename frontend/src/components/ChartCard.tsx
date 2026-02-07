interface ChartCardProps {
  title: string;
  subtitle?: string;
  children: React.ReactNode;
}

export function ChartCard({ title, subtitle, children }: ChartCardProps) {
  return (
    <div className="bg-slate-800 border border-slate-700 rounded-lg p-6">
      <div className="mb-4">
        <h3 className="text-lg font-semibold text-slate-50">{title}</h3>
        {subtitle && (
          <p className="text-sm text-slate-400 mt-1">{subtitle}</p>
        )}
      </div>
      {children}
    </div>
  );
}
