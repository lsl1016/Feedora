import { LockOutlined, UserOutlined } from '@ant-design/icons';
import { Button, Card, Form, Input, Typography, message } from 'antd';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import { api } from '../../api';
import { useAuthStore } from '../../store/authStore';

export function LoginPage() {
  const [form] = Form.useForm(); const navigate = useNavigate(); const [params] = useSearchParams(); const setAuth = useAuthStore(s=>s.setAuth);
  async function submit(v:{account:string;password:string}){ try{const r=await api.login(v.account,v.password); setAuth(r.token,r.user); message.success('登录成功'); navigate(params.get('redirect')||'/');}catch(e){message.error(e instanceof Error?e.message:'登录失败')} }
  return <div className="auth-page"><Card className="auth-card"><Typography.Title level={2}>开发者知识社区</Typography.Title><Typography.Text type="secondary">登录后开始学习、交流和沉淀知识</Typography.Text><Form form={form} layout="vertical" onFinish={submit} initialValues={{account:'zhangsan',password:'123456'}}><Form.Item name="account" rules={[{required:true,message:'请输入账号'}]}><Input prefix={<UserOutlined/>} placeholder="账号"/></Form.Item><Form.Item name="password" rules={[{required:true,message:'请输入密码'}]}><Input.Password prefix={<LockOutlined/>} placeholder="密码"/></Form.Item><Button type="primary" htmlType="submit" block>登录</Button></Form><p>没有账号？<Link to="/register">去注册</Link></p><Typography.Text type="secondary">管理员：admin / 123456</Typography.Text></Card></div>;
}
