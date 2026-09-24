import type { PostStatus, CircleMemberRole, CircleMemberStatus } from './types';

export function formatTime(value?: string) {
  if (!value) return '';
  return value.slice(0, 16);
}
export function truncate(value: string, len = 90) { return value.length > len ? `${value.slice(0, len)}...` : value; }
export function isAdmin(user?: { role?: string } | null) { return user?.role === 'admin'; }
export function maskApiKey(apiKey: string) { if (!apiKey) return ''; if (apiKey.length <= 8) return '****'; return `${apiKey.slice(0, 3)}****${apiKey.slice(-4)}`; }
export function getErrorMessage(err: unknown) { return err instanceof Error ? err.message : '操作失败，请稍后再试'; }
export const postStatusText: Record<PostStatus, string> = { draft: '草稿', scheduled: '待发布', reviewing: '审核中', published: '已发布', hidden: '已隐藏', rejected: '审核拒绝', deleted: '已删除', takedown: '已下架' };
export const postStatusColor: Record<PostStatus, string> = { draft: 'default', scheduled: 'blue', reviewing: 'orange', published: 'green', hidden: 'purple', rejected: 'red', deleted: 'default', takedown: 'red' };
export const roleText: Record<CircleMemberRole, string> = { owner: '圈主', moderator: '管理员', reviewer: '管理员', member: '普通成员' };
export const statusText: Record<CircleMemberStatus, string> = { normal: '正常', muted: '禁言中', removed: '已移除' };
