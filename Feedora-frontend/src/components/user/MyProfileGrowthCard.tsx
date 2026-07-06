import { Button, Card, Progress, Space, Typography } from 'antd';
import { useNavigate } from 'react-router-dom';
import type { User } from '../../types';

export function MyProfileGrowthCard({ user, onCheckIn }: { user: User; onCheckIn: () => void }) {
  const navigate = useNavigate();
  const percent = Math.round((user.experience / user.nextLevelExperience) * 100);
  return (
    <Card className="soft-card my-profile-growth-card" title="成长信息">
      <Space direction="vertical" style={{ width: '100%' }}>
        <Space className="toolbar-space" wrap>
          <Typography.Text strong>Lv{user.level} {user.levelName} · 积分 {user.points}</Typography.Text>
          <Space>
            <Button type="primary" onClick={onCheckIn}>{user.checkedInToday ? '已签到' : '签到'}</Button>
            <Button onClick={() => navigate('/tasks')}>任务中心</Button>
            <Button onClick={() => navigate('/creator')}>创作者中心</Button>
          </Space>
        </Space>
        <Progress percent={percent} status="active" />
        <Typography.Text type="secondary">经验值：{user.experience} / {user.nextLevelExperience} · 连续签到 {user.continuousCheckInDays} 天</Typography.Text>
      </Space>
    </Card>
  );
}
