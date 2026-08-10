import { useState } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { Button } from '@/components/ui/Button';
import { Card } from '@/components/ui/Card';
import { Input } from '@/components/ui/Input';
import { Modal } from '@/components/ui/Modal';
import { Table } from '@/components/ui/Table';
import { useAuth } from '@/features/auth/AuthProvider';
import { createLocation, useLocations, useUsers } from '@/features/hooks/useApi';

export function SettingsPage() {
  const { user } = useAuth();
  const { data: locations } = useLocations();
  const usersQuery = useUsers();
  const qc = useQueryClient();
  const [open, setOpen] = useState(false);
  const [form, setForm] = useState({ name: '', type: 'depot', address: '', lat: '', lng: '' });
  const [err, setErr] = useState('');

  async function onCreate(e: React.FormEvent) {
    e.preventDefault();
    setErr('');
    try {
      await createLocation({
        name: form.name,
        type: form.type,
        address: form.address,
        lat: form.lat ? Number(form.lat) : null,
        lng: form.lng ? Number(form.lng) : null,
      });
      setOpen(false);
      await qc.invalidateQueries({ queryKey: ['locations'] });
    } catch (ex) {
      setErr(ex instanceof Error ? ex.message : 'Create failed');
    }
  }

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold text-[var(--color-navy)]">Settings</h1>
      <Card>
        <h2 className="mb-2 font-semibold">Locations</h2>
        <div className="mb-3">
          <Button type="button" onClick={() => setOpen(true)}>
            Add location
          </Button>
        </div>
        <Table headers={['Name', 'Type', 'Address']}>
          {(locations ?? []).map((l) => (
            <tr key={l.id}>
              <td className="px-3 py-2">{l.name}</td>
              <td className="px-3 py-2">{l.type}</td>
              <td className="px-3 py-2">{l.address}</td>
            </tr>
          ))}
        </Table>
      </Card>

      {user?.role === 'admin' ? (
        <Card>
          <h2 className="mb-2 font-semibold">Users</h2>
          {usersQuery.error ? (
            <p className="text-sm text-slate-600">Unable to load users (admin only).</p>
          ) : (
            <Table headers={['Name', 'Email', 'Role']}>
              {(usersQuery.data ?? []).map((u) => (
                <tr key={u.id}>
                  <td className="px-3 py-2">{u.name}</td>
                  <td className="px-3 py-2">{u.email}</td>
                  <td className="px-3 py-2">{u.role}</td>
                </tr>
              ))}
            </Table>
          )}
        </Card>
      ) : null}

      <Modal open={open} title="New location" onClose={() => setOpen(false)}>
        <form className="space-y-3" onSubmit={(e) => void onCreate(e)}>
          <Input label="Name" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required />
          <Input label="Type" value={form.type} onChange={(e) => setForm({ ...form, type: e.target.value })} />
          <Input label="Address" value={form.address} onChange={(e) => setForm({ ...form, address: e.target.value })} />
          <Input label="Lat" value={form.lat} onChange={(e) => setForm({ ...form, lat: e.target.value })} />
          <Input label="Lng" value={form.lng} onChange={(e) => setForm({ ...form, lng: e.target.value })} />
          {err ? <p className="text-sm text-red-600">{err}</p> : null}
          <Button type="submit">Create</Button>
        </form>
      </Modal>
    </div>
  );
}
