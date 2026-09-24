import { Button, Empty, Result, Skeleton } from 'antd';

export function LoadingState({ rows = 2 }: { rows?: number }) {
  return <div className="soft-card state-card">{Array.from({ length: rows }).map((_, i) => <Skeleton key={i} active avatar paragraph={{ rows: 3 }} />)}</div>;
}
export function EmptyState({ description = '暂无数据', actionText, onAction }: { description?: string; actionText?: string; onAction?: () => void }) {
  return <div className="soft-card state-card"><Empty description={description} />{actionText && <Button type="primary" onClick={onAction}>{actionText}</Button>}</div>;
}
export function ErrorState({ message = '加载失败', onRetry }: { message?: string; onRetry?: () => void }) {
  return <Result status="warning" title={message} extra={onRetry && <Button type="primary" onClick={onRetry}>重新加载</Button>} />;
}
