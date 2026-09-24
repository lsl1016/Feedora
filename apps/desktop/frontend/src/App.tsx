import { ConfigProvider } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import { RouterProvider } from 'react-router-dom';
import { feedoraTokens } from '@feedora/ui';
import { router } from './router';

export default function App() {
  return (
    <ConfigProvider locale={zhCN} theme={{ token: { colorPrimary: feedoraTokens.colorPrimary, borderRadius: feedoraTokens.radiusMd } }}>
      <RouterProvider router={router} />
    </ConfigProvider>
  );
}
