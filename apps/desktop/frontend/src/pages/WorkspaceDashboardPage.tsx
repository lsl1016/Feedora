import { FileAddOutlined, FolderAddOutlined, ImportOutlined, RocketOutlined } from '@ant-design/icons';
import { Button, Card, List, Progress, Space, Statistic, Tag, Typography, message } from 'antd';
import { useNavigate } from 'react-router-dom';
import { useDesktopStore } from '../store/desktopStore';

export function WorkspaceDashboardPage() {
  const nav = useNavigate();
  const { notes, knowledgeBases, repositories, notifications, addContext, publishNote } = useDesktopStore();

  const metrics = [
    ['笔记', notes.length],
    ['知识库', knowledgeBases.length],
    ['资源', 42],
    ['仓库', repositories.length],
    ['待发布', 3],
  ];

  return (
    <div className="desktop-page">
      <div className="page-heading-row">
        <div>
          <Typography.Title level={2}>工作台</Typography.Title>
          <Typography.Text type="secondary">沉淀笔记、知识库、仓库资料与 AI 产出，形成可复用的个人知识工作流。</Typography.Text>
        </div>
        <Space>
          <Button type="primary" icon={<FileAddOutlined />} onClick={() => nav('/workspace/notes')}>新建笔记</Button>
          <Button icon={<ImportOutlined />}>导入文件</Button>
          <Button icon={<FolderAddOutlined />} onClick={() => nav('/repository')}>添加仓库</Button>
        </Space>
      </div>

      <div className="metric-grid-desktop">
        {metrics.map(([label, value]) => <Card key={label as string}><Statistic title={label} value={value as number} /></Card>)}
      </div>

      <div className="workspace-grid">
        <Card title="最近编辑" className="panel-card">
          <List dataSource={notes.slice(0, 4)} renderItem={(note:any) => (
            <List.Item actions={[
              <Button size="small" onClick={() => nav('/workspace/notes?note=' + note.noteId)}>继续编辑</Button>,
              <Button size="small" onClick={() => addContext({ id:'note-'+note.noteId, type:'note', label:'笔记: '+note.title })}>Context</Button>,
            ]}>
              <List.Item.Meta title={note.title} description={(note.updatedAt || '刚刚') + ' · ' + (note.knowledgeBaseName || '未分类')} />
            </List.Item>
          )} />
        </Card>

        <Card title="置顶知识库" className="panel-card">
          <div className="kb-tile-grid">
            {knowledgeBases.map((kb:any) => (
              <div className="kb-tile" key={kb.knowledgeBaseId} onClick={() => nav('/workspace/knowledge-bases')}>
                <b>{kb.name}</b>
                <p>{kb.description}</p>
                <span>{kb.noteCount} 篇笔记 · {kb.memberCount} 位成员</span>
              </div>
            ))}
          </div>
        </Card>

        <Card title="待发布笔记" className="panel-card">
          <List dataSource={notes.slice(0, 3)} renderItem={(note:any) => (
            <List.Item actions={[
              <Button type="primary" size="small" icon={<RocketOutlined />} onClick={() => {
                publishNote(note.noteId);
                message.success('已生成社区帖子（Mock）');
              }}>发布成帖子</Button>,
            ]}>
              <List.Item.Meta title={note.title} description={<Space><Tag>知识沉淀</Tag><Tag>草稿</Tag></Space>} />
            </List.Item>
          )} />
        </Card>

        <Card title="索引任务状态" className="panel-card">
          <div className="job-row"><span>feedora 仓库索引</span><Progress percent={100} size="small" /></div>
          <div className="job-row"><span>product_prd.pdf 文档解析</span><Progress percent={82} size="small" /></div>
          <div className="job-row"><span>notes 增量索引</span><Tag color="green">Ready</Tag></div>
        </Card>

        <Card title="最近 AI 会话" className="panel-card">
          <List dataSource={['总结架构设计文档','生成发布草稿','解释 repository indexing 流程']} renderItem={(item) => <List.Item>{item}</List.Item>} />
        </Card>

        <Card title="今日状态" className="panel-card">
          <p>未读通知：{notifications.filter((item:any) => item.readStatus === 'unread').length}</p>
          <p>本地索引：Ready</p>
          <p>Runtime：Mock Connected</p>
        </Card>
      </div>
    </div>
  );
}
