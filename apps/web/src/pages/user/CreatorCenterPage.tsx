import { Card, List } from 'antd';
import { Page } from '../../components/common/Page';
export function CreatorCenterPage(){return <Page><div className="metric-grid"><Card>累计发帖<h2>32</h2></Card><Card>累计获赞<h2>2345</h2></Card><Card>累计收藏<h2>560</h2></Card><Card>精华内容<h2>8</h2></Card></div><Card className="soft-card" title="优质内容列表"><List dataSource={['如何高效学习一门新的编程语言？','做程序员社区时，冷启动内容应该怎么准备？']} renderItem={x=><List.Item>{x}</List.Item>}/></Card></Page>}
