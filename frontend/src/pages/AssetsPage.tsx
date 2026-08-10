import { useState } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { Badge } from '@/components/ui/Badge';
import { Button } from '@/components/ui/Button';
import { Card } from '@/components/ui/Card';
import { Input } from '@/components/ui/Input';
import { Modal } from '@/components/ui/Modal';
import { Select } from '@/components/ui/Select';
import { Table } from '@/components/ui/Table';
import { createAsset, useAssets, useLocations } from '@/features/hooks/useApi';

export function AssetsPage() {
  const { data, isLoading, error } = useAssets();
  const { data: locations } = useLocations();
  const qc = useQueryClient();
  const [open, setOpen] = useState(false);
  const [form, setForm] = useState({
    asset_code: '',
    name: '',
    model: '',
    type: 'excavator',
    category: 'heavy_equipment',
    current_location_id: '',
    weight_t: '20',
  });
  const [err, setErr] = useState('');

  async function onCreate(e: React.FormEvent) {
    e.preventDefault();
    setErr('');
    try {
      await createAsset({
        asset_code: form.asset_code,
        name: form.name,
        model: form.model,
        type: form.type,
        category: form.category,
        current_location_id: form.current_location_id || null,
        specs: { weight_t: Number(form.weight_t) },
        status: 'available',
      });
      setOpen(false);
      await qc.invalidateQueries({ queryKey: ['assets'] });
    } catch (ex) {
      setErr(ex instanceof Error ? ex.message : 'Create failed');
    }
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-[var(--color-navy)]">Assets</h1>
        <Button type="button" onClick={() => setOpen(true)}>
          Add asset
        </Button>
      </div>
      <Card>
        {isLoading ? <p>Loading…</p> : null}
        {error ? <p className="text-red-600">{(error as Error).message}</p> : null}
        <Table headers={['Code', 'Name', 'Type', 'Status', 'Location']}>
          {(data ?? []).map((a) => (
            <tr key={a.id}>
              <td className="px-3 py-2 font-medium">{a.asset_code}</td>
              <td className="px-3 py-2">{a.name}</td>
              <td className="px-3 py-2">{a.type}</td>
              <td className="px-3 py-2">
                <Badge status={a.status} />
              </td>
              <td className="px-3 py-2">{a.current_location?.name ?? '—'}</td>
            </tr>
          ))}
        </Table>
      </Card>

      <Modal open={open} title="New asset" onClose={() => setOpen(false)}>
        <form className="space-y-3" onSubmit={(e) => void onCreate(e)}>
          <Input label="Asset code" value={form.asset_code} onChange={(e) => setForm({ ...form, asset_code: e.target.value })} required />
          <Input label="Name" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required />
          <Input label="Model" value={form.model} onChange={(e) => setForm({ ...form, model: e.target.value })} />
          <Select
            label="Type"
            value={form.type}
            onChange={(e) => setForm({ ...form, type: e.target.value })}
            options={[
              { value: 'excavator', label: 'Excavator' },
              { value: 'truck', label: 'Truck' },
              { value: 'roller', label: 'Roller' },
            ]}
          />
          <Select
            label="Location"
            value={form.current_location_id}
            onChange={(e) => setForm({ ...form, current_location_id: e.target.value })}
            options={[
              { value: '', label: 'None' },
              ...(locations ?? []).map((l) => ({ value: l.id, label: l.name })),
            ]}
          />
          <Input label="Weight (t)" value={form.weight_t} onChange={(e) => setForm({ ...form, weight_t: e.target.value })} />
          {err ? <p className="text-sm text-red-600">{err}</p> : null}
          <Button type="submit">Create</Button>
        </form>
      </Modal>
    </div>
  );
}
