import { NavLink, Outlet } from 'react-router-dom';
import { useAuth } from '@/features/auth/AuthProvider';
import { Button } from '@/components/ui/Button';

const links = [
  { to: '/', label: 'Dashboard' },
  { to: '/calendar', label: 'Calendar' },
  { to: '/dispatch', label: 'Dispatch' },
  { to: '/assets', label: 'Assets' },
  { to: '/projects', label: 'Projects' },
  { to: '/maintenance', label: 'Maintenance' },
  { to: '/reports', label: 'Reports' },
  { to: '/settings', label: 'Settings' },
];

export function AppShell() {
  const { user, logout } = useAuth();

  return (
    <div className="min-h-screen md:flex">
      <aside className="w-full border-b border-slate-200 bg-[var(--color-navy)] text-white md:w-56 md:border-b-0 md:border-r">
        <div className="px-4 py-5">
          <div className="text-xs uppercase tracking-widest text-slate-300">Clements</div>
          <div className="text-lg font-semibold">Fleet Availability</div>
        </div>
        <nav className="flex gap-1 overflow-x-auto px-2 pb-3 md:flex-col md:overflow-visible">
          {links.map((l) => (
            <NavLink
              key={l.to}
              to={l.to}
              end={l.to === '/'}
              className={({ isActive }) =>
                `whitespace-nowrap rounded px-3 py-2 text-sm ${
                  isActive ? 'bg-[var(--color-steel)]' : 'hover:bg-white/10'
                }`
              }
            >
              {l.label}
            </NavLink>
          ))}
        </nav>
      </aside>
      <div className="flex min-w-0 flex-1 flex-col">
        <header className="flex items-center justify-between border-b border-slate-200 bg-white px-4 py-3">
          <div className="text-sm text-slate-600">Equipment Availability System</div>
          <div className="flex items-center gap-3 text-sm">
            <span>
              {user?.name} ({user?.role})
            </span>
            <Button variant="secondary" type="button" onClick={() => void logout()}>
              Log out
            </Button>
          </div>
        </header>
        <main className="flex-1 p-4 md:p-6">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
