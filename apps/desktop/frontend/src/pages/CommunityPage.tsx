import {
  BookOutlined,
  CommentOutlined,
  HeartFilled,
  HeartOutlined,
  LikeFilled,
  LikeOutlined,
  RobotOutlined,
} from '@ant-design/icons';
import { Avatar, Button, Card, Segmented, Space, Tag, Typography, message } from 'antd';
import { useMemo, useState } from 'react';
import { useDesktopStore } from '../store/desktopStore';

export function CommunityPage() {
  const { posts, toggleLike, toggleFavorite, savePostAsNote, addContext, askAi } = useDesktopStore();
  const [tab, setTab] = useState<'recommend' | 'latest' | 'hot'>('recommend');

  const rows = useMemo(() => {
    if (tab === 'hot') return [...posts].sort((a, b) => b.hotScore - a.hotScore);
    if (tab === 'latest') return [...posts].sort((a, b) => String(b.createdAt).localeCompare(String(a.createdAt)));
    return posts;
  }, [posts, tab]);

  return (
    <div className="desktop-page desktop-page-community">
      <div className="page-heading-row">
        <div>
          <Typography.Title level={2}>社区首页</Typography.Title>
          <Typography.Text type="secondary">发现高质量开发者内容，并把有价值的信息沉淀到本地 Workspace。</Typography.Text>
        </div>
        <Button type="primary">发布帖子</Button>
      </div>

      <Card className="feed-toolbar">
        <Segmented value={tab} onChange={(value) => setTab(value as any)} options={[
          { value: 'recommend', label: '推荐' },
          { value: 'latest', label: '最新' },
          { value: 'hot', label: '热门' },
        ]} />
        <Space><Tag>Go</Tag><Tag>AI Agent</Tag><Tag>桌面应用</Tag></Space>
      </Card>

      <div className="feed-list">
        {rows.map((post, index) => (
          <Card key={post.postId} className={`feed-card ${index === 0 ? 'selected-card' : ''}`}>
            <div className="feed-card-meta">
              <Space>
                <Avatar src={post.author.avatar} />
                <div><b>{post.author.nickname}</b><div className="muted">{post.createdAt}</div></div>
              </Space>
              <Space>{post.circle && <Tag>{post.circle.name}</Tag>}{post.isFeatured && <Tag color="blue">精选</Tag>}</Space>
            </div>
            <Typography.Title level={4}>{post.title}</Typography.Title>
            <Typography.Paragraph type="secondary">{post.summary}</Typography.Paragraph>
            <Space wrap>{post.tags.map((tag:any) => <Tag key={tag.tagId}>{tag.tagName}</Tag>)}</Space>
            <div className="feed-actions">
              <Space>
                <Button type="text" icon={post.liked ? <LikeFilled /> : <LikeOutlined />} onClick={() => toggleLike(post.postId)}>{post.likeCount}</Button>
                <Button type="text" icon={<CommentOutlined />}>{post.commentCount}</Button>
                <Button type="text" icon={post.favorited ? <HeartFilled /> : <HeartOutlined />} onClick={() => toggleFavorite(post.postId)}>{post.favoriteCount}</Button>
              </Space>
              <Space>
                <Button icon={<BookOutlined />} onClick={() => {
                  savePostAsNote(post.postId);
                  message.success('已保存到本地笔记');
                }}>保存为笔记</Button>
                <Button onClick={() => {
                  addContext({ id: 'post-' + post.postId, type: 'post', label: '帖子: ' + post.title });
                  message.success('已加入 AI Context');
                }}>添加到 Context</Button>
                <Button type="primary" icon={<RobotOutlined />} onClick={() => {
                  addContext({ id: 'post-' + post.postId, type: 'post', label: '帖子: ' + post.title });
                  askAi('请总结这篇帖子，并给出可以沉淀为笔记的结构。');
                }}>让 AI 分析</Button>
              </Space>
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
}
