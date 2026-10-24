import { afterEach, describe, expect, it, vi } from "vitest";

vi.mock("@/lib/mock/enabled", () => ({ MOCK: true }));

import { api } from "../api";
import { findings, tasks } from "./data";
import { sampleTaskReport } from "./reports";

afterEach(() => vi.restoreAllMocks());

describe("sample reports through the demo API", () => {
  it.each([
    ["t-acme-web", "총 13건 · 치명적 2 · 높음 6 · 중간 4 · 낮음 1"],
    ["t-acme-api", "총 3건 · 치명적 0 · 높음 2 · 중간 0 · 낮음 1"],
    ["t-shop-pay", "총 1건 · 치명적 0 · 높음 0 · 중간 1 · 낮음 0"],
    ["t-vpn-edge", "총 0건 · 치명적 0 · 높음 0 · 중간 0 · 낮음 0"],
  ])("returns task-specific findings for %s without network access", async (id, counts) => {
    const fetch = vi.spyOn(globalThis, "fetch").mockRejectedValue(new Error("Network prohibited"));
    const report = await api.report(id);
    expect(report).toContain(`작업 ID: ${id}`);
    expect(report).toContain(counts);
    const included = [...report.matchAll(/^### (f-\d+) ·/gm)].map((match) => match[1]);
    expect(included).toEqual(findings.filter((finding) => finding.task_id === id).map((finding) => finding.id));
    expect(fetch).not.toHaveBeenCalled();
  });

  it("uses the active sample by default and never substitutes it for an unknown task", async () => {
    expect(await api.report()).toBe(await api.report("t-acme-web"));
    expect(await api.report("t-missing")).toBe("");
  });

  it("preserves evidence and fixtures, including Markdown fence characters", () => {
    const evidence = "<script>example</script>\n```\n사용자 원문 {payload} %s";
    const item = { ...findings[0], evidence };
    const before = structuredClone(item);
    const report = sampleTaskReport(tasks[0], [item]);
    expect(report).toContain(`\`\`\`\`text\n${evidence}\n\`\`\`\``);
    expect(item).toEqual(before);
    expect(report).toContain("총 1건");
  });

  it("does not treat inherited object properties as sample task IDs", () => {
    for (const id of ["constructor", "toString", "__proto__"]) {
      expect(sampleTaskReport({ ...tasks[0], id }, findings)).toBe("");
    }
    expect(sampleTaskReport(undefined, findings)).toBe("");
  });
});
