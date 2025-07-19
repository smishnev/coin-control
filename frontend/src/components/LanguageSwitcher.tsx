import React, { useState, useRef, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

const languages = [
  { code: 'de', label: 'Deutschland', flag: '🇩🇪' },
  { code: 'en', label: 'English', flag: '🇬🇧' },
];

const LanguageSwitcher: React.FC = () => {
  const { i18n } = useTranslation();
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  const currentLang = languages.find(l => l.code === i18n.language) || languages[0];

  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    };
    if (open) document.addEventListener('mousedown', handler);
    return () => document.removeEventListener('mousedown', handler);
  }, [open]);

  return (
    <div className="relative" ref={ref}>
      <button
        className="flex items-center gap-2 px-3 py-2 rounded-lg bg-background border border-border shadow-sm text-foreground font-sans font-medium hover:bg-menu-active transition"
        onClick={() => setOpen((v) => !v)}
        style={{ minWidth: 56 }}
      >
        <span className="text-xl font-normal">{currentLang.flag}</span>
        <span className="hidden sm:inline">{currentLang.code.toUpperCase()}</span>
        <svg className="ml-1 w-4 h-4 text-muted-foreground" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
          <path d="M19 9l-7 7-7-7" strokeLinecap="round" strokeLinejoin="round"/>
        </svg>
      </button>
      {open && (
        <div className="absolute left-0 mt-2 w-36 bg-menu border border-border rounded-lg shadow-lg z-10 animate-fade-in">
          {languages.map((l) => (
            <button
              key={l.code}
              className={`w-full flex items-center gap-2 px-3 py-2 text-left rounded-lg hover:bg-menu-active transition ${
                i18n.language === l.code ? 'bg-menu-active text-primary' : 'text-foreground'
              }`}
              onClick={() => {
                i18n.changeLanguage(l.code);
                setOpen(false);
              }}
            >
              <span className="text-xl font-normal">{l.flag}</span>
              <span>{l.label}</span>
            </button>
          ))}
        </div>
      )}
    </div>
  );
};

export default LanguageSwitcher;