import { isMockMode } from './config';

export function createApiAdapter<T extends Record<string, any>>(mockApi: T, realApi: T): T {
  return isMockMode ? mockApi : realApi;
}
