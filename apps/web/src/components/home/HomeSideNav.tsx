import { BellOutlined, BookOutlined, FireOutlined, HomeOutlined, NotificationOutlined, ReadOutlined, StarOutlined, TeamOutlined, TrophyOutlined } from '@ant-design/icons';
import { NavLink } from 'react-router-dom';

const items = [
  { path: '/', label: '首页', icon: <HomeOutlined /> },
  { path: '/following-center', label: '我的关注', icon: <StarOutlined /> },
  { path: '/workspace', label: '工作空间', icon: <BookOutlined /> },
  { path: '/circles', label: '圈子', icon: <TeamOutlined /> },
  { path: '/topics', label: '话题', icon: <ReadOutlined /> },
  { path: '/activities', label: '活动', icon: <NotificationOutlined /> },
  { path: '/rankings', label: '排行榜', icon: <TrophyOutlined /> },
  { path: '/announcements', label: '官方公告', icon: <BellOutlined /> },
  { path: '/hot', label: '热门', icon: <FireOutlined /> },
];
export function HomeSideNav() {
  return <aside className="home-side-nav soft-card">{items.map(item => <NavLink key={item.path} to={item.path} end={item.path === '/'}>{item.icon}<span>{item.label}</span></NavLink>)}</aside>;
}
