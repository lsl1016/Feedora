import { Button, Card, Popconfirm, Space, Tag, Typography, message } from 'antd';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../../api';
import type { MyCommentItem } from '../../types';
import { EmptyState, LoadingState } from '../common/States';

export function MyCommentList() {
  const navigate = useNavigate(); const [list, setList] = useState<MyCommentItem[]>([]); const [loading, setLoading] = useState(true);
  async function load(){ setLoading(true); const r=await api.getMyComments({page:1,pageSize:50}); setList(r.list); setLoading(false); }
  useEffect(()=>{load();},[]);
  async function remove(id:number){ await api.deleteComment(id); message.success('删除成功'); load(); }
  if(loading)return <LoadingState />; if(list.length===0)return <EmptyState description="你还没有发表过评论" />;
  return <div className="list-stack">{list.map(c=><Card key={c.commentId} className="soft-card"><Typography.Paragraph>{c.content}</Typography.Paragraph><Typography.Text type="secondary">评论对象：《{c.postTitle}》</Typography.Text><br/><Space wrap style={{marginTop:12}}><Tag color={c.status==='normal'?'green':c.status==='deleted'?'default':'red'}>{c.status==='normal'?'正常':c.status==='deleted'?'已删除':'审核拒绝'}</Tag><span>点赞 {c.likeCount}</span><span>{c.createdAt}</span><Button size="small" onClick={()=>navigate(`/posts/${c.postId}#comment-${c.commentId}`)}>查看原帖</Button><Popconfirm title="确认删除评论？" onConfirm={()=>remove(c.commentId)}><Button size="small" danger>删除评论</Button></Popconfirm></Space></Card>)}</div>;
}
