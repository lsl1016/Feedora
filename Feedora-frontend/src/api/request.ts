import axios, { type AxiosError, type AxiosRequestConfig } from 'axios';
import { message } from 'antd';
import { apiConfig } from './config';

export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
  traceId?: string;
  timestamp?: string;
}

export class ApiBusinessError extends Error {
  code: number;
  traceId?: string;
  data?: unknown;
  constructor(code: number, msg: string, traceId?: string, data?: unknown) {
    super(msg);
    this.name = 'ApiBusinessError';
    this.code = code;
    this.traceId = traceId;
    this.data = data;
  }
}

export const request = axios.create({
  baseURL: apiConfig.baseURL,
  timeout: apiConfig.timeout,
});

request.interceptors.request.use((config) => {
  const token = localStorage.getItem('community_v21_token') || localStorage.getItem('community_token');
  if (token) config.headers.Authorization = `Bearer ${token}`;
  if (apiConfig.enableLog) {
    console.info('[API Request]', config.method?.toUpperCase(), config.url, config.params || config.data || '');
  }
  return config;
});

request.interceptors.response.use(
  (response) => {
    const result = response.data as ApiResponse<unknown>;
    if (result && typeof result === 'object' && 'code' in result) {
      if (result.code !== 0) throw new ApiBusinessError(result.code, result.message || '业务处理失败', result.traceId, result.data);
      if (apiConfig.enableLog) console.info('[API Response]', response.config.url, result.traceId || '', result.data);
      return result.data;
    }
    return response.data;
  },
  (error: AxiosError<ApiResponse<unknown>>) => {
    const status = error.response?.status;
    const traceId = error.response?.data?.traceId;
    let text = '网络异常，请稍后再试';
    if (status === 401) {
      localStorage.removeItem('community_v21_token');
      localStorage.removeItem('community_v21_user');
      text = '登录已过期，请重新登录';
      if (!location.pathname.startsWith('/login')) location.href = `/login?redirect=${encodeURIComponent(location.pathname + location.search)}`;
    } else if (status === 403) text = '无权限访问';
    else if (status === 404) text = '资源不存在';
    else if (status === 422) text = error.response?.data?.message || '业务校验失败';
    else if (status === 429) text = '请求过于频繁，请稍后再试';
    else if (status && status >= 500) text = '服务异常，请稍后再试';
    message.error(traceId ? `${text}（traceId: ${traceId}）` : text);
    return Promise.reject(error);
  },
);

export function httpGet<T>(url: string, params?: unknown, config?: AxiosRequestConfig): Promise<T> {
  return request.get(url, { ...config, params }) as Promise<T>;
}
export function httpPost<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
  return request.post(url, data, config) as Promise<T>;
}
export function httpPut<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
  return request.put(url, data, config) as Promise<T>;
}
export function httpDelete<T>(url: string, params?: unknown, config?: AxiosRequestConfig): Promise<T> {
  return request.delete(url, { ...config, params }) as Promise<T>;
}
