import { Button, Card, Popconfirm, Space, Typography, message } from 'antd';
import { useEffect, useState } from 'react';
import { api } from '../../api';
import type { CircleSummarySubscription } from '../../types';
import { CircleSummarySubscribeModal } from '../ai/CircleSummarySubscribeModal';

function frequencyText(value?: string) {
  if (value === 'once') return '推送一次';
  if (value === 'daily') return '每日';
  if (value === 'weekly') return '每周';
  if (value === 'monthly') return '每月';
  return '-';
}

export function CircleSummaryManageCard({ circleId }: { circleId: number }) {
  const [subscription, setSubscription] = useState<CircleSummarySubscription>();
  const [open, setOpen] = useState(false);

  async function load() {
    setSubscription(await api.getSummarySubscription(circleId));
  }

  useEffect(() => { load(); }, [circleId]);

  async function cancel() {
    await api.cancelSummarySubscription(circleId);
    message.success('已取消智能摘要订阅');
    load();
  }

  const enabled = subscription?.enabled;
  const pushText = subscription?.frequency === 'once' ? subscription?.oneTimePushAt : subscription?.pushTime;

  return (
    <Card className="soft-card circle-summary-manage-card" title="我的智能摘要">
      {enabled ? (
        <Space direction="vertical" style={{ width: '100%' }}>
          <Typography.Text strong>{subscription?.circleName || '当前圈子'}</Typography.Text>
          <Typography.Text type="secondary">频率：{frequencyText(subscription?.frequency)}</Typography.Text>
          <Typography.Text type="secondary">推送时间：{pushText || '-'}</Typography.Text>
          <Typography.Text type="secondary">推送方式：{subscription?.channels?.map((c) => c === 'email' ? '邮箱' : '站内通知').join('、')}</Typography.Text>
          <Space>
            <Button type="primary" onClick={() => setOpen(true)}>查看配置</Button>
            <Popconfirm title="确认取消订阅？" okText="确认取消" cancelText="再想想" onConfirm={cancel}><Button danger>取消订阅</Button></Popconfirm>
          </Space>
        </Space>
      ) : (
        <Space direction="vertical" style={{ width: '100%' }}>
          <Typography.Text type="secondary">你还没有配置该圈子的智能摘要。</Typography.Text>
          <Button type="primary" onClick={() => setOpen(true)}>开启智能摘要</Button>
        </Space>
      )}
      <CircleSummarySubscribeModal open={open} circleId={circleId} onClose={() => setOpen(false)} onSaved={setSubscription} />
    </Card>
  );
}
