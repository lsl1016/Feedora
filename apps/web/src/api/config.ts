export const apiConfig = {
  mode: (import.meta.env.VITE_API_MODE || 'mock') as 'mock' | 'real',
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
  timeout: Number(import.meta.env.VITE_API_TIMEOUT || 15000),
  enableLog: String(import.meta.env.VITE_ENABLE_API_LOG || 'true') === 'true',
};

export const isMockMode = apiConfig.mode !== 'real';
