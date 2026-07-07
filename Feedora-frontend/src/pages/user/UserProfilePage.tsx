import { Avatar, Card, Space, Typography } from 'antd';
import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { api } from '../../api';
import { httpGet } from '../../api/request';
import { Page } from '../../components/common/Page';
import { PostList } from '../../components/post/PostList';
import type { PageResult, Post, User } from '../../types';
export function UserProfilePage(){const id=Number(useParams().userId);const [user,setUser]=useState<User>();const [data,setData]=useState<PageResult<Post>>();useEffect(()=>{httpGet<User>(`/users/${id}`).then(setUser).catch(()=>{});api.getPosts({page:1,pageSize:20,authorId:id,status:'published'}).then(setData)},[id]);return <Page narrow><Card className="soft-card"><Space><Avatar src={user?.avatar} size={72}/><div><Typography.Title level={3}>{user?.nickname}</Typography.Title><p>{user?.bio}</p><span>粉丝 {user?.followerCount} · 发帖 {user?.postCount} · 获赞 {user?.likeReceivedCount}</span></div></Space></Card><PostList data={data}/></Page>}
