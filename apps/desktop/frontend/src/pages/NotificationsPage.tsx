import { BellOutlined, CheckOutlined, CommentOutlined, HeartOutlined, InfoCircleOutlined } from '@ant-design/icons';
import { Button, Card, Segmented, Space, Typography } from 'antd';
import { useMemo, useState } from 'react';
import { useDesktopStore } from '../store/desktopStore';

function icon(category:string) {
  if (category === 'interaction') return <CommentOutlined />;
  if (category === 'circle') return <HeartOutlined />;
  if (category === 'system') return <InfoCircleOutlined />;
  return <BellOutlined />;
}

export function NotificationsPage() {
  const { notifications, markNotificationRead, markAllNotificationsRead, askAi } = useDesktopStore();
  const [filter, setFilter] = useState('all');
  const rows = useMemo(() => notifications.filter((item:any) => filter === 'all' || item.category === filter), [notifications, filter]);

  return (
    <div className="desktop-page">
      <div className="page-heading-row">
        <div><Typography.Title level={2}>通知中心</Typography.Title><Typography.Text type="secondary">社区互动、发布反馈与桌面系统消息。</Typography.Text></div>
        <Space><Button icon={<CheckOutlined />} onClick={markAllNotificationsRead}>全部标记已读</Button><Button onClick={() => askAi('请总结今天的重要通知，并指出哪些需要我处理。')}>AI 总结</Button></Space>
      </div>

      <Segmented value={filter} onChange={(value)=>setFilter(String(value))} options={[
        {value:'all',label:'全部'},
        {value:'interaction',label:'评论与互动'},
        {value:'circle',label:'圈子动态'},
        {value:'system',label:'系统消息'},
      ]} />

      <div className="notification-list">
        {rows.map((item:any) => (
          <Card key={item.notificationId} className={`notification-card ${item.readStatus === 'unread' ? 'unread' : ''}`} onClick={() => markNotificationRead(item.notificationId)}>
            <Space align="start" size={16}>
              <div className="notification-icon">{icon(item.category)}</div>
              <div>
                <Space><b>{item.title}</b>{item.readStatus === 'unread' && <span className="unread-dot" />}</Space>
                <Typography.Paragraph type="secondary">{item.content}</Typography.Paragraph>
                <span className="muted">{item.createdAt}</span>
              </div>
            </Space>
          </Card>
        ))}
      </div>
    </div>
  );
}
