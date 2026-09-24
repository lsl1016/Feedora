import { createCommunityApi, type HttpTransport } from '@feedora/api-client';
import { createMockFeedoraApi } from '@feedora/mock-data';
import { isMockMode } from './config';
import { httpDelete, httpGet, httpPost, httpPut } from './request';

const transport: HttpTransport = { get: httpGet, post: httpPost, put: httpPut, delete: httpDelete };

export const realApi = createCommunityApi(transport);
export const mockApi = createMockFeedoraApi();

// V1 可以完全脱离后端运行。生产接后端时只需要切 VITE_API_MODE=real。
export const api = isMockMode ? mockApi : realApi;
