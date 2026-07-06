import { Avatar, Button, Card, Space, Typography } from 'antd';
import type { User } from '../../types';

export function MyProfileHeroCard({ user, onEdit }: { user: User; onEdit: () => void }) {
  return (
    <Card className="soft-card my-profile-hero-card" bodyStyle={{ padding: 0 }}>
      <div className="my-profile-hero-bg" />
      <div className="my-profile-hero-content">
        <Space align="start" className="my-profile-hero-space">
          <Avatar src={user.avatar} size={82} className="my-profile-avatar" />
          <div className="my-profile-info">
            <Typography.Title level={3}>{user.nickname}</Typography.Title>
            <Typography.Paragraph>{user.bio}</Typography.Paragraph>
            <Space wrap className="profile-stats-line">
              <span>关注 {user.followingCount}</span>
              <span>粉丝 {user.followerCount}</span>
              <span>获赞 {user.likeReceivedCount.toLocaleString()}</span>
            </Space>
          </div>
          <Button onClick={onEdit}>编辑资料</Button>
        </Space>
      </div>
    </Card>
  );
}
