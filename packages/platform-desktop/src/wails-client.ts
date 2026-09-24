export interface WailsBridgeClient {
  call<T>(binding: string, method: string, ...args: unknown[]): Promise<T>;
}

type WailsMethod = (...args: unknown[]) => Promise<unknown>;
type WailsBindings = Record<string, Record<string, WailsMethod>>;

declare global {
  interface Window {
    go?: Record<string, WailsBindings>;
  }
}

export class GlobalWailsBridgeClient implements WailsBridgeClient {
  constructor(private readonly namespace = 'bridge') {}

  async call<T>(binding: string, method: string, ...args: unknown[]): Promise<T> {
    const target = window.go?.[this.namespace]?.[binding]?.[method];
    if (!target) throw new Error(`Wails binding unavailable: ${this.namespace}.${binding}.${method}`);
    return target(...args) as Promise<T>;
  }
}
