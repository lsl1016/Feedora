export interface NativePort {
  openExternal(url: string): Promise<void>;
  revealPath(path: string): Promise<void>;
  chooseDirectory(title?: string): Promise<string | null>;
}

export interface SecureStoragePort {
  get(key: string): Promise<string | null>;
  set(key: string, value: string): Promise<void>;
  delete(key: string): Promise<void>;
}
