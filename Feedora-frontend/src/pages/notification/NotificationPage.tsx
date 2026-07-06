import { Card, List, Tabs } from 'antd';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../../api';
import { Page } from '../../components/common/Page';
import type { NotificationItem } from '../../types';
export function NotificationPage(){const nav=useNavigate();const [type,setType]=useState('all');const [list,setList]=useState<NotificationItem[]>([]);useEffect(()=>{api.getNotificationList().then(setList)},[]);const filtered=list.filter(n=>type==='all'||n.category===type);return <Page narrow><Card className="soft-card" title="通知中心"><Tabs activeKey={type} onChange={setType} items={[{key:'all',label:'全部'},{key:'interaction',label:'互动'},{key:'follow',label:'关注'},{key:'circle',label:'圈子'},{key:'review',label:'审核'},{key:'growth',label:'积分'},{key:'activity',label:'活动'},{key:'system',label:'系统'}]}/><List dataSource={filtered} renderItem={n=><List.Item className="clickable" onClick={()=>{api.markNotificationRead(n.notificationId); if(n.targetUrl)nav(n.targetUrl)}}><List.Item.Meta avatar={n.readStatus==='unread'?<span className="unread-dot"/>:null} title={n.title} description={`${n.content} · ${n.createdAt}`}/></List.Item>}/></Card></Page>}
