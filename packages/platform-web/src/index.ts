import { WEB_CAPABILITIES, createCapabilityRegistry, type CapabilityRegistry } from '@feedora/app-core';

export function createWebCapabilityRegistry(): CapabilityRegistry {
  return createCapabilityRegistry(WEB_CAPABILITIES);
}
