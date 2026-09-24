import type { SearchQuery, SearchResult } from '../domain/search';
import type { SearchPort } from '../ports/search';

export interface SearchSourcePort {
  name: string;
  port: SearchPort;
}

export async function searchAll(
  sources: readonly SearchSourcePort[],
  query: SearchQuery,
): Promise<SearchResult[]> {
  const groups = await Promise.all(sources.map(({ port }) => port.search(query)));
  return groups
    .flat()
    .sort((left, right) => (right.score ?? 0) - (left.score ?? 0))
    .slice(0, query.limit ?? 50);
}
