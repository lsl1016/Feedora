export type SearchSource = 'note' | 'code' | 'resource' | 'community';

export interface SearchQuery {
  keyword: string;
  sources?: SearchSource[];
  limit?: number;
}

export interface SearchResult {
  id: string;
  source: SearchSource;
  title: string;
  summary?: string;
  location?: string;
  score?: number;
  metadata?: Record<string, unknown>;
}
