import type {
  AgentContextRef,
  AgentMessage,
  AgentRuntimePort,
  LocalRepository,
  NativePort,
  RepositoryFile,
  RepositoryPort,
  SearchPort,
  SearchQuery,
  SearchResult,
  SecureStoragePort,
  WorkspacePort,
  NoteQuery,
  CreateNoteInput,
} from '@feedora/app-core';
import type { KnowledgeBase, PageResult, WorkspaceNote } from '@feedora/contracts';
import type { WailsBridgeClient } from './wails-client';

export class WailsWorkspaceAdapter implements WorkspacePort {
  constructor(private readonly bridge: WailsBridgeClient) {}
  listNotes(query: NoteQuery) { return this.bridge.call<PageResult<WorkspaceNote>>('WorkspaceBridge', 'ListNotes', query); }
  getNote(noteId: number) { return this.bridge.call<WorkspaceNote>('WorkspaceBridge', 'GetNote', noteId); }
  createNote(input: CreateNoteInput) { return this.bridge.call<WorkspaceNote>('WorkspaceBridge', 'CreateNote', input); }
  updateNote(noteId: number, input: Partial<CreateNoteInput>) { return this.bridge.call<WorkspaceNote>('WorkspaceBridge', 'UpdateNote', noteId, input); }
  deleteNote(noteId: number) { return this.bridge.call<void>('WorkspaceBridge', 'DeleteNote', noteId); }
  listKnowledgeBases() { return this.bridge.call<KnowledgeBase[]>('WorkspaceBridge', 'ListKnowledgeBases'); }
}

export class WailsRepositoryAdapter implements RepositoryPort {
  constructor(private readonly bridge: WailsBridgeClient) {}
  listRepositories() { return this.bridge.call<LocalRepository[]>('RepositoryBridge', 'ListRepositories'); }
  addRepository(path: string) { return this.bridge.call<LocalRepository>('RepositoryBridge', 'AddRepository', path); }
  removeRepository(repositoryId: string) { return this.bridge.call<void>('RepositoryBridge', 'RemoveRepository', repositoryId); }
  readFile(repositoryId: string, path: string) { return this.bridge.call<RepositoryFile>('RepositoryBridge', 'ReadFile', repositoryId, path); }
  reindex(repositoryId: string) { return this.bridge.call<void>('RepositoryBridge', 'Reindex', repositoryId); }
}

export class WailsSearchAdapter implements SearchPort {
  constructor(private readonly bridge: WailsBridgeClient) {}
  search(query: SearchQuery) { return this.bridge.call<SearchResult[]>('SearchBridge', 'Search', query); }
}

export class WailsAgentRuntimeAdapter implements AgentRuntimePort {
  constructor(private readonly bridge: WailsBridgeClient) {}
  createSession(title?: string) { return this.bridge.call<string>('AgentBridge', 'CreateSession', title ?? ''); }
  sendMessage(sessionId: string, message: string, context?: AgentContextRef[]) {
    return this.bridge.call<AgentMessage[]>('AgentBridge', 'SendMessage', sessionId, message, context ?? []);
  }
  cancel(sessionId: string) { return this.bridge.call<void>('AgentBridge', 'Cancel', sessionId); }
  health() { return this.bridge.call<{ connected: boolean; version?: string }>('AgentBridge', 'Health'); }
}

export class WailsNativeAdapter implements NativePort {
  constructor(private readonly bridge: WailsBridgeClient) {}
  openExternal(url: string) { return this.bridge.call<void>('SystemBridge', 'OpenExternal', url); }
  revealPath(path: string) { return this.bridge.call<void>('SystemBridge', 'RevealPath', path); }
  chooseDirectory(title?: string) { return this.bridge.call<string | null>('SystemBridge', 'ChooseDirectory', title ?? ''); }
}

export class WailsSecureStorageAdapter implements SecureStoragePort {
  constructor(private readonly bridge: WailsBridgeClient) {}
  get(key: string) { return this.bridge.call<string | null>('SystemBridge', 'SecureGet', key); }
  set(key: string, value: string) { return this.bridge.call<void>('SystemBridge', 'SecureSet', key, value); }
  delete(key: string) { return this.bridge.call<void>('SystemBridge', 'SecureDelete', key); }
}
