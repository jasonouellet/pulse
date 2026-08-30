import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import { defaultNS } from './types';

import frCommon from './locales/fr/common.json';
import frRoster from './locales/fr/roster.json';
import enCommon from './locales/en/common.json';
import enRoster from './locales/en/roster.json';

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    fallbackLng: 'fr',
    defaultNS,
    resources: {
      fr: { common: frCommon, roster: frRoster },
      en: { common: enCommon, roster: enRoster },
    },
    interpolation: {
      escapeValue: false, // Inutile avec React qui échappe nativement le XSS
    },
    detection: {
      order: ['localStorage', 'navigator'],
      caches: ['localStorage'],
    },
  });

// Synchronise l'attribut html lang pour l'accessibilité WCAG / Loi 25
i18n.on('languageChanged', (lng) => {
  document.documentElement.lang = lng;
});

// Applique la langue initiale au premier chargement
if (typeof document !== 'undefined') {
  document.documentElement.lang = i18n.language || 'fr';
}

export default i18n;
