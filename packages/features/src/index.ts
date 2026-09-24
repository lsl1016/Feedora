import type { ClientCapability } from '@feedora/app-core';

export interface ClientFeature {
  id: 'community' | 'workspace' | 'repository' | 'search' | 'notification' | 'settings';
  label: string;
  route: string;
  capability?: ClientCapability;
}

export const desktopPrimaryNavigation: readonly ClientFeature[] = [
  { id: 'community', label: '社区', route: '/community', capability: 'community' },
  { id: 'workspace', label: '工作台', route: '/workspace', capability: 'local-workspace' },
  { id: 'repository', label: '仓库', route: '/repository', capability: 'repository' },
  { id: 'search', label: '搜索', route: '/search', capability: 'local-search' },
  { id: 'notification', label: '通知', route: '/notifications', capability: 'community' },
  { id: 'settings', label: '设置', route: '/settings' },
];
