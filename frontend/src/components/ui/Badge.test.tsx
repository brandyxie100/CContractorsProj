import { describe, expect, it } from 'vitest';
import { render, screen } from '@testing-library/react';
import { Badge } from '@/components/ui/Badge';

describe('Badge', () => {
  it('should render available status with label text', () => {
    render(<Badge status="available" />);
    expect(screen.getByText(/available/i)).toBeInTheDocument();
  });

  it('should show custom label for maintenance', () => {
    render(<Badge status="maintenance" label="out of service" />);
    expect(screen.getByText('out of service')).toBeInTheDocument();
  });
});
