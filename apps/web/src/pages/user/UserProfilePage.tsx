import { Avatar, Button, Card, Space, Typography, message } from 'antd';
import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { api } from '../../api';
import { httpGet } from '../../api/request';
import { useAuthStore } from '../../store/authStore';
import { Page } from '../../components/common/Page';
import { PostList } from '../../components/post/PostList';
import type { PageResult, Post, User } from '../../types';
export function UserProfilePage(){const id=Number(useParams().userId);const {isLogin,currentUser}=useAuthStore();const [user,setUser]=useState<User>();const [data,setData]=useState<PageResult<Post>>();const [following,setFollowing]=useState(false);const isMine=id===currentUser?.userId;useEffect(()=>{httpGet<User>(`/users/${id}`).then(setUser).catch(()=>{});api.getPosts({page:1,pageSize:20,authorId:id,status:'published'}).then(setData).catch(()=>{});if(isLogin&&!isMine){api.getUserFollowState(id).then((r:{following?:boolean}|undefined)=>setFollowing(!!r?.following)).catch(()=>setFollowing(false))}},[id,isLogin,isMine]);async function toggleFollow(){try{if(following){await api.unfollowUser(id);setFollowing(false);message.success('已取消关注')}else{await api.followUser(id);setFollowing(true);message.success('关注成功')}httpGet<User>(`/users/${id}`).then(setUser).catch(()=>{})}catch{message.error('操作失败，请重试')}}return <Page narrow><Card className="soft-card"><Space><Avatar src={user?.avatar} size={72}/><div><Typography.Title level={3}>{user?.nickname}</Typography.Title><p>{user?.bio}</p><span>关注 {user?.followingCount} · 粉丝 {user?.followerCount} · 发帖 {user?.postCount} · 获赞 {user?.likeReceivedCount}</span></div>{isLogin&&!isMine&&<Button type={following?'default':'primary'} onClick={toggleFollow}>{following?'已关注':'+ 关注'}</Button>}</Space></Card><PostList data={data}/></Page>}
