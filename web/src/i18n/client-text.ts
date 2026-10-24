import { createTranslator } from "next-intl";

import type { CompanyScopeErrorCode } from "@/lib/company-scope";
import { MAX_COMPANY_SCOPE_VALUE_LENGTH } from "@/lib/company-scope";
import { getLocalStorageValue } from "@/lib/local-storage.client";

import en from "../../messages/en.json";
import es from "../../messages/es.json";
import ja from "../../messages/ja.json";
import ko from "../../messages/ko.json";
import zh from "../../messages/zh.json";
import { isLocale, LOCALE_STORAGE_KEY, resolveLocale } from "./config";

const messages = { en, es, ja, ko, zh };

function currentLocale() {
  const cached = getLocalStorageValue(LOCALE_STORAGE_KEY);
  return isLocale(cached) ? cached : resolveLocale();
}

// Non-React API paths read the same locale cache as the rendering provider.
export function apiErrorText(key: keyof typeof en.apiErrors): string {
  return messages[currentLocale()].apiErrors[key];
}

// Mock errors are new product messages, not stored demo evidence.
export function mockText(key: keyof typeof en.mockMessages, values?: Record<string, string | number>): string {
  const locale = currentLocale();
  return createTranslator({ locale, messages: messages[locale].mockMessages })(key, values);
}

export function scopeErrorText(code: CompanyScopeErrorCode | "" | undefined): string {
  if (!code) return "";
  const locale = currentLocale();
  return createTranslator({ locale, messages: messages[locale].scopeEditor.error })(code, {
    max: MAX_COMPANY_SCOPE_VALUE_LENGTH,
  });
}
