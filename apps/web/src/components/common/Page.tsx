import type { ReactNode } from 'react';

export function Page({ children, narrow = false, className = '' }: { children: ReactNode; narrow?: boolean; className?: string }) {
  return <main className={`page-container ${narrow ? 'page-narrow' : ''} ${className}`}>{children}</main>;
}
