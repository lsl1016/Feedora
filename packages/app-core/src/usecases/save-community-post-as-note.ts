import type { WorkspacePort } from '../ports/workspace';

export interface CommunityPostSnapshot {
  title: string;
  content: string;
  sourceUrl?: string;
}

export async function saveCommunityPostAsNote(
  workspace: WorkspacePort,
  post: CommunityPostSnapshot,
  knowledgeBaseId?: number,
) {
  const source = post.sourceUrl ? `\n\n> 来源：${post.sourceUrl}` : '';
  return workspace.createNote({
    title: post.title,
    content: `${post.content}${source}`,
    knowledgeBaseId,
  });
}
