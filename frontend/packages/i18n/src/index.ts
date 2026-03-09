import i18next from 'i18next';
import { initReactI18next } from 'react-i18next';

const defaultResources = {
  en: {
    translation: {
      common: {
        save: 'Save',
        cancel: 'Cancel',
        delete: 'Delete',
        edit: 'Edit',
        create: 'Create',
        search: 'Search',
        loading: 'Loading...',
        noResults: 'No results found',
        confirm: 'Confirm',
        back: 'Back',
      },
    },
  },
  tr: {
    translation: {
      common: {
        save: 'Kaydet',
        cancel: 'Iptal',
        delete: 'Sil',
        edit: 'Duzenle',
        create: 'Olustur',
        search: 'Ara',
        loading: 'Yukleniyor...',
        noResults: 'Sonuc bulunamadi',
        confirm: 'Onayla',
        back: 'Geri',
      },
    },
  },
};

export function initI18n(resources = defaultResources, lng = 'en') {
  return i18next.use(initReactI18next).init({
    resources,
    lng,
    fallbackLng: 'en',
    interpolation: {
      escapeValue: false,
    },
  });
}

export { i18next };
export { useTranslation } from 'react-i18next';
