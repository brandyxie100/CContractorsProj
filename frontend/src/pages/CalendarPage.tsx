import { useMemo } from 'react';
import { Badge } from '@/components/ui/Badge';
import { Card } from '@/components/ui/Card';
import { useAssets, useAssignments } from '@/features/hooks/useApi';

function startOfWeek(): Date {
  const d = new Date();
  d.setHours(0, 0, 0, 0);
  d.setDate(d.getDate() - d.getDay() + 1);
  return d;
}

export function CalendarPage() {
  const from = startOfWeek();
  const to = new Date(from);
  to.setDate(to.getDate() + 7);
  const fromISO = from.toISOString();
  const toISO = to.toISOString();

  const { data: assets } = useAssets();
  const { data: assignments, isLoading } = useAssignments(fromISO, toISO);

  const days = useMemo(() => {
    return Array.from({ length: 7 }, (_, i) => {
      const d = new Date(from);
      d.setDate(from.getDate() + i);
      return d;
    });
  }, [from]);

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold text-[var(--color-navy)]">Availability Calendar</h1>
      <p className="text-sm text-slate-600">Week view — color bars show assignments per asset</p>
      {isLoading ? <p>Loading…</p> : null}
      <Card className="overflow-x-auto">
        <div className="min-w-[720px]">
          <div className="grid grid-cols-8 gap-1 border-b border-slate-200 pb-2 text-xs font-semibold text-slate-600">
            <div>Asset</div>
            {days.map((d) => (
              <div key={d.toISOString()}>{d.toLocaleDateString(undefined, { weekday: 'short', day: 'numeric' })}</div>
            ))}
          </div>
          <div className="mt-2 space-y-2">
            {(assets ?? []).map((asset) => {
              const rows = (assignments ?? []).filter((a) => a.asset_id === asset.id);
              return (
                <div key={asset.id} className="grid grid-cols-8 items-center gap-1 text-sm">
                  <div>
                    <div className="font-medium">{asset.asset_code}</div>
                    <Badge status={asset.status} />
                  </div>
                  {days.map((day) => {
                    const dayStart = new Date(day);
                    const dayEnd = new Date(day);
                    dayEnd.setDate(dayEnd.getDate() + 1);
                    const hit = rows.find((a) => {
                      const s = new Date(a.start_time);
                      const e = new Date(a.end_time);
                      return s < dayEnd && e > dayStart;
                    });
                    return (
                      <div
                        key={day.toISOString()}
                        className={`h-10 rounded border ${
                          hit ? 'border-sky-300 bg-sky-200' : 'border-slate-100 bg-slate-50'
                        }`}
                        title={hit?.project?.name ?? ''}
                      />
                    );
                  })}
                </div>
              );
            })}
          </div>
        </div>
      </Card>
    </div>
  );
}
