"use client";

import { useTranslations } from "next-intl";

import { Label } from "@/components/ui/label";
import { isLocale, LOCALE_LABELS, LOCALES } from "@/i18n/config";
import { setActiveLocale, useActiveLocale } from "@/i18n/locale-store";

export function SetupLanguage() {
  const t = useTranslations("onboarding");
  const locale = useActiveLocale();
  return (
    <div className="grid gap-2">
      <Label htmlFor="setup-language">{t("language")}</Label>
      <select
        id="setup-language"
        value={locale}
        className="h-9 rounded-md border bg-background px-3 text-sm"
        onChange={(event) => {
          const value = event.target.value;
          if (!isLocale(value)) return;
          setActiveLocale(value);
          // Only a display preference; never persist credentials or change server settings before authentication.
          try {
            sessionStorage.setItem("artex-onboarding-language", value);
          } catch {
            /* Storage may be disabled. */
          }
        }}
      >
        {LOCALES.map((value) => (
          <option key={value} value={value}>
            {LOCALE_LABELS[value]}
          </option>
        ))}
      </select>
    </div>
  );
}
