import { Card, Segmented, Table, Typography } from 'antd';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../../api';
import { EmptyState, ErrorState } from '../../components/common/States';
import { Page } from '../../components/common/Page';
import { HotPageFooterInfo } from '../../components/home/HotPageFooterInfo';

type RankType = 'post' | 'circle' | 'topic';
type TimeRange = 'today' | 'week' | 'all';

export function HotPage() {
  const navigate = useNavigate();
  const [rankType, setRankType] = useState<RankType>('post');
  const [timeRange, setTimeRange] = useState<TimeRange>('today');
  const [list, setList] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  async function load() {
    setLoading(true);
    setError('');
    try {
      setList(await api.getHotRanks(rankType, timeRange));
    } catch (err) {
      setError(err instanceof Error ? err.message : '榜单加载失败');
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => { load(); }, [rankType, timeRange]);

  const postColumns = [
    { title: '排名', dataIndex: 'rank', width: 80, render: (rank: number) => <span className={`rank-num rank-${rank}`}>{rank}</span> },
    { title: '标题', dataIndex: 'title', render: (title: string, record: any) => <Typography.Link onClick={() => navigate(`/posts/${record.postId}`)}>{title}</Typography.Link> },
    { title: '作者', dataIndex: 'authorName' },
    { title: '热度', dataIndex: 'hotScore' },
    { title: '点赞', dataIndex: 'likeCount' },
    { title: '评论', dataIndex: 'commentCount' },
  ];
  const circleColumns = [
    { title: '排名', dataIndex: 'rank', width: 80, render: (rank: number) => <span className={`rank-num rank-${rank}`}>{rank}</span> },
    { title: '圈子', dataIndex: 'name', render: (name: string, record: any) => <Typography.Link onClick={() => navigate(`/circles/${record.circleId}`)}>{name}</Typography.Link> },
    { title: '成员数', dataIndex: 'memberCount' },
    { title: '内容数', dataIndex: 'postCount' },
    { title: '精华数', dataIndex: 'featuredPostCount' },
    { title: '热度', dataIndex: 'hotScore' },
  ];
  const topicColumns = [
    { title: '排名', dataIndex: 'rank', width: 80, render: (rank: number) => <span className={`rank-num rank-${rank}`}>{rank}</span> },
    { title: '话题', dataIndex: 'name', render: (name: string, record: any) => <Typography.Link onClick={() => navigate(`/topics/${record.topicId}`)}>#{name}#</Typography.Link> },
    { title: '参与人数', dataIndex: 'participantCount' },
    { title: '内容数', dataIndex: 'postCount' },
    { title: '热度', dataIndex: 'hotScore' },
  ];

  const columns = rankType === 'post' ? postColumns : rankType === 'circle' ? circleColumns : topicColumns;

  return (
    <Page narrow>
      <Card className="soft-card hot-rank-page-card">
        <div className="hot-page-header">
          <Typography.Title level={2}>热门榜单</Typography.Title>
          <Segmented value={rankType} onChange={(v) => setRankType(v as RankType)} options={[{ value: 'post', label: '热门帖子' }, { value: 'circle', label: '热门圈子' }, { value: 'topic', label: '热门话题' }]} />
        </div>
        <div className="hot-time-tabs">
          <Segmented value={timeRange} onChange={(v) => setTimeRange(v as TimeRange)} options={[{ value: 'today', label: '今日热门' }, { value: 'week', label: '本周热门' }, { value: 'all', label: '总榜' }]} />
        </div>
        {error ? <ErrorState message={error} onRetry={load} /> : list.length === 0 && !loading ? <EmptyState description="暂无热门内容" /> : <Table rowKey={(record) => `${rankType}-${record.postId || record.circleId || record.topicId}`} loading={loading} dataSource={list} pagination={false} columns={columns as any} />}
      </Card>
      <HotPageFooterInfo />
    </Page>
  );
}
