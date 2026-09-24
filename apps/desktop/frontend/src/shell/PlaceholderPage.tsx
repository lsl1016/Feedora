import { Card, Typography } from 'antd';

export function PlaceholderPage(props: { title: string; description: string }) {
  return (
    <div className="desktop-page">
      <Typography.Title level={2}>{props.title}</Typography.Title>
      <Typography.Paragraph type="secondary">{props.description}</Typography.Paragraph>
      <Card className="desktop-placeholder-card">
        该页面已进入 Desktop Shell 路由，下一阶段按冻结的高保真页面逐页实现。
      </Card>
    </div>
  );
}
