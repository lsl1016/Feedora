import { Button, Card, Carousel, Space, Typography } from 'antd';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../../api';
import { AiQuickAskCard } from '../../components/home/AiQuickAskCard';
import { HomeFeedFilterBar, defaultHomeFeedFilter, type HomeFeedFilterValue } from '../../components/home/HomeFeedFilterBar';
import { HomeRightSidebar } from '../../components/home/HomeRightSidebar';
import { HomeSideNav } from '../../components/home/HomeSideNav';
import { Page } from '../../components/common/Page';
import { ErrorState } from '../../components/common/States';
import { PostList } from '../../components/post/PostList';
import type { Banner, PageResult, Post } from '../../types';

export function HomePage() {
  const navigate = useNavigate();
  const [data, setData] = useState<PageResult<Post>>();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [filter, setFilter] = useState<HomeFeedFilterValue>(defaultHomeFeedFilter);
  const [banners, setBanners] = useState<Banner[]>([]);

  async function load() {
    setLoading(true);
    setError('');
    try {
      const feedType = filter.feedType === 'friend' ? 'following' : filter.feedType;
      const sort = filter.feedType === 'hot' ? 'hot' : filter.feedType === 'latest' ? 'latest' : 'recommend';
      setData(await api.getPosts({
        page: 1,
        pageSize: 10,
        feedType: feedType as any,
        sort: sort as any,
        timeRange: filter.timeRange,
        tagId: filter.tagId,
        circleId: filter.circleId,
        topicId: filter.topicId,
      }));
    } catch (err) {
      setError(err instanceof Error ? err.message : '加载失败');
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => { load(); }, [filter]);
  useEffect(() => { api.getBanners().then(setBanners); }, []);

  return (
    <Page className="home-three">
      <HomeSideNav />
      <main className="home-main">
        {banners.length > 0 && (
          <Carousel autoplay className="soft-card home-banner">
            {banners.map((banner) => (
              <div key={banner.bannerId}>
                <div
                  className="banner-slide"
                  style={{ backgroundImage: `linear-gradient(90deg,rgba(37,99,235,.85),rgba(79,70,229,.45)),url(${banner.imageUrl})` }}
                  onClick={() => navigate(banner.targetUrl)}
                >
                  <Typography.Title level={3}>{banner.title}</Typography.Title>
                  <Typography.Text>技术学习、知识沉淀、AI 助手和圈子协作</Typography.Text>
                </div>
              </div>
            ))}
          </Carousel>
        )}
        <AiQuickAskCard />
        <Card className="soft-card publish-entry">
          <Typography.Text>分享你的技术经验、项目复盘或学习笔记...</Typography.Text>
          <Space>
            <Button type="primary" onClick={() => navigate('/posts/create')}>发布帖子</Button>
            <Button onClick={() => navigate('/circles/create')}>创建圈子</Button>
          </Space>
        </Card>
        <HomeFeedFilterBar value={filter} onChange={setFilter} />
        {error ? <ErrorState message={error} onRetry={load} /> : <PostList data={data} loading={loading} onChange={(id, patch) => setData((d) => d ? { ...d, list: d.list.map((p) => p.postId === id ? { ...p, ...patch } : p) } : d)} />}
      </main>
      <HomeRightSidebar />
    </Page>
  );
}
