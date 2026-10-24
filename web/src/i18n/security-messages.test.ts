import { createTranslator } from "next-intl";
import { describe, expect, it } from "vitest";

import en from "../../messages/en.json";
import es from "../../messages/es.json";
import ja from "../../messages/ja.json";
import ko from "../../messages/ko.json";
import zh from "../../messages/zh.json";

const namespaces = [
  "purple",
  "investigation",
  "investigationAgent",
  "Aegis",
  "ClosedLoop",
  "LiveLoop",
  "generatedMessages",
  "apiErrors",
  "builtinToolSummaries",
  "builtinSkillDescriptions",
  "builtinParameterSummaries",
  "mockMessages",
];

describe("security workflow messages", () => {
  for (const [locale, messages] of Object.entries({ en, es, ja, ko, zh })) {
    it(`formats ${locale} messages without ICU errors`, () => {
      const errors: string[] = [];
      const visit = (node: Record<string, unknown>, prefix: string) => {
        for (const [key, value] of Object.entries(node)) {
          const path = prefix ? `${prefix}.${key}` : key;
          if (typeof value === "object" && value !== null) {
            visit(value as Record<string, unknown>, path);
          } else if (typeof value === "string") {
            const t = createTranslator({
              locale,
              messages: { message: value },
              onError: (error) => errors.push(`${path}: ${error.message}`),
            });
            // This is a literal JSON example, consumed with t.raw in the UI.
            if (path === "LiveLoop.rulePlaceholder") {
              expect(() => JSON.parse(t.raw("message"))).not.toThrow();
              continue;
            }
            const values = Object.fromEntries(
              [...value.matchAll(/\{([A-Za-z_]\w*)[,}]/g)].map((match) => [match[1], 1]),
            );
            t("message", values);
          }
        }
      };
      for (const namespace of namespaces) {
        visit((messages as Record<string, unknown>)[namespace] as Record<string, unknown>, namespace);
      }
      expect(errors).toEqual([]);
    });
  }
});
