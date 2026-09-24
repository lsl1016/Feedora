import { Button, Card, Form, Input, message } from 'antd';
import { useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { api } from '../../api';
import { Page } from '../../components/common/Page';
export function PostEditPage(){const id=Number(useParams().postId);const nav=useNavigate();const [form]=Form.useForm();useEffect(()=>{api.getPost(id).then(p=>form.setFieldsValue(p))},[id]);async function save(){const v=await form.validateFields();await api.updatePost(id,v);message.success('保存成功');nav(`/posts/${id}`)}return <Page narrow><Card className="soft-card" title="编辑帖子"><Form form={form} layout="vertical"><Form.Item name="title" label="标题" rules={[{required:true}]}><Input/></Form.Item><Form.Item name="content" label="正文" rules={[{required:true}]}><Input.TextArea rows={10}/></Form.Item><Button type="primary" onClick={save}>保存修改</Button></Form></Card></Page>}
