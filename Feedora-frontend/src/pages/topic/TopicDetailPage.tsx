import { Button, Card, Select, Space, Typography } from 'antd';
import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { api } from '../../api';
import { TopicFloatingAgent } from '../../components/ai/TopicFloatingAgent';
import { Page } from '../../components/common/Page';
import { PostList } from '../../components/post/PostList';
import type { PageResult, Post, Topic } from '../../types';
export function TopicDetailPage(){const id=Number(useParams().topicId);const nav=useNavigate();const [topic,setTopic]=useState<Topic>();const [data,setData]=useState<PageResult<Post>>();const [sort,setSort]=useState<any>('latest');useEffect(()=>{api.getTopic(id).then(setTopic);api.getPosts({page:1,pageSize:20,topicId:id,sort}).then(setData)},[id,sort]);if(!topic)return null;return <Page><Card className="soft-card page-title-card"><div><Typography.Title level={2}># {topic.name}</Typography.Title><p>{topic.description}</p><span>{topic.participantCount} 人参与 · {topic.postCount} 篇内容</span></div><Button type="primary" onClick={()=>nav(`/posts/create?topicId=${id}`)}>参与话题</Button></Card><Card className="soft-card filter-card"><Space><Select value={sort} onChange={setSort} options={[{value:'latest',label:'最新'},{value:'hot',label:'热门'}]}/></Space></Card><PostList data={data}/><TopicFloatingAgent topicId={id} topicName={topic.name}/></Page>}
