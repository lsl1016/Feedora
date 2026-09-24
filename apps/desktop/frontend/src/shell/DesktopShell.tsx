import {
  BellOutlined, FolderOutlined, HomeOutlined, SearchOutlined, SettingOutlined, ToolOutlined,
  PlusOutlined, DeleteOutlined,
} from '@ant-design/icons';
import { Button, Input, Space, Tag, Tooltip, Typography, message } from 'antd';
import { KeyboardEvent, useState } from 'react';
import { NavLink, Outlet, useLocation, useNavigate } from 'react-router-dom';
import { desktopPrimaryNavigation } from '@feedora/features';
import { useDesktopStore } from '../store/desktopStore';

const icons:any = {
  community: <HomeOutlined />,
  workspace: <ToolOutlined />,
  repository: <FolderOutlined />,
  search: <SearchOutlined />,
  notification: <BellOutlined />,
  settings: <SettingOutlined />,
};

const sidebars:Record<string,Array<{label:string;route?:string}>> = {
  community: [{label:'推荐 Feed',route:'/community'},{label:'最新',route:'/community'},{label:'热门',route:'/community'},{label:'圈子'},{label:'话题'}],
  workspace: [{label:'总览',route:'/workspace'},{label:'笔记',route:'/workspace/notes'},{label:'知识库',route:'/workspace/knowledge-bases'},{label:'资源'},{label:'最近编辑'},{label:'待发布'}],
  repository: [{label:'本地仓库',route:'/repository'},{label:'最近文件'},{label:'索引任务'}],
  search: [{label:'统一搜索',route:'/search'},{label:'语义搜索'},{label:'代码搜索'},{label:'AI 搜索'}],
  notification: [{label:'全部',route:'/notifications'},{label:'评论与回复'},{label:'点赞与收藏'},{label:'发布反馈'},{label:'系统消息'}],
  settings: [{label:'账户',route:'/settings'},{label:'Feedora Server',route:'/settings'},{label:'Agent Runtime',route:'/settings'},{label:'Workspace',route:'/settings'},{label:'Repositories',route:'/settings'},{label:'Indexing',route:'/settings'},{label:'隐私与安全',route:'/settings'},{label:'诊断',route:'/settings'}],
};

export function DesktopShell() {
  const location=useLocation();
  const navigate=useNavigate();
  const active=desktopPrimaryNavigation.find((item)=>location.pathname.startsWith(item.route))??desktopPrimaryNavigation[0];
  const { context, removeContext, clearContext, aiMessages, askAi, saveLastAiAsNote, selectedRepositoryId, repositories, notifications }=useDesktopStore();
  const [query,setQuery]=useState('');
  const [input,setInput]=useState('');
  const unread=notifications.filter((item:any)=>item.readStatus==='unread').length;
  const repo=repositories.find((item:any)=>item.id===selectedRepositoryId);

  const goSearch=()=>{ if(query.trim()) navigate('/search?q='+encodeURIComponent(query.trim())); };
  const onCommandKey=(event:KeyboardEvent<HTMLInputElement>)=>{ if(event.key==='Enter') goSearch(); };

  return (
    <div className="desktop-root">
      <header className="desktop-titlebar">
        <div className="brand-mark">◇</div><strong>Feedora Desktop</strong><span className="desktop-version">V1</span>
        <Input className="desktop-command" prefix={<SearchOutlined />} value={query} onChange={(e)=>setQuery(e.target.value)} onKeyDown={onCommandKey} placeholder="搜索社区、笔记、仓库、文件或提问..." suffix="Ctrl K" />
      </header>

      <div className="desktop-body">
        <nav className="activity-rail">
          {desktopPrimaryNavigation.map((item)=>(
            <Tooltip title={item.label} placement="right" key={item.id}>
              <NavLink to={item.route} className={({isActive})=>`activity-item ${isActive?'active':''}`}>
                <span className="activity-icon">{icons[item.id]}</span><span>{item.label}{item.id==='notification'&&unread>0?<sup className="rail-dot">{unread}</sup>:null}</span>
              </NavLink>
            </Tooltip>
          ))}
        </nav>

        <aside className="context-sidebar">
          <div className="context-title-row"><div className="context-title">{active.label}</div><Button type="text" size="small" icon={<PlusOutlined />} /></div>
          {sidebars[active.id].map((item,index)=>(
            <div className={`context-item ${index===0?'active':''}`} key={item.label} onClick={()=>item.route&&navigate(item.route)}>{item.label}</div>
          ))}
          {active.id==='workspace'&&<><div className="sidebar-section-title">置顶知识库</div>{['Feedora Desktop','Go Backend','RAG / AI','产品设计'].map((item)=><div className="sidebar-mini-row" key={item}>{item}</div>)}</>}
          {active.id==='repository'&&<><div className="sidebar-section-title">仓库</div>{repositories.map((item:any)=><div className="sidebar-mini-row" key={item.id}>{item.name}<span>{item.branch}</span></div>)}</>}
        </aside>

        <main className="main-canvas"><Outlet /></main>

        <aside className="ai-panel">
          <div className="ai-panel-title"><span>✦ AI 助手</span><Space><Button type="text" size="small" onClick={()=>message.info('Mock 会话历史')}>↶</Button><Button type="text" size="small" icon={<PlusOutlined />} /></Space></div>
          <div className="ai-context-card">
            <div className="ai-context-head"><b>当前上下文</b>{context.length>0&&<Button type="text" size="small" icon={<DeleteOutlined />} onClick={clearContext}>清空</Button>}</div>
            {context.length===0?<Typography.Text type="secondary">从帖子、笔记、仓库文件或搜索结果添加 Context。</Typography.Text>:<Space wrap>{context.map((item)=><Tag closable onClose={()=>removeContext(item.id)} key={item.id}>{item.label}</Tag>)}</Space>}
          </div>

          <div className="ai-chat-scroll">
            {aiMessages.map((item)=>(
              <div key={item.id} className={`ai-message-row ${item.role}`}>
                <b>{item.role==='user'?'我':'✦ AI 助手'}</b>
                <div className="ai-message-bubble">{item.content}</div>
              </div>
            ))}
          </div>

          <div className="ai-panel-actions">
            <Button size="small" onClick={()=>{const id=saveLastAiAsNote(); if(id)message.success('AI 回答已保存为笔记');}}>保存为笔记</Button>
            <Button size="small" onClick={()=>message.success('Mock：已插入当前笔记')}>插入当前笔记</Button>
            <Button size="small" onClick={()=>askAi('请把当前上下文整理成可发布到社区的草稿。')}>生成发布草稿</Button>
          </div>

          <Space.Compact className="ai-panel-input">
            <Input.TextArea autoSize={{minRows:2,maxRows:4}} value={input} onChange={(e)=>setInput(e.target.value)} onPressEnter={(e)=>{if(!e.shiftKey){e.preventDefault();askAi(input);setInput('');}}} placeholder="继续提问，或让 AI 帮你总结、分析、生成内容..." />
            <Button type="primary" onClick={()=>{askAi(input);setInput('');}}>发送</Button>
          </Space.Compact>
          <div className="ai-model-row">DeepSeek R1 · Mock Runtime</div>
        </aside>
      </div>

      <footer className="status-bar">
        <span className="status-ok">● Feedora Server: Online</span>
        <span className="status-ok">● Local Index: Ready</span>
        <span className="status-ok">● Agent Runtime: Connected</span>
        <span>Repo: {repo?.name || '-'} @ {repo?.branch || '-'}</span>
        <span className="status-version">v1.0.0-mock</span>
      </footer>
    </div>
  );
}
