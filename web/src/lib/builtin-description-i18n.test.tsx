import type { ReactNode } from "react";

import { renderHook } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { describe, expect, it } from "vitest";

import en from "../../messages/en.json";
import es from "../../messages/es.json";
import ja from "../../messages/ja.json";
import ko from "../../messages/ko.json";
import zh from "../../messages/zh.json";
import { useBuiltinDescriptions } from "./builtin-description-i18n";
import tools from "./builtin-description-sources.json";
import parameters from "./builtin-parameter-sources.json";
import skills from "./builtin-skill-sources.json";

describe("built-in display summaries preserve editable source", () => {
  for (const [locale, messages] of Object.entries({ en, es, ja, ko, zh })) {
    it(`translates known defaults but preserves custom text in ${locale}`, () => {
      const { result, unmount } = renderHook(useBuiltinDescriptions, {
        wrapper: ({ children }: { children: ReactNode }) => (
          <NextIntlClientProvider locale={locale} messages={messages} timeZone="UTC">
            {children}
          </NextIntlClientProvider>
        ),
      });
      const [toolKey, toolSource] = Object.entries(tools)[0];
      const [skillName, skillSource] = Object.entries(skills)[0];
      const custom = "사용자 원문 / Evidence {payload} %s 中文";
      expect(result.current.tool(toolKey, toolSource, true)).not.toBe(toolSource);
      expect(result.current.tool(toolKey, toolSource, false)).toBe(toolSource);
      expect(result.current.tool("custom-tool", toolSource, true)).toBe(toolSource);
      expect(result.current.tool(toolKey, custom, true)).toBe(custom);
      expect(result.current.skill(skillName, skillSource)).toBe(
        messages.builtinSkillDescriptions[skillName as keyof typeof messages.builtinSkillDescriptions],
      );
      expect(result.current.skill(skillName, custom)).toBe(custom);
      expect(result.current.skill("custom-skill", skillSource)).toBe(skillSource);
      expect(result.current.parameter(parameters[0], false)).toBe(parameters[0]);
      expect(result.current.parameter(custom, true)).toBe(custom);
      unmount();
    });
  }
});
