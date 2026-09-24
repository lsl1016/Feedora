import { resolveCommunityEndpoint } from './community-endpoints';
import type { HttpTransport } from './http';

export type ApiMethod = (...args: any[]) => Promise<any>;

export function createCommunityApi(transport: HttpTransport): Record<string, ApiMethod> {
  return new Proxy({} as Record<string, ApiMethod>, {
    get(_target, prop) {
      const methodName = String(prop);
      return async (...args: any[]) => {
        const endpoint = resolveCommunityEndpoint(methodName, args);
        if (endpoint.url.startsWith('/__unimplemented__')) {
          console.warn('[api] method is not mapped to Feedora Server:', methodName);
          return Promise.reject(new Error(`api.${methodName} is not implemented`));
        }

        let response: any;
        if (endpoint.method === 'get') response = await transport.get(endpoint.url, endpoint.params);
        else if (endpoint.method === 'post') response = await transport.post(endpoint.url, endpoint.data);
        else if (endpoint.method === 'put') response = await transport.put(endpoint.url, endpoint.data);
        else response = await transport.delete(endpoint.url, endpoint.params);

        return endpoint.transform ? endpoint.transform(response) : response;
      };
    },
  });
}
