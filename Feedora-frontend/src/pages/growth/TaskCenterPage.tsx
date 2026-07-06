import { Button, Card, Progress, Space, Tabs, Tag, message } from 'antd';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../../api';
import { Page } from '../../components/common/Page';
import type { Task } from '../../types';
export function TaskCenterPage(){const nav=useNavigate();const [type,setType]=useState('newbie');const [list,setList]=useState<Task[]>([]);function load(){api.getTasks(type).then(setList)}useEffect(load,[type]);async function claim(t:Task){await api.claimTask(t.taskId);message.success(`获得 ${t.rewardPoints} 积分`);load()}return <Page><Card className="soft-card" title="任务中心"><Tabs activeKey={type} onChange={setType} items={[{key:'newbie',label:'新手任务'},{key:'daily',label:'每日任务'},{key:'growth',label:'成长任务'}]}/><div className="task-grid">{list.map(t=><Card key={t.taskId}><h3>{t.title}</h3><p>{t.description}</p><Tag color="blue">奖励 {t.rewardPoints} 积分</Tag><Progress percent={Math.min(100,Math.round(t.currentValue/t.targetValue*100))}/>{t.status==='done'?<Button type="primary" onClick={()=>claim(t)}>领取奖励</Button>:t.status==='claimed'?<Button disabled>已完成</Button>:<Button onClick={()=>nav(t.actionUrl)}>{t.actionText}</Button>}</Card>)}</div></Card></Page>}
