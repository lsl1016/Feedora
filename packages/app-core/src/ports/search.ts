import type { SearchQuery, SearchResult } from '../domain/search';

export interface SearchPort {
  search(query: SearchQuery): Promise<SearchResult[]>;
}
