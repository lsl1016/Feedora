export type EditorMode = 'edit' | 'preview' | 'split';

export interface EditorReference {
  type: 'community' | 'repository' | 'resource';
  label: string;
  target: string;
}

export interface MarkdownDocumentState {
  title: string;
  content: string;
  mode: EditorMode;
  references: EditorReference[];
  dirty: boolean;
}
