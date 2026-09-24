import { SendOutlined } from '@ant-design/icons';
import { Avatar, Button, Card, Input, Space, Typography, message } from 'antd';
import { useEffect, useRef, useState } from 'react';
import { api } from '../../api';
import { useAuthStore } from '../../store/authStore';
import type { Circle, CircleChatMessage } from '../../types';

export function CircleChatPanel({ circle }: { circle: Circle }) {
  const { currentUser } = useAuthStore(); const [messages,setMessages]=useState<CircleChatMessage[]>([]); const [input,setInput]=useState(''); const [sending,setSending]=useState(false); const bottom=useRef<HTMLDivElement>(null);
  async function load(){setMessages(await api.getChatMessages(circle.circleId)); setTimeout(()=>bottom.current?.scrollIntoView({behavior:'smooth'}),50)}
  useEffect(()=>{load();},[circle.circleId]);
  const disabled = !circle.isJoined || circle.myStatus==='muted';
  async function send(){ if(disabled)return; if(!input.trim()){message.warning('请输入消息内容');return;} setSending(true); const m=await api.sendChatMessage(circle.circleId,input.trim()); setMessages(prev=>[...prev,m]); setInput(''); setSending(false); setTimeout(()=>bottom.current?.scrollIntoView({behavior:'smooth'}),50); }
  return <div className="circle-chat-panel"><div className="chat-list">{messages.length===0&&<Typography.Text type="secondary">暂无聊天消息，来发第一条吧</Typography.Text>}{messages.map(m=>m.messageType==='system'?<div key={m.messageId} className="chat-system">[系统] {m.content}</div>:<div key={m.messageId} className={`chat-message ${m.senderId===currentUser?.userId?'mine':''}`}><Avatar src={m.sender.avatar}/><div><Space><b>{m.sender.nickname}</b><span>{m.createdAt.slice(11,16)}</span></Space><Card size="small">{m.content}</Card></div></div>)}<div ref={bottom}/></div><Space.Compact className="chat-input"><Input.TextArea rows={2} value={input} onChange={e=>setInput(e.target.value)} disabled={disabled} placeholder={disabled ? (circle.isJoined?'你已被禁言，无法发言':'加入圈子后可参与聊天') : '请输入消息...'} /><Button type="primary" icon={<SendOutlined/>} loading={sending} onClick={send} disabled={disabled}>发送</Button></Space.Compact></div>;
}
