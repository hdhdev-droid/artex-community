"use client";

import { api } from "@/lib/api";

import { isLocale, type Locale } from "./config";
import { setActiveLocale } from "./locale-store";

// One installation-wide language, matching API errors and model output.
// A late settings read must not overwrite a more recent language save.
let revision = 0;
let saving = false;

export async function syncServerLocale() {
  if (saving) return;
  const request = ++revision;
  const settings = await api.settings();
  if (request === revision && isLocale(settings.output_language)) {
    setActiveLocale(settings.output_language);
  }
}

export async function saveServerLocale(locale: Locale) {
  ++revision;
  saving = true;
  try {
    const settings = await api.setSettings({ output_language: locale });
    if (isLocale(settings.output_language)) setActiveLocale(settings.output_language);
    return settings;
  } finally {
    saving = false;
  }
}
