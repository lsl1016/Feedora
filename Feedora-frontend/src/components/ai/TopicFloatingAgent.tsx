import { RobotOutlined, CloseOutlined } from '@ant-design/icons';
import { Button, Card, Input, Space, Typography } from 'antd';
import { useState } from 'react';
import { api } from '../../api';

const actions = [
  { key: 'summarize_topic', label: '总结当前话题' },
  { key: 'recommend_posts', label: '推荐高质量帖子' },
  { key: 'generate_ideas', label: '生成发帖灵感' },
  { key: 'generate_learning_path', label: '生成学习路线' },
  { key: 'generate_interview_questions', label: '整理面试题' },
];
export function TopicFloatingAgent({ topicId, topicName }: { topicId: number; topicName: string }) {
  const [open, setOpen] = useState(false); const [answer, setAnswer] = useState(''); const [question, setQuestion] = useState(''); const [loading, setLoading] = useState(false);
  async function run(action: string, q?: string) { setLoading(true); const r = await api.topicAgent(topicId, action, q); setAnswer(r.answer); setLoading(false); }
  return <div className="topic-agent">{!open ? <Button type="primary" icon={<RobotOutlined />} onClick={() => setOpen(true)}>AI Agent</Button> : <Card className="topic-agent-panel" title="AI Agent" extra={<CloseOutlined onClick={() => setOpen(false)} />}><Typography.Text type="secondary">当前话题：#{topicName}#</Typography.Text><Space wrap style={{ margin: '12px 0' }}>{actions.map(a => <Button key={a.key} size="small" onClick={() => run(a.key)}>{a.label}</Button>)}</Space><Space.Compact style={{ width: '100%' }}><Input value={question} onChange={e => setQuestion(e.target.value)} placeholder="问问这个话题的 AI Agent..." /><Button loading={loading} type="primary" onClick={() => run('custom_question', question)}>发送</Button></Space.Compact>{answer && <Typography.Paragraph className="ai-answer-box">{answer}</Typography.Paragraph>}</Card>}</div>;
}
