export interface HttpTransport {
  get<T>(url: string, params?: unknown): Promise<T>;
  post<T>(url: string, data?: unknown): Promise<T>;
  put<T>(url: string, data?: unknown): Promise<T>;
  delete<T>(url: string, params?: unknown): Promise<T>;
}

export interface ResolvedEndpoint {
  method: 'get' | 'post' | 'put' | 'delete';
  url: string;
  params?: unknown;
  data?: unknown;
  transform?: (response: any) => any;
}
