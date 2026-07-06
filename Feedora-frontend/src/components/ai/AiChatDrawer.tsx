import { EditOutlined, FileAddOutlined, SendOutlined } from '@ant-design/icons';
import { Button, Drawer, Form, Input, Modal, Select, Space, Typography, message } from 'antd';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api, knowledgeBases } from '../../api';
import { useAppStore } from '../../store/appStore';
import type { AiChatMessage } from '../../types';

export function AiChatDrawer() {
  const { aiDrawerOpen, aiQuestion, closeAiDrawer } = useAppStore();
  const navigate = useNavigate();
  const [sessionId, setSessionId] = useState<number>();
  const [messages, setMessages] = useState<AiChatMessage[]>([]);
  const [input, setInput] = useState('');
  const [sending, setSending] = useState(false);
  const [editing, setEditing] = useState<AiChatMessage>();
  const [saveTarget, setSaveTarget] = useState<AiChatMessage>();
  const [form] = Form.useForm();

  useEffect(() => {
    if (!aiDrawerOpen) return;
    api.createAiSession(aiQuestion || 'AI 快问').then((session) => {
      setSessionId(session.sessionId);
      setMessages([]);
      if (aiQuestion) ask(session.sessionId, aiQuestion);
    });
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [aiDrawerOpen]);

  async function ask(id = sessionId, text = input) {
    if (!id || !text.trim()) { message.warning('请输入问题'); return; }
    setSending(true);
    const result = await api.sendAiMessage(id, text.trim());
    setMessages((prev) => [...prev, ...result]);
    setInput('');
    setSending(false);
  }

  async function saveEdit() {
    if (!editing) return;
    const value = form.getFieldValue('content');
    const updated = await api.updateAiMessage(editing.messageId, value);
    setMessages((prev) => prev.map((m) => m.messageId === editing.messageId ? updated! : m));
    setEditing(undefined);
  }

  async function saveToNote(values: { title: string; kbId?: number }) {
    if (!saveTarget) return;
    await api.saveAiToNote(saveTarget.messageId, values.title, values.kbId);
    message.success('已保存到笔记');
    setSaveTarget(undefined);
  }

  const lastAssistant = [...messages].reverse().find((m) => m.role === 'assistant');

  return <>
    <Drawer title="AI 助手" open={aiDrawerOpen} onClose={closeAiDrawer} width={520} maskClosable={false} extra={<Button onClick={() => navigate('/workspace/model-config')}>模型配置</Button>}>
      <div className="ai-chat-list">{messages.map(m => <div key={m.messageId} className={`ai-message ${m.role}`}><b>{m.role === 'user' ? '你' : 'AI'}</b><Typography.Paragraph>{m.content}</Typography.Paragraph>{m.role === 'assistant' && <Space><Button size="small" icon={<EditOutlined />} onClick={() => { setEditing(m); form.setFieldsValue({ content: m.content }); }}>编辑</Button><Button size="small" icon={<FileAddOutlined />} onClick={() => setSaveTarget(m)}>保存到笔记</Button><Button size="small" onClick={() => navigate(`/posts/create?sourceType=ai&answer=${encodeURIComponent(m.content)}`)}>发布成帖子</Button></Space>}</div>)}</div>
      {lastAssistant && <div className="ai-actions-sticky"><Space><Button onClick={() => setEditing(lastAssistant)}>编辑最新回答</Button><Button type="primary" onClick={() => setSaveTarget(lastAssistant)}>保存到笔记</Button></Space></div>}
      <Space.Compact className="ai-input"><Input.TextArea rows={2} value={input} onChange={e => setInput(e.target.value)} placeholder="继续追问..." onPressEnter={(e) => { if (!e.shiftKey) { e.preventDefault(); ask(); } }} /><Button type="primary" loading={sending} icon={<SendOutlined />} onClick={() => ask()}>发送</Button></Space.Compact>
    </Drawer>
    <Modal title="编辑 AI 回答" open={!!editing} onCancel={() => setEditing(undefined)} onOk={saveEdit} maskClosable={false}><Form form={form} layout="vertical"><Form.Item name="content" label="回答内容"><Input.TextArea rows={8} /></Form.Item></Form></Modal>
    <Modal title="保存到笔记" open={!!saveTarget} onCancel={() => setSaveTarget(undefined)} onOk={() => form.validateFields(['title','kbId']).then(saveToNote)} maskClosable={false}><Form form={form} layout="vertical"><Form.Item name="title" label="笔记标题" rules={[{ required: true, message: '请输入笔记标题' }]}><Input placeholder="AI 总结：学习路线" /></Form.Item><Form.Item name="kbId" label="所属知识库"><Select allowClear options={knowledgeBases.map(k => ({ value: k.knowledgeBaseId, label: k.name }))} /></Form.Item></Form></Modal>
  </>;
}
