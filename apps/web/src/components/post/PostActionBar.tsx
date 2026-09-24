import { CommentOutlined, HeartFilled, HeartOutlined, RetweetOutlined, ShareAltOutlined, StarFilled, StarOutlined } from '@ant-design/icons';
import { Button, Input, Modal, Space, message } from 'antd';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../../api';
import type { Post } from '../../types';

export function PostActionBar({ post, onChange }: { post: Post; onChange?: (patch: Partial<Post>) => void }) {
  const navigate = useNavigate();
  const [shareOpen, setShareOpen] = useState(false);
  const [repostOpen, setRepostOpen] = useState(false);
  const [repostText, setRepostText] = useState('');
  async function like() { const p = await api.likePost(post.postId); onChange?.(p!); }
  async function favorite() { const p = await api.favoritePost(post.postId); onChange?.(p!); }
  async function share() { await navigator.clipboard.writeText(`${location.origin}/posts/${post.postId}`).catch(()=>{}); const p = await api.sharePost(post.postId); onChange?.(p!); message.success('链接已复制'); setShareOpen(false); }
  async function repost() { const p = await api.repostPost(post.postId, repostText); message.success('转发成功'); setRepostOpen(false); setRepostText(''); onChange?.({ repostCount: post.repostCount + 1 }); navigate(`/posts/${p.postId}`); }
  return <><Space wrap className="post-actions"><Button type="text" icon={post.liked ? <HeartFilled /> : <HeartOutlined />} onClick={like}>点赞 {post.likeCount}</Button><Button type="text" icon={<CommentOutlined />} onClick={() => navigate(`/posts/${post.postId}#comments`)}>评论 {post.commentCount}</Button><Button type="text" icon={post.favorited ? <StarFilled /> : <StarOutlined />} onClick={favorite}>收藏 {post.favoriteCount}</Button><Button type="text" icon={<ShareAltOutlined />} onClick={() => setShareOpen(true)}>分享 {post.shareCount}</Button><Button type="text" icon={<RetweetOutlined />} onClick={() => setRepostOpen(true)}>转发 {post.repostCount}</Button></Space><Modal title="分享帖子" open={shareOpen} onCancel={() => setShareOpen(false)} footer={<Button onClick={() => setShareOpen(false)}>关闭</Button>} maskClosable={false}><Input.Group compact><Input value={`${location.origin}/posts/${post.postId}`} readOnly style={{ width: 'calc(100% - 96px)' }} /><Button type="primary" onClick={share}>复制</Button></Input.Group><Button style={{ marginTop: 16 }} block onClick={() => { setShareOpen(false); setRepostOpen(true); }}>转发到我的动态</Button></Modal><Modal title="转发帖子" open={repostOpen} onCancel={() => setRepostOpen(false)} onOk={repost} okText="确认转发" cancelText="取消" maskClosable={false}><div className="source-post"><b>{post.author.nickname}</b><h4>{post.title}</h4><p>{post.summary}</p></div><Input.TextArea rows={4} value={repostText} onChange={e=>setRepostText(e.target.value)} maxLength={300} showCount placeholder="说点什么..." /></Modal></>;
}
