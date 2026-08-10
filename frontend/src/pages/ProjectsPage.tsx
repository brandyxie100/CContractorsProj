import { useState } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { Button } from '@/components/ui/Button';
import { Card } from '@/components/ui/Card';
import { Input } from '@/components/ui/Input';
import { Modal } from '@/components/ui/Modal';
import { Table } from '@/components/ui/Table';
import { createProject, useProjects } from '@/features/hooks/useApi';

export function ProjectsPage() {
  const { data, isLoading, error } = useProjects();
  const qc = useQueryClient();
  const [open, setOpen] = useState(false);
  const [form, setForm] = useState({
    name: '',
    address: '',
    start_date: new Date().toISOString().slice(0, 10),
    end_date: new Date(Date.now() + 14 * 86400000).toISOString().slice(0, 10),
    description: '',
  });
  const [err, setErr] = useState('');

  async function onCreate(e: React.FormEvent) {
    e.preventDefault();
    setErr('');
    try {
      await createProject({ ...form, status: 'active' });
      setOpen(false);
      await qc.invalidateQueries({ queryKey: ['projects'] });
    } catch (ex) {
      setErr(ex instanceof Error ? ex.message : 'Create failed');
    }
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-[var(--color-navy)]">Projects</h1>
        <Button type="button" onClick={() => setOpen(true)}>
          New project
        </Button>
      </div>
      <Card>
        {isLoading ? <p>Loading…</p> : null}
        {error ? <p className="text-red-600">{(error as Error).message}</p> : null}
        <Table headers={['Name', 'Status', 'Start', 'End', 'Address']}>
          {(data ?? []).map((p) => (
            <tr key={p.id}>
              <td className="px-3 py-2 font-medium">{p.name}</td>
              <td className="px-3 py-2">{p.status}</td>
              <td className="px-3 py-2">{String(p.start_date).slice(0, 10)}</td>
              <td className="px-3 py-2">{String(p.end_date).slice(0, 10)}</td>
              <td className="px-3 py-2">{p.address}</td>
            </tr>
          ))}
        </Table>
      </Card>
      <Modal open={open} title="New project" onClose={() => setOpen(false)}>
        <form className="space-y-3" onSubmit={(e) => void onCreate(e)}>
          <Input label="Name" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required />
          <Input label="Address" value={form.address} onChange={(e) => setForm({ ...form, address: e.target.value })} />
          <Input label="Start" type="date" value={form.start_date} onChange={(e) => setForm({ ...form, start_date: e.target.value })} />
          <Input label="End" type="date" value={form.end_date} onChange={(e) => setForm({ ...form, end_date: e.target.value })} />
          <Input label="Description" value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} />
          {err ? <p className="text-sm text-red-600">{err}</p> : null}
          <Button type="submit">Create</Button>
        </form>
      </Modal>
    </div>
  );
}
