export interface LocalRepository {
  id: string;
  name: string;
  path: string;
  branch?: string;
  status?: 'ready' | 'indexing' | 'error';
}

export interface RepositoryFile {
  path: string;
  content: string;
  language?: string;
}

export interface RepositoryPort {
  listRepositories(): Promise<LocalRepository[]>;
  addRepository(path: string): Promise<LocalRepository>;
  removeRepository(repositoryId: string): Promise<void>;
  readFile(repositoryId: string, path: string): Promise<RepositoryFile>;
  reindex(repositoryId: string): Promise<void>;
}
