import { CloudOutlined, DatabaseOutlined, FolderOutlined, SafetyOutlined, ToolOutlined } from '@ant-design/icons';
import { Button, Card, Checkbox, Form, Input, Select, Space, Switch, Typography, message } from 'antd';
import { useEffect } from 'react';
import { useDesktopStore } from '../store/desktopStore';

export function SettingsPage() {
  const { settings, updateSettings, repositories } = useDesktopStore();
  const [form] = Form.useForm();

  useEffect(() => form.setFieldsValue(settings), [settings]);

  async function save() {
    const values = await form.validateFields();
    updateSettings(values);
    message.success('设置已保存到 Mock Desktop State');
  }

  return (
    <div className="desktop-page settings-page">
      <div className="page-heading-row">
        <div><Typography.Title level={2}>设置</Typography.Title><Typography.Text type="secondary">配置账户、运行时、本地工作区、仓库、索引与隐私安全。</Typography.Text></div>
        <Button type="primary" onClick={save}>保存全部设置</Button>
      </div>

      <Form form={form} layout="vertical">
        <div className="settings-grid">
          <Card title={<Space><CloudOutlined />Feedora Server</Space>} className="panel-card">
            <Form.Item name="serverUrl" label="服务器地址"><Input /></Form.Item>
            <Space><span className="status-ok">● 已连接（Mock）</span><Button onClick={()=>message.success('Feedora Server Mock 连接正常')}>测试连接</Button></Space>
          </Card>

          <Card title={<Space><ToolOutlined />Agent Runtime</Space>} className="panel-card">
            <Form.Item name="runtimeUrl" label="运行时地址"><Input /></Form.Item>
            <Space><span className="status-ok">● 已连接（Mock）</span><Button onClick={()=>message.success('Agent Runtime Mock 连接正常')}>测试连接</Button></Space>
          </Card>

          <Card title={<Space><FolderOutlined />Workspace 工作区</Space>} className="panel-card">
            <Form.Item name="workspacePath" label="工作区路径"><Input /></Form.Item>
            <Form.Item name="importPath" label="默认导入目录"><Input /></Form.Item>
          </Card>

          <Card title={<Space><DatabaseOutlined />Repositories 仓库</Space>} className="panel-card">
            {repositories.map((repo:any)=><div key={repo.id} className="settings-repo-row"><Checkbox defaultChecked /> <span>{repo.path}</span><TagLike>{repo.branch}</TagLike></div>)}
            <Button style={{marginTop:12}}>+ 添加仓库</Button>
          </Card>

          <Card title="Indexing 索引设置" className="panel-card">
            <Form.Item name="indexMode" label="索引模式"><Select options={[{value:'local',label:'纯本地索引'},{value:'hybrid',label:'混合模式（本地 + AI 增强）'}]} /></Form.Item>
            <Form.Item name="incrementalIndex" valuePropName="checked"><Switch /> <span className="setting-switch-label">增量索引</span></Form.Item>
            <Form.Item name="autoIndex" valuePropName="checked"><Switch /> <span className="setting-switch-label">自动索引新文件</span></Form.Item>
          </Card>

          <Card title={<Space><SafetyOutlined />隐私与安全</Space>} className="panel-card">
            <Form.Item name="allowAiLocalFiles" valuePropName="checked"><Switch /> <span className="setting-switch-label">允许 AI 处理本地文件内容</span></Form.Item>
            <Form.Item name="encryptLocalData" valuePropName="checked"><Switch /> <span className="setting-switch-label">本地数据加密</span></Form.Item>
            <Button onClick={()=>message.success('Mock 缓存已清理')}>清理缓存</Button>
          </Card>

          <Card title="诊断" className="panel-card">
            <p>Feedora Server：Online</p>
            <p>Local Index：Ready</p>
            <p>Agent Runtime：Connected</p>
            <Space><Button onClick={()=>message.success('诊断完成：一切正常')}>运行诊断检查</Button><Button>打开日志目录</Button></Space>
          </Card>
        </div>
      </Form>
    </div>
  );
}

function TagLike({children}:{children:any}) {
  return <span className="mini-tag">{children}</span>;
}
