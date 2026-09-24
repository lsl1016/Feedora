import { CodeOutlined, FileTextOutlined, GlobalOutlined, SearchOutlined } from '@ant-design/icons';
import { Button, Card, Input, Select, Space, Tag, Typography, message } from 'antd';
import { useMemo, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { useDesktopStore } from '../store/desktopStore';

const iconFor = (source:string) => source === 'code' ? <CodeOutlined /> : source === 'community' ? <GlobalOutlined /> : <FileTextOutlined />;

export function SearchPage() {
  const [params, setParams] = useSearchParams();
  const { searchResults, notes, posts, addContext, askAi } = useDesktopStore();
  const [keyword, setKeyword] = useState(params.get('q') || 'Wails Bridge');
  const [source, setSource] = useState('all');

  const rows = useMemo(() => {
    const dynamic = [
      ...notes.map((note:any) => ({id:'note-'+note.noteId,source:'note',title:note.title,summary:note.summary,location:note.knowledgeBaseName,score:.85})),
      ...posts.map((post:any) => ({id:'post-'+post.postId,source:'community',title:post.title,summary:post.summary,location:post.circle?.name || 'Community',score:.78})),
      ...searchResults,
    ];
    const q=keyword.toLowerCase();
    return dynamic.filter((item:any) => (source === 'all' || item.source === source) && (!q || (item.title + item.summary + item.location).toLowerCase().includes(q)));
  }, [keyword, source, searchResults, notes, posts]);

  const run = () => setParams({q:keyword});

  return (
    <div className="desktop-page">
      <div className="page-heading-row">
        <div><Typography.Title level={2}>统一搜索</Typography.Title><Typography.Text type="secondary">一起搜索 Notes、Code、Resources 与 Community。</Typography.Text></div>
      </div>

      <Card className="search-hero">
        <Space.Compact style={{width:'100%'}}>
          <Input size="large" prefix={<SearchOutlined />} value={keyword} onChange={(e)=>setKeyword(e.target.value)} onPressEnter={run} />
          <Button size="large" type="primary" onClick={run}>搜索</Button>
        </Space.Compact>
        <div className="search-filter-row">
          <Select value={source} onChange={setSource} style={{width:160}} options={[
            {value:'all',label:'全部来源'},
            {value:'note',label:'Notes'},
            {value:'code',label:'Code'},
            {value:'resource',label:'Resources'},
            {value:'community',label:'Community'},
          ]} />
          <Tag>找到 {rows.length} 条结果</Tag>
        </div>
      </Card>

      <div className="search-results">
        {rows.map((item:any) => (
          <Card key={item.id} className="search-result-card">
            <div className="search-result-head">
              <Space>{iconFor(item.source)}<Tag>{item.source}</Tag><b>{item.title}</b></Space>
              <Button size="small" onClick={() => {
                addContext({id:'search-'+item.id,type:'search',label:'搜索结果: '+item.title});
                message.success('已加入 AI Context');
              }}>加入 Context</Button>
            </div>
            <Typography.Paragraph type="secondary">{item.summary}</Typography.Paragraph>
            <div className="muted">{item.location} · score {(item.score || 0).toFixed(2)}</div>
          </Card>
        ))}
      </div>

      <Button icon={<SearchOutlined />} onClick={() => askAi('请把当前搜索结果按设计、实现、发布三个方向归类。')}>让 AI 整理当前结果</Button>
    </div>
  );
}
