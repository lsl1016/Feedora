import { Navigate, createHashRouter } from 'react-router-dom';
import { DesktopShell } from './shell/DesktopShell';
import { PlaceholderPage } from './shell/PlaceholderPage';

export const router = createHashRouter([
  {
    path: '/',
    element: <DesktopShell />,
    children: [
      { index: true, element: <Navigate to="/community" replace /> },
      { path: 'community', element: <PlaceholderPage title="社区" description="社区 Feed、帖子详情与内容沉淀入口。" /> },
      { path: 'workspace', element: <PlaceholderPage title="工作台" description="本地笔记、知识库、资源与最近工作。" /> },
      { path: 'workspace/notes', element: <PlaceholderPage title="笔记" description="Markdown 编辑器与本地知识沉淀。" /> },
      { path: 'workspace/knowledge-bases', element: <PlaceholderPage title="知识库" description="知识库目录、资源、协作与索引状态。" /> },
      { path: 'repository', element: <PlaceholderPage title="仓库" description="Git 仓库阅读、文件树、代码引用与索引。" /> },
      { path: 'search', element: <PlaceholderPage title="统一搜索" description="Notes / Code / Resources / Community 统一检索。" /> },
      { path: 'notifications', element: <PlaceholderPage title="通知" description="社区互动、发布反馈和桌面系统消息。" /> },
      { path: 'settings', element: <PlaceholderPage title="设置" description="Server、Runtime、Workspace、Indexing 与隐私设置。" /> },
    ],
  },
]);
