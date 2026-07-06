import { Button, Card, Tabs, Typography } from 'antd';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../../api';
import { Page } from '../../components/common/Page';
import type { Activity } from '../../types';
export function ActivityCenterPage(){const nav=useNavigate();const [type,setType]=useState('all');const [list,setList]=useState<Activity[]>([]);useEffect(()=>{api.getActivities({page:1,pageSize:50,type}).then(r=>setList(r.list))},[type]);return <Page><Card className="soft-card page-title-card"><Typography.Title level={2}>活动中心</Typography.Title></Card><Card className="soft-card"><Tabs activeKey={type} onChange={setType} items={[{key:'all',label:'全部'},{key:'topic',label:'话题活动'},{key:'vote',label:'投票活动'},{key:'essay',label:'征文活动'},{key:'checkin',label:'打卡活动'}]}/><div className="activity-grid">{list.map(a=><Card key={a.activityId} hoverable cover={<img src={a.coverImage}/>} onClick={()=>nav(`/activities/${a.activityId}`)}><Typography.Title level={4}>{a.title}</Typography.Title><p>{a.description}</p><p>{a.participantCount} 人参与</p><Button type="primary" block>{a.status==='ended'?'查看结果':'立即参与'}</Button></Card>)}</div></Card></Page>}
