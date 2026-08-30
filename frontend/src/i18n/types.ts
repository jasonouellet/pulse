import frCommon from './locales/fr/common.json';
import frRoster from './locales/fr/roster.json';

export const defaultNS = 'common';
export const resources = {
  common: frCommon,
  roster: frRoster,
} as const;

declare module 'i18next' {
  interface CustomTypeOptions {
    defaultNS: typeof defaultNS;
    resources: typeof resources;
  }
}
