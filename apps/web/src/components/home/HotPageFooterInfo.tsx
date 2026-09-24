import { Card, Typography } from 'antd';
export function HotPageFooterInfo() {
  return <Card className="soft-card footer-info" title="关于本站"><Typography.Paragraph>开发者知识社区，面向程序员、开发者、IT 从业者和计算机学生，用于技术知识分享、学习交流、求职经验、面试讨论和个人知识沉淀。</Typography.Paragraph><Typography.Paragraph><b>作者：</b>ll</Typography.Paragraph><Typography.Paragraph><b>项目说明：</b>本项目为社区 MVP 前端原型，使用 React + TypeScript + Ant Design 实现。</Typography.Paragraph><Typography.Text type="secondary">© 2026 Developer Knowledge Community. All rights reserved.</Typography.Text></Card>;
}
