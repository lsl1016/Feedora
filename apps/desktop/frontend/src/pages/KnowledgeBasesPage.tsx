import { DatabaseOutlined, FileAddOutlined, PlusOutlined, TeamOutlined } from '@ant-design/icons';
import { Button, Card, List, Modal, Form, Input, Space, Statistic, Tag, Typography, message } from 'antd';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useDesktopStore } from '../store/desktopStore';

export function KnowledgeBasesPage() {
  const nav = useNavigate();
  const { knowledgeBases, notes, createKnowledgeBase, addContext } = useDesktopStore();
  const [selectedId, setSelectedId] = useState(knowledgeBases[0]?.knowledgeBaseId);
  const [open, setOpen] = useState(false);
  const [form] = Form.useForm();
  const selected = knowledgeBases.find((item:any) => item.knowledgeBaseId === selectedId) || knowledgeBases[0];
  const kbNotes = notes.filter((note:any) => note.knowledgeBaseId === selected?.knowledgeBaseId);

  async function create() {
    const values = await form.validateFields();
    const id = createKnowledgeBase(values.name, values.description);
    setSelectedId(id);
    setOpen(false);
    form.resetFields();
    message.success('知识库已创建');
  }

  return (
    <div className="desktop-page">
      <div className="page-heading-row">
        <div><Typography.Title level={2}>知识库</Typography.Title><Typography.Text type="secondary">管理结构化笔记、资源与上下文。</Typography.Text></div>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setOpen(true)}>新建知识库</Button>
      </div>

      <div className="kb-browser">
        <Card className="kb-list-panel">
          {knowledgeBases.map((kb:any) => (
            <div key={kb.knowledgeBaseId} className={`kb-list-row ${kb.knowledgeBaseId === selected?.knowledgeBaseId ? 'active' : ''}`} onClick={() => setSelectedId(kb.knowledgeBaseId)}>
              <DatabaseOutlined /><div><b>{kb.name}</b><span>{kb.noteCount} 篇笔记</span></div>
            </div>
          ))}
        </Card>

        <div className="kb-detail-main">
          <Card className="kb-hero">
            <div className="page-heading-row">
              <div><Typography.Title level={3}>{selected?.name}</Typography.Title><Typography.Paragraph type="secondary">{selected?.description}</Typography.Paragraph></div>
              <Space>
                <Button icon={<FileAddOutlined />} onClick={() => nav('/workspace/notes')}>新建笔记</Button>
                <Button onClick={() => {
                  addContext({id:'kb-'+selected?.knowledgeBaseId,type:'knowledge-base',label:'知识库: '+selected?.name});
                  message.success('知识库已加入 Context');
                }}>加入 Context</Button>
              </Space>
            </div>
            <div className="kb-stats">
              <Statistic title="笔记" value={kbNotes.length} />
              <Statistic title="资源" value={8} />
              <Statistic title="协作者" value={selected?.memberCount || 1} prefix={<TeamOutlined />} />
              <Statistic title="本周更新" value={5} />
            </div>
          </Card>

          <div className="workspace-grid">
            <Card title="最近笔记" className="panel-card">
              <List dataSource={kbNotes} renderItem={(note:any) => <List.Item actions={[<Button size="small" onClick={() => nav('/workspace/notes?note='+note.noteId)}>打开</Button>]}><List.Item.Meta title={note.title} description={note.summary} /></List.Item>} />
            </Card>
            <Card title="相关资源" className="panel-card">
              <div className="resource-grid">
                <div className="resource-card">📄 product_prd.pdf<Tag>PDF</Tag></div>
                <div className="resource-card">🖼 ui-layout.png<Tag>Image</Tag></div>
                <div className="resource-card">📝 architecture.md<Tag>Markdown</Tag></div>
                <div className="resource-card">⌘ feedora<Tag>Repository</Tag></div>
              </div>
            </Card>
            <Card title="最近活动" className="panel-card">
              <p>Lin Chen 更新了《Feedora Desktop V1 架构拆解》</p>
              <p>Shirley 添加了资源 product_prd.pdf</p>
              <p>本地索引任务已完成</p>
            </Card>
          </div>
        </div>
      </div>

      <Modal title="新建知识库" open={open} onCancel={() => setOpen(false)} onOk={create}>
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="名称" rules={[{required:true}]}><Input /></Form.Item>
          <Form.Item name="description" label="描述"><Input.TextArea rows={3} /></Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
