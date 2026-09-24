import { RobotOutlined } from '@ant-design/icons';
import { Button, Card, Input, Space, Tag, Typography, message } from 'antd';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAppStore } from '../../store/appStore';
import { useAuthStore } from '../../store/authStore';

const quick = ['Java 后端面试路线', 'Go 项目怎么部署', 'React 学习路线', '如何写简历'];
export function AiQuickAskCard() {
  const [question, setQuestion] = useState('');
  const { isLogin } = useAuthStore();
  const navigate = useNavigate();
  const openAiDrawer = useAppStore(s => s.openAiDrawer);
  function ask(q = question) { if (!isLogin) { navigate('/login'); return; } if (!q.trim()) { message.warning('请输入问题'); return; } openAiDrawer(q); }
  return <Card className="soft-card ai-quick-card"><Typography.Title level={4}><RobotOutlined /> AI 快问</Typography.Title><Typography.Text type="secondary">用 AI 快速梳理学习路线、面试准备和项目思路</Typography.Text><Space.Compact className="ai-quick-input"><Input value={question} onChange={e => setQuestion(e.target.value)} placeholder="问问 AI：例如“如何准备 Java 后端面试？”" onPressEnter={() => ask()} /><Button type="primary" onClick={() => ask()}>发送</Button></Space.Compact><div className="quick-tags">{quick.map(q => <Tag key={q} color="blue" onClick={() => { setQuestion(q); ask(q); }}>{q}</Tag>)}</div></Card>;
}
