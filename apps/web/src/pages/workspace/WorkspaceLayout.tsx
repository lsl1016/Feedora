import { BookOutlined, DashboardOutlined, MessageOutlined, RobotOutlined, SettingOutlined } from '@ant-design/icons';
import { Layout, Menu } from 'antd';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';
const { Sider, Content } = Layout;
const items=[{key:'/workspace',label:'工作台',icon:<DashboardOutlined/>},{key:'/workspace/notes',label:'我的笔记',icon:<BookOutlined/>},{key:'/workspace/knowledge-bases',label:'我的知识库',icon:<MessageOutlined/>},{key:'/workspace/ai-chat',label:'AI 问答',icon:<RobotOutlined/>},{key:'/workspace/model-config',label:'模型配置',icon:<SettingOutlined/>}];
export function WorkspaceLayout(){const nav=useNavigate();const loc=useLocation();const selected=[...items].sort((a,b)=>b.key.length-a.key.length).find(i=>loc.pathname===i.key||loc.pathname.startsWith(i.key+'/'))?.key||'/workspace';return <Layout className="workspace-layout"><Sider width={220} className="workspace-sider"><h3>工作空间</h3><Menu mode="inline" selectedKeys={[selected]} items={items} onClick={({key})=>nav(key)}/></Sider><Content className="workspace-content"><Outlet/></Content></Layout>}
