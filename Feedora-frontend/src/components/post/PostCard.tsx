import { Avatar, Card, Image, Space, Tag, Typography } from 'antd';
import { useNavigate } from 'react-router-dom';
import type { Post } from '../../types';
import { postStatusColor, postStatusText } from '../../utils';
import { PostActionBar } from './PostActionBar';

export function PostCard({ post, onChange, mine = false }: { post: Post; onChange?: (patch: Partial<Post>) => void; mine?: boolean }) {
  const navigate = useNavigate();
  return <Card className="soft-card post-card" hoverable><Space className="post-author" onClick={()=>navigate(`/users/${post.authorId}`)}><Avatar src={post.author.avatar} /><div><b>{post.author.nickname}</b><span>Lv{post.author.level} {post.author.levelName}</span></div></Space>{post.circle && <Tag color="blue" onClick={()=>navigate(`/circles/${post.circleId}`)}>来自圈子：{post.circle.name}</Tag>}<div className="post-click-zone" onClick={()=>navigate(`/posts/${post.postId}`)}><Typography.Title level={4}>{post.isTop && <Tag color="red">置顶</Tag>}{post.isFeatured && <Tag color="gold">精华</Tag>}{mine && <Tag color={postStatusColor[post.status]}>{postStatusText[post.status]}</Tag>}{post.title}</Typography.Title>{post.topics.map(t=><Tag color="geekblue" key={t.topicId}>#{t.name}#</Tag>)}{post.postType === 'repost' && post.sourcePost ? <div className="source-post"><p>{post.repostComment}</p><b>原帖：{post.sourcePost.title}</b><p>{post.sourcePost.summary}</p></div> : <Typography.Paragraph ellipsis={{ rows: 3 }}>{post.summary}</Typography.Paragraph>}<Space wrap>{post.tags.map(t=><Tag key={t.tagId} onClick={(e)=>{e.stopPropagation();navigate(`/tags/${t.tagId}`)}}>{t.tagName}</Tag>)}</Space>{post.images.length>0 && <div className="post-image-grid">{post.images.slice(0,3).map(src=><Image key={src} src={src} preview={false} />)}</div>}</div><PostActionBar post={post} onChange={onChange} /></Card>;
}
