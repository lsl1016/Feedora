import { Avatar, Card, Segmented, Space, Table, Tag, Typography } from 'antd';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../../api';
import { EmptyState, ErrorState } from '../../components/common/States';
import { Page } from '../../components/common/Page';
import type { RankingItem } from '../../types';

type RankingType = 'active' | 'contribution' | 'creator' | 'circle';
type RankingTimeRange = 'day' | 'week' | 'all';

function rankTag(rank: number) {
  if (rank === 1) return <Tag color="gold">Top 1</Tag>;
  if (rank === 2) return <Tag color="default">Top 2</Tag>;
  if (rank === 3) return <Tag color="orange">Top 3</Tag>;
  return rank;
}

export function RankingPage() {
  const navigate = useNavigate();
  const [type, setType] = useState<RankingType>('active');
  const [range, setRange] = useState<RankingTimeRange>('week');
  const [data, setData] = useState<RankingItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  async function load() {
    setLoading(true);
    setError('');
    try {
      setData(await api.getRankings(type, range));
    } catch (err) {
      setError(err instanceof Error ? err.message : '排行榜加载失败');
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => { load(); }, [type, range]);

  const userColumns = [
    { title: '排名', dataIndex: 'rank', width: 90, render: rankTag },
    { title: '用户', render: (_: unknown, r: RankingItem) => <Space className="clickable" onClick={() => navigate(`/users/${r.targetId}`)}><Avatar src={r.avatar} />{r.name}</Space> },
    { title: '等级', render: (_: unknown, r: RankingItem) => r.level ? `Lv${r.level} ${r.levelName}` : '-' },
    { title: '积分', dataIndex: 'points' },
    { title: '发帖数', dataIndex: 'postCount' },
    { title: '获赞数', dataIndex: 'likeReceivedCount' },
    { title: '分数', dataIndex: 'score' },
  ];

  const circleColumns = [
    { title: '排名', dataIndex: 'rank', width: 90, render: rankTag },
    { title: '圈子', render: (_: unknown, r: RankingItem) => <Space className="clickable" onClick={() => navigate(`/circles/${r.targetId}`)}><Avatar src={r.avatar} />{r.name}</Space> },
    { title: '成员数', dataIndex: 'memberCount' },
    { title: '内容数', dataIndex: 'postCount' },
    { title: '精华数', render: () => Math.round(Math.random() * 100) + 20 },
    { title: '分数', dataIndex: 'score' },
  ];

  return (
    <Page>
      <Card className="soft-card ranking-page-card">
        <div className="ranking-page-header">
          <Typography.Title level={2}>排行榜</Typography.Title>
          <Space wrap>
            <Segmented value={type} onChange={(v) => setType(v as RankingType)} options={[{ value: 'active', label: '活跃榜' }, { value: 'contribution', label: '贡献榜' }, { value: 'creator', label: '创作者榜' }, { value: 'circle', label: '圈子榜' }]} />
            <Segmented value={range} onChange={(v) => setRange(v as RankingTimeRange)} options={[{ value: 'day', label: '日榜' }, { value: 'week', label: '周榜' }, { value: 'all', label: '总榜' }]} />
          </Space>
        </div>
        {error ? <ErrorState message={error} onRetry={load} /> : data.length === 0 && !loading ? <EmptyState description="暂无榜单数据" /> : <Table rowKey={(r) => `${r.targetType}-${r.targetId}`} loading={loading} dataSource={data} pagination={false} rowClassName={(r) => r.isCurrentUser ? 'current-user-row' : ''} columns={(type === 'circle' ? circleColumns : userColumns) as any} />}
      </Card>
    </Page>
  );
}
