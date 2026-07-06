import { Button, Modal, Typography } from 'antd';
export function RegisterAgreementModal({ open, loading, onAgree }: { open: boolean; loading?: boolean; onAgree: () => void }) {
  return <Modal title="欢迎加入开发者知识社区" open={open} closable={false} maskClosable={false} footer={<Button type="primary" loading={loading} block onClick={onAgree}>我已同意并进入社区</Button>}><Typography.Paragraph>本站主要面向程序员、开发者、IT 从业者和计算机学生，用于技术知识分享、学习交流、求职经验、面试讨论、项目复盘和个人知识沉淀。</Typography.Paragraph><ol className="agreement-list"><li>禁止发布政治敏感内容</li><li>禁止发布赌博、暴力、色情、违法违规内容</li><li>禁止广告刷屏、恶意引流</li><li>禁止人身攻击、辱骂、骚扰他人</li><li>禁止发布游戏外挂、灰产、诈骗等内容</li><li>请围绕技术学习、职业成长、知识分享进行交流</li></ol></Modal>;
}
