import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import { user } from '../wailsjs/go/models';

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
        userProfile: 'User Profile',
        firstName: 'First Name',
        lastName: 'Last Name',
        userProfileUpdateSuccess: 'User profile updated successfully.',
        userProfileUpdateError: 'Error updating user profile.',
        saveButton: 'Save',
        loading: 'Loading...',
      },
    },
    de: {
      translation: {
        item1: 'Item1',
        item2: 'Item2',
        content_item1: 'Inhalt für Artike1',
        content_item2: 'Inhalt für Artike2',
        language: 'Sprache',
        userProfile: 'Benutzerprofil',
        firstName: 'Vorname',
        lastName: 'Nachname',
        userProfileUpdateSuccess: 'Benutzerprofil erfolgreich aktualisiert.',
        userProfileUpdateError: 'Fehler beim Aktualisieren des Benutzerprofils.',
        saveButton: 'Speichern',
        loading: 'Laden...',
      },
    },
  },
  interpolation: { escapeValue: false },
});

export default i18n;