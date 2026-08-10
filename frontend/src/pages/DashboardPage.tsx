import { Bar, BarChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts';
import { Card } from '@/components/ui/Card';
import { useDashboard } from '@/features/hooks/useApi';

export function DashboardPage() {
  const { data, isLoading, error } = useDashboard();

  if (isLoading) return <p>Loading dashboard…</p>;
  if (error) return <p className="text-red-600">{(error as Error).message}</p>;
  if (!data) return null;

  const chartData = Object.entries(data.fleet_by_status ?? {}).map(([status, count]) => ({
    status,
    count,
  }));

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold text-[var(--color-navy)]">Dashboard</h1>
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <Card>
          <div className="text-xs uppercase text-slate-500">Utilization</div>
          <div className="text-3xl font-semibold">{data.utilization_percent.toFixed(0)}%</div>
        </Card>
        <Card>
          <div className="text-xs uppercase text-slate-500">Active assignments</div>
          <div className="text-3xl font-semibold">{data.active_assignments}</div>
        </Card>
        <Card>
          <div className="text-xs uppercase text-slate-500">Upcoming (7d)</div>
          <div className="text-3xl font-semibold">{data.upcoming_assignments}</div>
        </Card>
        <Card>
          <div className="text-xs uppercase text-slate-500">Overdue maintenance</div>
          <div className="text-3xl font-semibold text-orange-700">{data.overdue_maintenance}</div>
        </Card>
      </div>
      <Card>
        <h2 className="mb-3 font-semibold text-[var(--color-steel)]">Fleet by status</h2>
        <div className="h-64">
          <ResponsiveContainer width="100%" height="100%">
            <BarChart data={chartData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="status" />
              <YAxis allowDecimals={false} />
              <Tooltip />
              <Bar dataKey="count" fill="#2c5f7c" />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </Card>
    </div>
  );
}
