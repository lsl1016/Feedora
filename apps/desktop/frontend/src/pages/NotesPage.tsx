import { FileTextOutlined, PlusOutlined, RobotOutlined, SaveOutlined, SendOutlined } from '@ant-design/icons';
import { Button, Input, List, Select, Space, Tag, Typography, message } from 'antd';
import { useEffect, useMemo, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { useDesktopStore } from '../store/desktopStore';

function markdownPreview(content: string) {
  return content
    .replace(/^### (.*)$/gm, '<h3>$1</h3>')
    .replace(/^## (.*)$/gm, '<h2>$1</h2>')
    .replace(/^# (.*)$/gm, '<h1>$1</h1>')
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
    .replace(/^- (.*)$/gm, '<li>$1</li>')
    .replace(/\n/g, '<br/>');
}

export function NotesPage() {
  const [params] = useSearchParams();
  const { notes, knowledgeBases, saveNote, publishNote, addContext, askAi } = useDesktopStore();
  const initial = Number(params.get('note')) || notes[0]?.noteId;
  const [selectedId, setSelectedId] = useState(initial);
  const selected = useMemo(() => notes.find((item:any) => item.noteId === selectedId) || notes[0], [notes, selectedId]);
  const [title, setTitle] = useState(selected?.title || '');
  const [content, setContent] = useState(selected?.content || '');
  const [kbId, setKbId] = useState(selected?.knowledgeBaseId || knowledgeBases[0]?.knowledgeBaseId);

  useEffect(() => {
    if (!selected) return;
    setTitle(selected.title);
    setContent(selected.content);
    setKbId(selected.knowledgeBaseId);
  }, [selectedId, selected?.updatedAt]);

  const persist = () => {
    const id = saveNote({ noteId:selected?.noteId, title, content, knowledgeBaseId:kbId });
    setSelectedId(id);
    message.success('笔记已保存');
  };

  const create = () => {
    const id = saveNote({ title:'未命名笔记', content:'# 未命名笔记\n\n开始记录你的知识...', knowledgeBaseId:knowledgeBases[0]?.knowledgeBaseId });
    setSelectedId(id);
  };

  return (
    <div className="notes-layout">
      <aside className="note-tree">
        <div className="note-tree-head"><b>Feedora Desktop</b><Button size="small" icon={<PlusOutlined />} onClick={create} /></div>
        <Input.Search placeholder="搜索笔记..." size="small" />
        <div className="note-group-title">产品文档</div>
        <List size="small" dataSource={notes} renderItem={(note:any) => (
          <List.Item className={`note-list-item ${note.noteId === selectedId ? 'active' : ''}`} onClick={() => setSelectedId(note.noteId)}>
            <Space><FileTextOutlined />{note.title}</Space>
          </List.Item>
        )} />
      </aside>

      <section className="note-editor">
        <div className="editor-header">
          <div>
            <Input value={title} onChange={(e) => setTitle(e.target.value)} className="editor-title-input" />
            <Space><Tag>{selected?.knowledgeBaseName || 'Feedora Desktop'}</Tag><span className="muted">上次编辑：{selected?.updatedAt}</span></Space>
          </div>
          <Space>
            <Button icon={<SaveOutlined />} type="primary" onClick={persist}>保存</Button>
            <Button icon={<RobotOutlined />} onClick={() => {
              addContext({id:'note-'+selected?.noteId,type:'note',label:'笔记: '+title});
              askAi('请优化当前笔记结构，并补充关键技术要点。');
              message.success('已让 AI 分析当前笔记');
            }}>AI 润色</Button>
            <Button icon={<SendOutlined />} onClick={() => { publishNote(selected?.noteId); message.success('已发布为社区帖子（Mock）'); }}>发布为帖子</Button>
          </Space>
        </div>
        <Select value={kbId} onChange={setKbId} options={knowledgeBases.map((kb:any) => ({value:kb.knowledgeBaseId,label:kb.name}))} style={{width:220,marginBottom:12}} />
        <Input.TextArea className="markdown-editor" value={content} onChange={(e) => setContent(e.target.value)} />
        <div className="reference-strip">
          <b>关联内容</b>
          <Tag>社区帖子 · Feedora Desktop V1</Tag>
          <Tag>代码片段 · app.go</Tag>
          <Tag>资源 · product_prd.pdf</Tag>
        </div>
      </section>

      <section className="note-preview">
        <div className="preview-tabs">预览</div>
        <div className="markdown-preview" dangerouslySetInnerHTML={{__html:markdownPreview(content)}} />
      </section>
    </div>
  );
}
