import { CodeOutlined, FolderOpenOutlined, ReloadOutlined, SearchOutlined } from '@ant-design/icons';
import { Button, Card, Input, Space, Tag, Tree, Typography, message } from 'antd';
import { useMemo, useState } from 'react';
import { useDesktopStore } from '../store/desktopStore';

export function RepositoryPage() {
  const { repositories, repositoryFiles, selectedRepositoryId, selectRepository, addContext, askAi } = useDesktopStore();
  const [filePath, setFilePath] = useState(repositoryFiles[0]?.path);
  const [filter, setFilter] = useState('');
  const repo = repositories.find((item:any) => item.id === selectedRepositoryId) || repositories[0];

  const files = useMemo(() => repositoryFiles.filter((item:any) =>
    item.repositoryId === repo?.id && (!filter || item.path.toLowerCase().includes(filter.toLowerCase()))
  ), [repositoryFiles, repo?.id, filter]);

  const file = files.find((item:any) => item.path === filePath) || files[0] || repositoryFiles[0];

  const treeData = [
    { title:'apps', key:'apps', children:[{title:'desktop',key:'desktop',children:[{title:'main.go',key:'apps/desktop/main.go'}]}] },
    { title:'packages', key:'packages', children:[{title:'app-core',key:'core',children:[{title:'ports',key:'ports',children:[{title:'search.ts',key:'packages/app-core/src/ports/search.ts'}]}]}] },
    { title:'docs', key:'docs', children:[{title:'client-platform-architecture.md',key:'docs/client-platform-architecture.md'}] },
  ];

  return (
    <div className="repo-page">
      <aside className="repo-sidebar">
        <div className="context-title">仓库</div>
        {repositories.map((item:any) => (
          <div key={item.id} className={`repo-item ${item.id === repo?.id ? 'active' : ''}`} onClick={() => {
            selectRepository(item.id);
            const first = repositoryFiles.find((f:any) => f.repositoryId === item.id);
            if (first) setFilePath(first.path);
          }}>
            <FolderOpenOutlined /><div><b>{item.name}</b><span>{item.branch}</span></div>
          </div>
        ))}
        <Input prefix={<SearchOutlined />} placeholder="筛选文件..." value={filter} onChange={(e) => setFilter(e.target.value)} />
        <Tree treeData={treeData as any} defaultExpandAll onSelect={(keys) => {
          const key=String(keys[0]||'');
          if (repositoryFiles.some((item:any) => item.path === key)) setFilePath(key);
        }} />
      </aside>

      <main className="repo-main">
        <div className="page-heading-row repo-header">
          <div><Typography.Title level={3}>{repo?.name}</Typography.Title><Typography.Text type="secondary">本地仓库阅读、索引与 AI 辅助理解 · {repo?.branch}</Typography.Text></div>
          <Space><Button icon={<ReloadOutlined />}>重新索引</Button><Button>打开目录</Button></Space>
        </div>

        <Card className="code-card" title={file?.path} extra={<Space>
          <Tag>{file?.language}</Tag>
          <Button size="small" onClick={() => {
            addContext({id:'file-'+file?.path,type:'file',label:'文件: '+file?.path});
            message.success('文件已加入 AI Context');
          }}>加入 Context</Button>
          <Button size="small" type="primary" icon={<CodeOutlined />} onClick={() => {
            addContext({id:'file-'+file?.path,type:'file',label:'文件: '+file?.path});
            askAi('请解释这个文件负责什么，以及它和 Client Platform 其他模块的关系。');
          }}>AI 解释</Button>
        </Space>}>
          <pre className="code-preview">{file?.content}</pre>
        </Card>

        <Card title="最近索引" className="panel-card">
          <Space wrap>{files.map((item:any) => <Tag key={item.path} onClick={() => setFilePath(item.path)}>{item.path}</Tag>)}</Space>
        </Card>
      </main>
    </div>
  );
}
