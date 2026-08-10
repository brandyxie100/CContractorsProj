import { useState } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { Button } from '@/components/ui/Button';
import { Card } from '@/components/ui/Card';
import { Input } from '@/components/ui/Input';
import { Modal } from '@/components/ui/Modal';
import { Select } from '@/components/ui/Select';
import { Table } from '@/components/ui/Table';
import { createMaintenance, useAssets, useMaintenance } from '@/features/hooks/useApi';

export function MaintenancePage() {
  const { data, isLoading, error } = useMaintenance();
  const { data: assets } = useAssets();
  const qc = useQueryClient();
  const [open, setOpen] = useState(false);
  const [form, setForm] = useState({
    asset_id: '',
    type: 'preventive',
    description: '',
    scheduled_start: new Date().toISOString().slice(0, 16),
    scheduled_end: new Date(Date.now() + 86400000).toISOString().slice(0, 16),
  });
  const [err, setErr] = useState('');

  async function onCreate(e: React.FormEvent) {
    e.preventDefault();
    setErr('');
    try {
      await createMaintenance({
        ...form,
        scheduled_start: new Date(form.scheduled_start).toISOString(),
        scheduled_end: new Date(form.scheduled_end).toISOString(),
        status: 'scheduled',
      });
      setOpen(false);
      await qc.invalidateQueries({ queryKey: ['maintenance'] });
      await qc.invalidateQueries({ queryKey: ['assets'] });
    } catch (ex) {
      setErr(ex instanceof Error ? ex.message : 'Create failed');
    }
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-[var(--color-navy)]">Maintenance</h1>
        <Button type="button" onClick={() => setOpen(true)}>
          Schedule
        </Button>
      </div>
      <Card>
        {isLoading ? <p>Loading…</p> : null}
        {error ? <p className="text-red-600">{(error as Error).message}</p> : null}
        <Table headers={['Asset', 'Type', 'Status', 'Start', 'End', 'Description']}>
          {(data ?? []).map((m) => (
            <tr key={m.id}>
              <td className="px-3 py-2">{m.asset?.asset_code ?? m.asset_id.slice(0, 8)}</td>
              <td className="px-3 py-2">{m.type}</td>
              <td className="px-3 py-2">{m.status}</td>
              <td className="px-3 py-2">{new Date(m.scheduled_start).toLocaleString()}</td>
              <td className="px-3 py-2">{new Date(m.scheduled_end).toLocaleString()}</td>
              <td className="px-3 py-2">{m.description}</td>
            </tr>
          ))}
        </Table>
      </Card>
      <Modal open={open} title="Schedule maintenance" onClose={() => setOpen(false)}>
        <form className="space-y-3" onSubmit={(e) => void onCreate(e)}>
          <Select
            label="Asset"
            value={form.asset_id}
            onChange={(e) => setForm({ ...form, asset_id: e.target.value })}
            options={[
              { value: '', label: 'Select…' },
              ...(assets ?? []).map((a) => ({ value: a.id, label: `${a.asset_code} — ${a.name}` })),
            ]}
          />
          <Select
            label="Type"
            value={form.type}
            onChange={(e) => setForm({ ...form, type: e.target.value })}
            options={[
              { value: 'preventive', label: 'Preventive' },
              { value: 'repair', label: 'Repair' },
              { value: 'inspection', label: 'Inspection' },
            ]}
          />
          <Input label="Start" type="datetime-local" value={form.scheduled_start} onChange={(e) => setForm({ ...form, scheduled_start: e.target.value })} />
          <Input label="End" type="datetime-local" value={form.scheduled_end} onChange={(e) => setForm({ ...form, scheduled_end: e.target.value })} />
          <Input label="Description" value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} />
          {err ? <p className="text-sm text-red-600">{err}</p> : null}
          <Button type="submit">Save</Button>
        </form>
      </Modal>
    </div>
  );
}
