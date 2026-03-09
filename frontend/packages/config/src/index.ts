import { create } from 'zustand';

export interface TenantTheme {
  primaryColor: string;
  logoUrl: string;
  faviconUrl: string;
}

export interface TenantFeatures {
  liveRating: boolean;
  ecommerce: boolean;
  franchise: boolean;
  invoicing: boolean;
  reporting: boolean;
}

export interface TenantConfig {
  tenantId: string;
  tenantName: string;
  domain: string;
  theme: TenantTheme;
  features: TenantFeatures;
  defaultLanguage: string;
  supportedLanguages: string[];
  currency: string;
  timezone: string;
}

interface ConfigState {
  tenant: TenantConfig | null;
  isLoaded: boolean;
  setTenant: (config: TenantConfig) => void;
  reset: () => void;
}

export const useConfigStore = create<ConfigState>((set) => ({
  tenant: null,
  isLoaded: false,

  setTenant: (tenant) => {
    localStorage.setItem('tenant_id', tenant.tenantId);
    set({ tenant, isLoaded: true });
  },

  reset: () => {
    localStorage.removeItem('tenant_id');
    set({ tenant: null, isLoaded: false });
  },
}));

export function getTenantId(): string | null {
  return localStorage.getItem('tenant_id');
}
