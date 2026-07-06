import { BellOutlined, BookOutlined, DashboardOutlined, EditOutlined, FireOutlined, HomeOutlined, LogoutOutlined, NotificationOutlined, StarOutlined, TagsOutlined, TeamOutlined, TrophyOutlined, UserOutlined } from '@ant-design/icons';
import { Avatar, Badge, Button, Dropdown, Input, Layout, Space } from 'antd';
import type { MenuProps } from 'antd';
import { Link, NavLink, Outlet, useNavigate } from 'react-router-dom';
import { api } from '../../api';
import { useAuthStore } from '../../store/authStore';
import { AiChatDrawer } from '../ai/AiChatDrawer';

const { Header, Content } = Layout;

const navItems = [
  { to: '/', label: '首页', icon: <HomeOutlined />, end: true },
  { to: '/?feed=recommend', label: '推荐', icon: <FireOutlined /> },
  { to: '/hot', label: '热门', icon: <FireOutlined /> },
  { to: '/circles', label: '圈子', icon: <TeamOutlined /> },
  { to: '/topics', label: '话题', icon: <TagsOutlined /> },
  { to: '/activities', label: '活动', icon: <NotificationOutlined /> },
  { to: '/rankings', label: '排行榜', icon: <TrophyOutlined /> },
];

export function FrontLayout() {
  const navigate = useNavigate();
  const { isLogin, currentUser, logout } = useAuthStore();
  const menu: MenuProps['items'] = isLogin ? [
    { key: 'me', label: '我的主页', icon: <UserOutlined /> },
    { key: 'following', label: '我的关注', icon: <StarOutlined /> },
    { key: 'workspace', label: '工作空间', icon: <BookOutlined /> },
    { key: 'creator', label: '创作者中心' },
    { key: 'tasks', label: '任务中心' },
    { key: 'favorites', label: '我的收藏' },
    ...(currentUser?.role === 'admin' ? [{ key: 'admin', label: '后台管理', icon: <DashboardOutlined /> }] : []),
    { type: 'divider' },
    { key: 'logout', label: '退出登录', icon: <LogoutOutlined /> },
  ] : [{ key: 'login', label: '登录' }, { key: 'register', label: '注册' }];

  function onMenu({ key }: { key: string }) {
    if (key === 'me') navigate('/users/me');
    if (key === 'following') navigate('/following-center');
    if (key === 'workspace') navigate('/workspace');
    if (key === 'creator') navigate('/creator');
    if (key === 'tasks') navigate('/tasks');
    if (key === 'favorites') navigate('/users/me?tab=favorites');
    if (key === 'admin') navigate('/admin');
    if (key === 'login') navigate('/login');
    if (key === 'register') navigate('/register');
    if (key === 'logout') { logout(); navigate('/login'); }
  }

  return (
    <Layout className="front-layout">
      <Header className="front-header">
        <div className="front-header-inner">
          <Link className="brand" to="/"><span className="brand-mark">社</span><b>开发者知识社区</b></Link>
          <nav className="top-nav">
            {navItems.map((item) => (
              <NavLink key={item.to + item.label} to={item.to} end={item.end}>
                {item.icon}<span>{item.label}</span>
              </NavLink>
            ))}
          </nav>
          <Input.Search className="global-search" placeholder="搜索帖子、用户、话题、圈子" onSearch={(v) => v && navigate(`/search?keyword=${encodeURIComponent(v)}`)} onFocus={() => api.getHotKeywords()} />
          <Button type="primary" icon={<EditOutlined />} onClick={() => navigate(isLogin ? '/posts/create' : '/login')}>发布</Button>
          <Badge count={2}><Button shape="circle" icon={<BellOutlined />} onClick={() => navigate('/notifications')} /></Badge>
          <Dropdown menu={{ items: menu, onClick: onMenu }}><Avatar src={currentUser?.avatar} icon={<UserOutlined />} className="clickable" /></Dropdown>
        </div>
      </Header>
      <Content><Outlet /></Content>
      <AiChatDrawer />
    </Layout>
  );
}
