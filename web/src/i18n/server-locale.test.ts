import { beforeEach, describe, expect, it, vi } from "vitest";

const mocks = vi.hoisted(() => ({
  settings: vi.fn(),
  setSettings: vi.fn(),
  setActiveLocale: vi.fn(),
}));
vi.mock("@/lib/api", () => ({ api: mocks }));
vi.mock("./locale-store", () => ({ setActiveLocale: mocks.setActiveLocale }));

import { saveServerLocale, syncServerLocale } from "./server-locale";

describe("installation language synchronization", () => {
  beforeEach(() => vi.clearAllMocks());
  it("uses the server language on a different browser", async () => {
    mocks.settings.mockResolvedValue({ output_language: "ja" });
    await syncServerLocale();
    expect(mocks.setActiveLocale).toHaveBeenCalledWith("ja");
  });
  it("preserves the current UI language on save failure", async () => {
    mocks.setSettings.mockRejectedValue(new Error("offline"));
    await expect(saveServerLocale("es")).rejects.toThrow("offline");
    expect(mocks.setActiveLocale).not.toHaveBeenCalled();
  });
  it("ignores a settings read that predates a successful save", async () => {
    let finish!: (settings: { output_language: string }) => void;
    mocks.settings.mockReturnValue(
      new Promise((resolve) => {
        finish = resolve;
      }),
    );
    const pending = syncServerLocale();
    mocks.setSettings.mockResolvedValue({ output_language: "ko" });
    await saveServerLocale("ko");
    finish({ output_language: "en" });
    await pending;
    expect(mocks.setActiveLocale.mock.calls).toEqual([["ko"]]);
  });
});
