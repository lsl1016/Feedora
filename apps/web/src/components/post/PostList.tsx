import { Pagination } from 'antd';
import type { PageResult, Post } from '../../types';
import { EmptyState, LoadingState } from '../common/States';
import { PostCard } from './PostCard';

export function PostList({ data, loading, mine = false, onChange, onPageChange }: { data?: PageResult<Post>; loading?: boolean; mine?: boolean; onChange?: (id:number, patch:Partial<Post>)=>void; onPageChange?: (page:number,pageSize:number)=>void }) {
  if (loading) return <LoadingState />;
  if (!data || data.list.length === 0) return <EmptyState description="暂无内容" />;
  return <div className="list-stack">{data.list.map(post => <PostCard key={post.postId} post={post} mine={mine} onChange={(patch)=>onChange?.(post.postId,patch)} />)}{onPageChange && <Pagination current={data.page} pageSize={data.pageSize} total={data.total} onChange={onPageChange} />}</div>;
}
