import { Button, Card, Tag, Typography, message } from 'antd';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../../api';
import type { Circle, Post, Topic } from '../../types';

export function HomeRightSidebar() {
  const navigate = useNavigate();
  const [circles, setCircles] = useState<Circle[]>([]);
  const [topics, setTopics] = useState<Topic[]>([]);
  const [keywords, setKeywords] = useState<string[]>([]);
  const [hotPosts, setHotPosts] = useState<Post[]>([]);
  const [checked, setChecked] = useState(false);

  useEffect(() => {
    api.getCircles({ page: 1, pageSize: 3, scope: 'recommended' }).then((res) => setCircles(res.list));
    api.getTopics({ page: 1, pageSize: 4, type: 'hot' }).then((res) => setTopics(res.list));
    api.getHotKeywords().then(setKeywords);
    api.getPosts({ page: 1, pageSize: 3, feedType: 'hot', sort: 'hot' }).then((res) => setHotPosts(res.list));
  }, []);

  async function handleCheckIn() {
    if (checked) {
      message.info('今日已签到');
      return;
    }
    const res = await api.checkIn();
    setChecked(true);
    message.success(`签到成功，获得 ${res.points} 积分`);
  }

  return (
    <aside className="right-aside">
      <Card className="soft-card" title="每日签到">
        <Typography.Paragraph>签到、完成任务、发布内容都能获得积分和经验。</Typography.Paragraph>
        <Button type="primary" block onClick={handleCheckIn}>{checked ? '今日已签到' : '签到领积分'}</Button>
      </Card>
      <Card className="soft-card" title="推荐圈子">
        {circles.map((circle) => (
          <div className="side-item" key={circle.circleId} onClick={() => navigate(`/circles/${circle.circleId}`)}>
            <b>{circle.name}</b><span>{circle.memberCount.toLocaleString()} 成员</span>
          </div>
        ))}
      </Card>
      <Card className="soft-card" title="热门话题">
        {topics.map((topic) => <Tag className="clickable topic-side-tag" color="blue" key={topic.topicId} onClick={() => navigate(`/topics/${topic.topicId}`)}>#{topic.name}#</Tag>)}
      </Card>
      <Card className="soft-card" title="热门搜索">
        {keywords.slice(0, 5).map((keyword) => <div className="side-item compact" key={keyword} onClick={() => navigate(`/search?keyword=${encodeURIComponent(keyword)}`)}>{keyword}</div>)}
      </Card>
      <Card className="soft-card" title={<span className="clickable" onClick={() => navigate('/hot')}>热门榜单</span>}>
        {hotPosts.map((post, index) => (
          <div className="side-rank-item" key={post.postId} onClick={() => navigate(`/posts/${post.postId}`)}>
            <span className={`rank-num rank-${index + 1}`}>{index + 1}</span>
            <div><b>{post.title}</b><p>{post.likeCount} 点赞 · {post.commentCount} 评论 · {post.hotScore} 热度</p></div>
          </div>
        ))}
      </Card>
      <Card className="soft-card" title="项目信息">
        <Typography.Paragraph>开发者知识社区，面向程序员、开发者、IT 从业者和计算机学生，用于技术知识分享、学习交流、求职经验、面试讨论和个人知识沉淀。</Typography.Paragraph>
        <Typography.Paragraph><b>作者：</b>ll</Typography.Paragraph>
        <Typography.Paragraph><b>项目说明：</b>本项目为社区 MVP 前端原型，使用 React + TypeScript + Ant Design 实现。</Typography.Paragraph>
        <Typography.Text type="secondary">© 2026 Developer Knowledge Community. All rights reserved.</Typography.Text>
      </Card>
    </aside>
  );
}
