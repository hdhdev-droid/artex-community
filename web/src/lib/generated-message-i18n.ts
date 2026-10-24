"use client";

import { useTranslations } from "next-intl";

import patterns from "./generated-message-patterns.json";

const escapePattern = (text: string) => text.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
const entries = Object.entries(patterns).map(([key, source]) => {
  const parameters = [...source.matchAll(/\{(p\d+)\}/g)].map((match) => match[1]);
  const expression = source
    .split(/\{p\d+\}/)
    .map(escapePattern)
    .join("([\\s\\S]*?)");
  return { key, parameters, pattern: new RegExp(`^${expression}$`) };
});

// Translate only known product-generated captions at render time. Tool output,
// user text, original logs, model prose and downloaded evidence stay verbatim.
export function useGeneratedMessage() {
  const t = useTranslations("generatedMessages");
  return (text: string): string => {
    const entry = entries.find(({ pattern }) => pattern.test(text));
    if (!entry) return text;
    const match = entry.pattern.exec(text);
    if (!match) return text;
    const values: Record<string, string> = {};
    entry.parameters.forEach((name, index) => {
      const raw = match[index + 1];
      // An upload omission reason is itself a known fixed caption.
      const nested = entries.find((item) => item.parameters.length === 0 && item.pattern.test(raw));
      values[name] = nested ? t(nested.key) : raw;
    });
    return t(entry.key, values);
  };
}
