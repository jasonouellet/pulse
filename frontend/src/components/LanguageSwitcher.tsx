import { useTranslation } from "react-i18next";
import { Globe } from "lucide-react";

export function LanguageSwitcher() {
  const { i18n } = useTranslation();

  // Langue courante (ex: "fr" ou "en")
  const currentLang = i18n.resolvedLanguage || i18n.language || "fr";

  const toggleLanguage = () => {
    const nextLanguage = currentLang.startsWith("fr") ? "en" : "fr";
    i18n.changeLanguage(nextLanguage);
  };

  return (
    <button
      onClick={toggleLanguage}
      type="button"
      className="inline-flex min-h-[44px] items-center gap-2 rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm font-medium text-slate-700 shadow-sm transition-colors hover:bg-slate-100 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-200 dark:hover:bg-slate-700"
    >
      <Globe className="h-4 w-4 text-slate-500" aria-hidden="true" />
      {/* Affiche FR ou EN selon la langue active */}
      <span className="uppercase font-semibold">{currentLang.slice(0, 2)}</span>
    </button>
  );
}