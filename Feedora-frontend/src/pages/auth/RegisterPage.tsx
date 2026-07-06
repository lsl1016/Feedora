import { Button, Card, Form, Input, Typography, message } from 'antd';
import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { api } from '../../api';
import { RegisterAgreementModal } from '../../components/auth/RegisterAgreementModal';
import { useAuthStore } from '../../store/authStore';

export function RegisterPage() {
  const [form] = Form.useForm(); const navigate = useNavigate(); const setAuth=useAuthStore(s=>s.setAuth); const [agreement,setAgreement]=useState<{open:boolean;account:string;password:string;loading:boolean}>({open:false,account:'',password:'',loading:false});
  async function submit(v:any){ try{await api.register(v); setAgreement({open:true,account:v.account,password:v.password,loading:false});}catch(e){message.error(e instanceof Error?e.message:'注册失败')} }
  async function agree(){ setAgreement(a=>({...a,loading:true})); const r=await api.login(agreement.account,agreement.password); setAuth(r.token,r.user); message.success('欢迎加入社区'); navigate('/'); }
  return <div className="auth-page"><Card className="auth-card"><Typography.Title level={2}>创建账号</Typography.Title><Typography.Text type="secondary">加入面向开发者的知识分享社区</Typography.Text><Form form={form} layout="vertical" onFinish={submit}><Form.Item name="account" rules={[{required:true,message:'请输入账号'},{min:4,message:'账号至少 4 位'}]}><Input placeholder="账号"/></Form.Item><Form.Item name="nickname" rules={[{required:true,message:'请输入昵称'}]}><Input placeholder="昵称"/></Form.Item><Form.Item name="password" rules={[{required:true,message:'请输入密码'},{min:6,message:'密码至少 6 位'}]}><Input.Password placeholder="密码"/></Form.Item><Form.Item name="confirmPassword" dependencies={['password']} rules={[{required:true,message:'请再次输入密码'},({getFieldValue})=>({validator(_,v){return !v||getFieldValue('password')===v?Promise.resolve():Promise.reject(new Error('两次密码不一致'))}})]}><Input.Password placeholder="确认密码"/></Form.Item><Button type="primary" block htmlType="submit">注册</Button></Form><p>已有账号？<Link to="/login">去登录</Link></p></Card><RegisterAgreementModal open={agreement.open} loading={agreement.loading} onAgree={agree}/></div>;
}
