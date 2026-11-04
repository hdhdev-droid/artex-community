import type { ReactNode } from "react";

import { cleanup, fireEvent, render, renderHook, screen, waitFor } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import MCPPage from "@/app/(main)/system/mcp/page";
import SkillsPage from "@/app/(main)/system/skills/page";

import en from "../../messages/en.json";
import es from "../../messages/es.json";
import ja from "../../messages/ja.json";
import ko from "../../messages/ko.json";
import zh from "../../messages/zh.json";
import { useAgentLabels } from "./agent-i18n";
import { api } from "./api";
import sources from "./builtin-agent-sources.json";

vi.mock("./api", () => ({
  api: {
    agents: vi.fn(),
    mcpServers: vi.fn(),
    resourceVisibility: vi.fn(),
    skills: vi.fn(),
    missingSkills: vi.fn(),
    skillVisibility: vi.fn(),
    skillUsage: vi.fn(),
  },
}));

const custom = "사용자 지정 / Evidence {payload} %s 中文";
const agents = Object.entries(sources).map(([key, source]) => ({
  ...source,
  key,
  id: key,
  builtin: key !== "retester",
  enabled: true,
}));

beforeEach(() => {
  vi.mocked(api.agents).mockResolvedValue([
    ...agents,
    { id: "custom", key: "custom", name: custom, description: custom, builtin: false },
  ] as Awaited<ReturnType<typeof api.agents>>);
  vi.mocked(api.mcpServers).mockResolvedValue([
    { id: 1, name: "User MCP", transport: "http", enabled: true },
  ] as Awaited<ReturnType<typeof api.mcpServers>>);
  vi.mocked(api.resourceVisibility).mockResolvedValue([]);
  vi.mocked(api.skills).mockResolvedValue([
    { name: "user-skill", description: custom, files: [], mcps: [], usage_agents: [], calls: 0, tasks: 0 },
  ]);
  vi.mocked(api.missingSkills).mockResolvedValue([]);
  vi.mocked(api.skillVisibility).mockResolvedValue([]);
  vi.mocked(api.skillUsage).mockResolvedValue([]);
});
afterEach(cleanup);

describe("agent labels in actual permission screens", () => {
  for (const [locale, messages] of Object.entries({ en, es, ja, ko, zh })) {
    const wrapper = ({ children }: { children: ReactNode }) => (
      <NextIntlClientProvider locale={locale} messages={messages} timeZone="UTC">
        {children}
      </NextIntlClientProvider>
    );

    it(`preserves customized source and translates defaults in ${locale}`, () => {
      const { result } = renderHook(useAgentLabels, { wrapper });
      for (const agent of agents) {
        const expected = messages.builtinAgents[agent.key as keyof typeof messages.builtinAgents];
        expect(result.current.name(agent)).toBe(expected.name);
        expect(result.current.description(agent)).toBe(expected.description);
        expect(result.current.name({ ...agent, name: custom })).toBe(custom);
        expect(result.current.description({ ...agent, description: custom })).toBe(custom);
      }
      // Compare against shipped defaults, not just the lookup fixture itself.
      expect(sources.auto.name).toBe("Auto");
      expect(sources.auto.description).toBe(zh.builtinAgents.auto.description);
      expect(result.current.name({ key: "goals", name: sources.goals.name, builtin: false })).toBe(sources.goals.name);
      expect(result.current.name({ key: "reporter", name: "보고서 작성" })).toBe(messages.builtinAgents.reporter.name);
      expect(result.current.name({ key: "reporter", name: custom })).toBe(custom);
    });

    it(`renders translated MCP permissions and unchanged custom names in ${locale}`, async () => {
      render(<MCPPage />, { wrapper });
      for (const agent of agents) {
        const expected = messages.builtinAgents[agent.key as keyof typeof messages.builtinAgents];
        expect(await screen.findByRole("checkbox", { name: expected.name, exact: true })).toBeTruthy();
      }
      expect(screen.getByRole("checkbox", { name: custom, exact: true })).toBeTruthy();
    });

    it(`renders translated skill detail and creation permissions in ${locale}`, async () => {
      render(<SkillsPage />, { wrapper });
      fireEvent.click(await screen.findByRole("button", { name: "user-skill", exact: true }));
      expect(
        await screen.findByRole("checkbox", { name: messages.builtinAgents.goals.name, exact: true }),
      ).toBeTruthy();
      expect(screen.getByRole("checkbox", { name: custom, exact: true })).toBeTruthy();
      fireEvent.click(screen.getByRole("button", { name: messages.skillsPage.newSkill, exact: true }));
      const tab = await screen.findByRole("tab", { name: messages.skillsPage.tabVisibility, exact: true });
      fireEvent.mouseDown(tab, { button: 0, ctrlKey: false });
      fireEvent.click(tab);
      await waitFor(() =>
        expect(screen.getByRole("checkbox", { name: messages.builtinAgents.retester.name, exact: true })).toBeTruthy(),
      );
      expect(screen.getByRole("checkbox", { name: custom, exact: true })).toBeTruthy();
    });
  }
});
