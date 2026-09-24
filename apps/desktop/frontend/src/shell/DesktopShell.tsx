import {
  BellOutlined, FolderOutlined, HomeOutlined, SearchOutlined, SettingOutlined, ToolOutlined,
} from '@ant-design/icons';
import { Input, Tooltip } from 'antd';
import { NavLink, Outlet, useLocation } from 'react-router-dom';
import { desktopPrimaryNavigation } from '@feedora/features';

const icons = {
  community: <HomeOutlined />,
  workspace: <ToolOutlined />,
  repository: <FolderOutlined />,
  search: <SearchOutlined />,
  notification: <BellOutlined />,
  settings: <SettingOutlined />,
};

const sidebarItems: Record<string, string[]> = {
  community: ['推荐 Feed', '最新', '热门', '圈子', '话题'],
  workspace: ['总览', '笔记', '知识库', '资源', '最近编辑', '待发布'],
  repository: ['feedora', 'go-backend', 'wails-bridge'],
  search: ['统一搜索', '语义搜索', '代码搜索', 'AI 搜索'],
  notification: ['全部', '评论与回复', '点赞与收藏', '发布反馈', '系统消息'],
  settings: ['账户', 'Feedora Server', 'Agent Runtime', 'Workspace', 'Repositories', 'Indexing', '隐私与安全', '诊断', '关于'],
};

export function DesktopShell() {
  const location = useLocation();
  const active = desktopPrimaryNavigation.find((item) => location.pathname.startsWith(item.route)) ?? desktopPrimaryNavigation[0];

  return (
    <div className="desktop-root">
      <header className="desktop-titlebar">
        <strong>Feedora Desktop</strong>
        <span className="desktop-version">V1</span>
        <Input className="desktop-command" prefix={<SearchOutlined />} placeholder="搜索社区、笔记、仓库、文件或提问..." suffix="Ctrl K" />
      </header>

      <div className="desktop-body">
        <nav className="activity-rail">
          {desktopPrimaryNavigation.map((item) => (
            <Tooltip title={item.label} placement="right" key={item.id}>
              <NavLink to={item.route} className={({ isActive }) => `activity-item ${isActive ? 'active' : ''}`}>
                {icons[item.id]}<span>{item.label}</span>
              </NavLink>
            </Tooltip>
          ))}
        </nav>

        <aside className="context-sidebar">
          <div className="context-title">{active.label}</div>
          {sidebarItems[active.id].map((label, index) => (
            <div className={`context-item ${index === 0 ? 'active' : ''}`} key={label}>{label}</div>
          ))}
        </aside>

        <main className="main-canvas"><Outlet /></main>

        <aside className="ai-panel">
          <div className="ai-panel-title">✦ AI 助手</div>
          <div className="ai-context">当前上下文将由页面通过 AgentRuntimePort 注入。</div>
          <div className="ai-empty">选择帖子、笔记、仓库文件或搜索结果后，可在这里继续分析。</div>
        </aside>
      </div>

      <footer className="status-bar">
        <span>● Feedora Server: Offline</span>
        <span>● Local Index: Ready</span>
        <span>● Agent Runtime: Disconnected</span>
        <span>Repo: -</span>
      </footer>
    </div>
  );
}
