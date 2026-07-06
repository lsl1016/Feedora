import { Card, Select, Typography } from 'antd';
import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { api, tags } from '../../api';
import { Page } from '../../components/common/Page';
import { PostList } from '../../components/post/PostList';
import type { PageResult, Post } from '../../types';
export function TagDetailPage(){const id=Number(useParams().tagId);const tag=tags.find(t=>t.tagId===id)||tags[0];const [sort,setSort]=useState<any>('latest');const [data,setData]=useState<PageResult<Post>>();useEffect(()=>{api.getPosts({page:1,pageSize:20,tagId:id,sort}).then(setData)},[id,sort]);return <Page><Card className="soft-card page-title-card"><div><Typography.Title level={2}># {tag.tagName}</Typography.Title><p>{tag.description}</p><span>相关内容数量：{tag.useCount}</span></div><Select value={sort} onChange={setSort} options={[{value:'latest',label:'最新'},{value:'hot',label:'热门'}]}/></Card><PostList data={data}/></Page>}
