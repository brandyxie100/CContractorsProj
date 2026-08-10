import { useState } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { Badge } from '@/components/ui/Badge';
import { Button } from '@/components/ui/Button';
import { Card } from '@/components/ui/Card';
import { Input } from '@/components/ui/Input';
import { Select } from '@/components/ui/Select';
import { createAssignment, searchDispatch, useProjects } from '@/features/hooks/useApi';
import type { DispatchResult } from '@/types';

export function DispatchPage() {
  const { data: projects } = useProjects('active');
  const qc = useQueryClient();
  const [equipmentType, setEquipmentType] = useState('excavator');
  const [start, setStart] = useState(() => new Date().toISOString().slice(0, 16));
  const [end, setEnd] = useState(() => {
    const d = new Date();
    d.setDate(d.getDate() + 2);
    return d.toISOString().slice(0, 16);
  });
  const [projectId, setProjectId] = useState('');
  const [minWeight, setMinWeight] = useState('20');
  const [results, setResults] = useState<DispatchResult[]>([]);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  async function onSearch(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError('');
    try {
      const body: Record<string, unknown> = {
        equipment_type: equipmentType,
        start_time: new Date(start).toISOString(),
        end_time: new Date(end).toISOString(),
        include_alternatives: true,
      };
      if (projectId) body.project_id = projectId;
      if (minWeight) body.min_weight_t = Number(minWeight);
      const res = await searchDispatch(body);
      setResults(res.results);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Search failed');
    } finally {
      setLoading(false);
    }
  }

  async function assign(item: DispatchResult) {
    if (!projectId) {
      setError('Select a project before assigning');
      return;
    }
    try {
      await createAssignment({
        asset_id: item.asset.id,
        project_id: projectId,
        start_time: new Date(start).toISOString(),
        end_time: new Date(end).toISOString(),
        notes: 'Assigned via Smart Dispatch Finder',
      });
      await qc.invalidateQueries({ queryKey: ['assignments'] });
      await qc.invalidateQueries({ queryKey: ['assets'] });
      await onSearch({ preventDefault() {} } as React.FormEvent);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Assign failed');
    }
  }

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold text-[var(--color-navy)]">Smart Dispatch Finder</h1>
      <Card>
        <form className="grid gap-3 md:grid-cols-2 lg:grid-cols-3" onSubmit={(e) => void onSearch(e)}>
          <Select
            label="Equipment type"
            value={equipmentType}
            onChange={(e) => setEquipmentType(e.target.value)}
            options={[
              { value: 'excavator', label: 'Excavator' },
              { value: 'truck', label: 'Truck' },
              { value: 'roller', label: 'Roller' },
            ]}
          />
          <Input label="Start" type="datetime-local" value={start} onChange={(e) => setStart(e.target.value)} />
          <Input label="End" type="datetime-local" value={end} onChange={(e) => setEnd(e.target.value)} />
          <Select
            label="Project"
            value={projectId}
            onChange={(e) => setProjectId(e.target.value)}
            options={[
              { value: '', label: 'Select project…' },
              ...(projects ?? []).map((p) => ({ value: p.id, label: p.name })),
            ]}
          />
          <Input label="Min weight (t)" value={minWeight} onChange={(e) => setMinWeight(e.target.value)} />
          <div className="flex items-end">
            <Button type="submit" disabled={loading}>
              {loading ? 'Searching…' : 'Find equipment'}
            </Button>
          </div>
        </form>
        {error ? <p className="mt-3 text-sm text-red-600">{error}</p> : null}
      </Card>

      <div className="space-y-3">
        {results.map((r) => (
          <Card key={r.asset.id}>
            <div className="flex flex-wrap items-start justify-between gap-3">
              <div>
                <div className="font-semibold">
                  {r.asset.asset_code} — {r.asset.name}
                  {r.is_alternative ? (
                    <span className="ml-2 text-xs font-normal text-amber-700">(alternative)</span>
                  ) : null}
                </div>
                <div className="mt-1 flex flex-wrap gap-2 text-sm text-slate-600">
                  <Badge status={r.available ? 'available' : 'reserved'} label={r.available ? 'available now' : 'later'} />
                  <span>From: {new Date(r.available_from).toLocaleString()}</span>
                  {r.distance_km != null ? <span>{r.distance_km.toFixed(1)} km</span> : null}
                  <span>Score: {r.rank_score.toFixed(2)}</span>
                </div>
                {r.warnings?.length ? (
                  <p className="mt-1 text-xs text-amber-700">{r.warnings.join(', ')}</p>
                ) : null}
              </div>
              <Button type="button" onClick={() => void assign(r)} disabled={!r.available}>
                Assign to job
              </Button>
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
}
