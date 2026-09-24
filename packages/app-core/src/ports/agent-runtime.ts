export interface AgentContextRef {
  type: 'post' | 'note' | 'knowledge-base' | 'repository' | 'file' | 'selection';
  id: string;
  label: string;
  metadata?: Record<string, unknown>;
}

export interface AgentMessage {
  role: 'user' | 'assistant' | 'system' | 'tool';
  content: string;
}

export interface AgentRuntimePort {
  createSession(title?: string): Promise<string>;
  sendMessage(sessionId: string, message: string, context?: AgentContextRef[]): Promise<AgentMessage[]>;
  cancel(sessionId: string): Promise<void>;
  health(): Promise<{ connected: boolean; version?: string }>;
}
