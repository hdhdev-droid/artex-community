"use client";

import { useCallback } from "react";

import { useTranslations } from "next-intl";

import toolSources from "./builtin-description-sources.json";
import parameterSources from "./builtin-parameter-sources.json";
import skillSources from "./builtin-skill-sources.json";

const parameterKeys = new Map(parameterSources.map((source, index) => [source, `p${index}`]));

// These are display summaries, never the descriptions sent to models or saved
// by editors. Both identity and the exact shipped source must match.
export function useBuiltinDescriptions() {
  const tools = useTranslations("builtinToolSummaries");
  const skills = useTranslations("builtinSkillDescriptions");
  const parameters = useTranslations("builtinParameterSummaries");
  const tool = useCallback(
    (key: string, source: string, system: boolean) => {
      if (!system || toolSources[key as keyof typeof toolSources] !== source || !tools.has(key)) return source;
      return tools(key);
    },
    [tools],
  );
  const skill = useCallback(
    (name: string, source: string) => {
      if (skillSources[name as keyof typeof skillSources] !== source || !skills.has(name)) return source;
      return skills(name);
    },
    [skills],
  );
  const parameter = useCallback(
    (source: string, system: boolean) => {
      const key = system ? parameterKeys.get(source) : undefined;
      return key && parameters.has(key) ? parameters(key) : source;
    },
    [parameters],
  );
  return { tool, skill, parameter };
}
