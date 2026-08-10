import type { AssetStatus } from '@/types';

const statusClass: Record<string, string> = {
  available: 'bg-emerald-100 text-emerald-800',
  assigned: 'bg-sky-100 text-sky-800',
  in_transit: 'bg-amber-100 text-amber-900',
  maintenance: 'bg-orange-100 text-orange-900',
  reserved: 'bg-violet-100 text-violet-900',
  retired: 'bg-slate-200 text-slate-700',
  sold: 'bg-slate-200 text-slate-700',
};

export function Badge({ status, label }: { status: AssetStatus | string; label?: string }) {
  const cls = statusClass[status] ?? 'bg-slate-100 text-slate-700';
  return (
    <span className={`inline-flex rounded px-2 py-0.5 text-xs font-semibold uppercase tracking-wide ${cls}`}>
      {label ?? status.replace('_', ' ')}
    </span>
  );
}
