import { createCommunityApi, type HttpTransport } from '@feedora/api-client';
import { httpDelete, httpGet, httpPost, httpPut } from './request';

const transport: HttpTransport = {
  get: httpGet,
  post: httpPost,
  put: httpPut,
  delete: httpDelete,
};

// Web keeps browser-specific auth/error handling in request.ts.
// Endpoint resolution is now shared and platform-neutral.
export const api = createCommunityApi(transport);

// Compatibility exports for the existing module facades.
export const realApi = api;
export const mockApi = api;
