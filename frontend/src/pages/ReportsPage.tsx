import { Card } from '@/components/ui/Card';
import { useAssets, useDashboard } from '@/features/hooks/useApi';

export function ReportsPage() {
  const { data: summary } = useDashboard();
  const { data: assets } = useAssets();

  const total = assets?.length ?? 0;
  const available = assets?.filter((a) => a.status === 'available').length ?? 0;
  const inUse = total - available;

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold text-[var(--color-navy)]">Reports</h1>
      <p className="text-sm text-slate-600">Basic utilization snapshot for fleet planning</p>
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <div className="text-xs uppercase text-slate-500">Fleet size</div>
          <div className="text-3xl font-semibold">{total}</div>
        </Card>
        <Card>
          <div className="text-xs uppercase text-slate-500">Available now</div>
          <div className="text-3xl font-semibold text-emerald-700">{available}</div>
        </Card>
        <Card>
          <div className="text-xs uppercase text-slate-500">In use / reserved</div>
          <div className="text-3xl font-semibold text-sky-800">{inUse}</div>
        </Card>
      </div>
      <Card>
        <h2 className="mb-2 font-semibold">Utilization</h2>
        <p className="text-4xl font-bold text-[var(--color-steel)]">
          {(summary?.utilization_percent ?? 0).toFixed(1)}%
        </p>
        <p className="mt-2 text-sm text-slate-600">
          Active assignments: {summary?.active_assignments ?? 0} · Upcoming week:{' '}
          {summary?.upcoming_assignments ?? 0}
        </p>
      </Card>
    </div>
  );
}
