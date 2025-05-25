import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

i18n.use(initReactI18next).init({
  lng: 'en', 
  fallbackLng: 'en',
  resources: {
    en: {
      translation: {
        item1: 'Item1',
        item2: 'Item2',
        content_item1: 'Content for Item1',
        content_item2: 'Content for Item2',
        language: 'Language',
      },
    },
    de: {
      translation: {
        item1: 'Item1',
        item2: 'Item2',
        content_item1: 'Inhalt für Artike1',
        content_item2: 'Inhalt für Artike2',
        language: 'Sprache',
      },
    },
  },
  interpolation: { escapeValue: false },
});

export default i18n;