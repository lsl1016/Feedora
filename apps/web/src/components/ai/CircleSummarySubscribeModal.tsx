import { Checkbox, DatePicker, Form, Input, Modal, Radio, Switch, TimePicker, Typography, message } from 'antd';
import dayjs from 'dayjs';
import { useEffect } from 'react';
import { api } from '../../api';
import type { CircleSummarySubscription } from '../../types';

export function CircleSummarySubscribeModal({ open, circleId, onClose, onSaved }: { open: boolean; circleId: number; onClose: () => void; onSaved?: (subscription: CircleSummarySubscription) => void }) {
  const [form] = Form.useForm();
  const enabled = Form.useWatch('enabled', form);
  const channels = Form.useWatch('channels', form) || [];
  const frequency = Form.useWatch('frequency', form);

  useEffect(() => {
    if (!open) return;
    api.getSummarySubscription(circleId).then((subscription) => {
      form.setFieldsValue({
        ...subscription,
        pushTime: subscription.pushTime ? dayjs(subscription.pushTime, 'HH:mm') : dayjs('09:00', 'HH:mm'),
        oneTimePushAt: subscription.oneTimePushAt ? dayjs(subscription.oneTimePushAt) : undefined,
      });
    });
  }, [open, circleId, form]);

  async function save() {
    const values = await form.validateFields();
    if (values.enabled && (!values.channels || values.channels.length === 0)) {
      message.warning('请至少选择一种推送方式');
      return;
    }
    if (values.enabled && (!values.contentScopes || values.contentScopes.length === 0)) {
      message.warning('请至少选择一种摘要内容');
      return;
    }
    if (values.enabled && values.frequency === 'once' && values.oneTimePushAt?.isBefore(dayjs())) {
      message.warning('推送时间必须晚于当前时间');
      return;
    }
    const payload = {
      ...values,
      pushTime: values.frequency === 'once' ? undefined : values.pushTime?.format('HH:mm'),
      oneTimePushAt: values.frequency === 'once' ? values.oneTimePushAt?.format('YYYY-MM-DD HH:mm:ss') : undefined,
    };
    const result = await api.saveSummarySubscription(circleId, payload);
    message.success('智能摘要配置已保存');
    onSaved?.(result);
    onClose();
  }

  return (
    <Modal title="圈子智能摘要" open={open} onCancel={onClose} onOk={save} okText="保存配置" maskClosable={false}>
      <Form form={form} layout="vertical" initialValues={{ enabled: true, frequency: 'daily', channels: ['notification'], contentScopes: ['post', 'chat'], pushTime: dayjs('09:00', 'HH:mm') }}>
        <Form.Item name="enabled" label="是否开启" valuePropName="checked"><Switch /></Form.Item>
        <Form.Item name="frequency" label="摘要频率" rules={[{ required: enabled, message: '请选择摘要频率' }]}>
          <Radio.Group disabled={!enabled} options={[{ value: 'once', label: '推送一次' }, { value: 'daily', label: '每日' }, { value: 'weekly', label: '每周' }, { value: 'monthly', label: '每月' }]} />
        </Form.Item>
        {frequency === 'once' ? (
          <Form.Item name="oneTimePushAt" label="一次性推送时间" rules={[{ required: enabled, message: '请选择一次性推送时间' }]}>
            <DatePicker showTime format="YYYY-MM-DD HH:mm" style={{ width: '100%' }} disabled={!enabled} />
          </Form.Item>
        ) : (
          <Form.Item name="pushTime" label="推送时间" rules={[{ required: enabled, message: '请选择推送时间' }]}>
            <TimePicker format="HH:mm" style={{ width: '100%' }} disabled={!enabled} />
          </Form.Item>
        )}
        {frequency === 'once' && <Typography.Paragraph type="secondary">推送一次将汇总从当前时间到所选推送时间之间的圈子内容。</Typography.Paragraph>}
        <Form.Item name="channels" label="推送方式"><Checkbox.Group disabled={!enabled} options={[{ value: 'notification', label: '站内通知' }, { value: 'email', label: '邮箱' }]} /></Form.Item>
        {channels.includes('email') && <Form.Item name="email" label="邮箱地址" rules={[{ required: true, message: '请输入邮箱地址' }, { type: 'email', message: '请输入正确的邮箱地址' }]}><Input /></Form.Item>}
        <Form.Item name="contentScopes" label="摘要内容"><Checkbox.Group disabled={!enabled} options={[{ value: 'post', label: '帖子' }, { value: 'comment', label: '评论' }, { value: 'chat', label: '聊天' }, { value: 'featured', label: '精华内容' }]} /></Form.Item>
      </Form>
    </Modal>
  );
}
