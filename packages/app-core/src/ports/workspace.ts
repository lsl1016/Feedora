import type { KnowledgeBase, PageResult, WorkspaceNote } from '@feedora/contracts';

export interface NoteQuery {
  page: number;
  pageSize: number;
  keyword?: string;
  knowledgeBaseId?: number;
}

export interface CreateNoteInput {
  title: string;
  content: string;
  knowledgeBaseId?: number;
}

export interface WorkspacePort {
  listNotes(query: NoteQuery): Promise<PageResult<WorkspaceNote>>;
  getNote(noteId: number): Promise<WorkspaceNote>;
  createNote(input: CreateNoteInput): Promise<WorkspaceNote>;
  updateNote(noteId: number, input: Partial<CreateNoteInput>): Promise<WorkspaceNote>;
  deleteNote(noteId: number): Promise<void>;
  listKnowledgeBases(): Promise<KnowledgeBase[]>;
}
