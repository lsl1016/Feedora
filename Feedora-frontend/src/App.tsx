import { ConfigProvider } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import { RouterProvider } from 'react-router-dom';
import { router } from './router';
import './styles/global.css';
export default function App(){return <ConfigProvider locale={zhCN} theme={{token:{colorPrimary:'#2563eb',borderRadius:12}}}><RouterProvider router={router}/></ConfigProvider>}
