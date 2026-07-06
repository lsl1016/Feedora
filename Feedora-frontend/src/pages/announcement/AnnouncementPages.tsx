import { Card, List, Typography } from 'antd';
import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { api } from '../../api';
import { Page } from '../../components/common/Page';
import type { OfficialAnnouncement } from '../../types';
export function AnnouncementListPage(){const nav=useNavigate();const [list,setList]=useState<OfficialAnnouncement[]>([]);useEffect(()=>{api.getOfficialAnnouncements({page:1,pageSize:50}).then(r=>setList(r.list))},[]);return <Page narrow><Card className="soft-card" title="官方公告"><List dataSource={list} renderItem={a=><List.Item onClick={()=>nav(`/announcements/${a.announcementId}`)} className="clickable"><List.Item.Meta title={a.title} description={`${a.summary} · ${a.publishedAt}`}/></List.Item>}/></Card></Page>}
export function AnnouncementDetailPage(){const id=Number(useParams().announcementId);const [a,setA]=useState<OfficialAnnouncement>();useEffect(()=>{api.getOfficialAnnouncements({page:1,pageSize:100}).then(r=>setA(r.list.find(x=>x.announcementId===id)||r.list[0]))},[id]);return <Page narrow><Card className="soft-card"><Typography.Title level={2}>{a?.title}</Typography.Title><Typography.Text type="secondary">{a?.publisherName} · {a?.publishedAt}</Typography.Text><Typography.Paragraph className="post-body-text">{a?.content}</Typography.Paragraph></Card></Page>}
