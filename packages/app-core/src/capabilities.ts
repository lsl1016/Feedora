export const CLIENT_CAPABILITIES = [
  'community',
  'remote-search',
  'remote-agent-runtime',
  'local-workspace',
  'local-files',
  'repository',
  'local-search',
  'secure-storage',
  'native-notification',
  'local-agent-runtime',
] as const;

export type ClientCapability = (typeof CLIENT_CAPABILITIES)[number];

export interface CapabilityRegistry {
  has(capability: ClientCapability): boolean;
  require(capability: ClientCapability): void;
  list(): readonly ClientCapability[];
}

export function createCapabilityRegistry(enabled: readonly ClientCapability[]): CapabilityRegistry {
  const values = new Set<ClientCapability>(enabled);
  return {
    has: (capability) => values.has(capability),
    require: (capability) => {
      if (!values.has(capability)) throw new Error(`client capability is unavailable: ${capability}`);
    },
    list: () => [...values],
  };
}

export const WEB_CAPABILITIES: readonly ClientCapability[] = [
  'community',
  'remote-search',
  'remote-agent-runtime',
];

export const DESKTOP_CAPABILITIES: readonly ClientCapability[] = [
  ...WEB_CAPABILITIES,
  'local-workspace',
  'local-files',
  'repository',
  'local-search',
  'secure-storage',
  'native-notification',
  'local-agent-runtime',
];
