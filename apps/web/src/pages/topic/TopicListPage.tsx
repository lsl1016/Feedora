import { Button, Card, Image, Tabs, Typography } from 'antd';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../../api';
import { EmptyState, ErrorState, LoadingState } from '../../components/common/States';
import { Page } from '../../components/common/Page';
import type { Topic } from '../../types';

type TopicTab = 'all' | 'official' | 'hot' | 'latest';

export function TopicListPage() {
  const nav = useNavigate();
  const [type, setType] = useState<TopicTab>('hot');
  const [list, setList] = useState<Topic[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  async function load() {
    setLoading(true);
    setError('');
    try {
      const result = await api.getTopics({ page: 1, pageSize: 50, type });
      setList(result.list);
    } catch (err) {
      setError(err instanceof Error ? err.message : '话题加载失败');
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => { load(); }, [type]);

  return (
    <Page>
      <Card className="soft-card page-title-card">
        <div>
          <Typography.Title level={2}>话题广场</Typography.Title>
          <Typography.Text type="secondary">发现热门话题，参与内容聚合和讨论。</Typography.Text>
        </div>
      </Card>
      <Card className="soft-card topic-square-card">
        <Tabs activeKey={type} onChange={(key) => setType(key as TopicTab)} items={[{ key: 'all', label: '全部' }, { key: 'official', label: '官方话题' }, { key: 'hot', label: '热门话题' }, { key: 'latest', label: '最新话题' }]} />
        {loading ? <LoadingState /> : error ? <ErrorState message={error} onRetry={load} /> : list.length === 0 ? <EmptyState description="暂无话题" /> : (
          <div className="topic-cover-grid">
            {list.map((topic) => (
              <Card
                key={topic.topicId}
                hoverable
                className="topic-cover-card"
                cover={<Image src={topic.coverImage} height={150} preview={false} fallback="https://images.unsplash.com/photo-1515879218367-8466d910aaa4?auto=format&fit=crop&w=1200&q=80" />}
                onClick={() => nav(`/topics/${topic.topicId}`)}
              >
                <Typography.Title level={4}>#{topic.name}#</Typography.Title>
                <Typography.Paragraph ellipsis={{ rows: 2 }}>{topic.description}</Typography.Paragraph>
                <Typography.Text type="secondary">{topic.participantCount.toLocaleString()} 人参与 · {topic.postCount.toLocaleString()} 篇内容</Typography.Text>
                <Button type="primary" block style={{ marginTop: 14 }} onClick={(event) => { event.stopPropagation(); nav(`/posts/create?topicId=${topic.topicId}`); }}>参与话题</Button>
              </Card>
            ))}
          </div>
        )}
      </Card>
    </Page>
  );
}
