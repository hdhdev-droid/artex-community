import { cleanup, render } from "@testing-library/react";
import { afterEach, expect, it, vi } from "vitest";

import { META_BY_LOCALE } from "@/config/app-config";

import type { Locale } from "./config";
import { IntlRuntimeProvider } from "./intl-runtime-provider";

const state = vi.hoisted(() => ({ locale: "en" as Locale }));
vi.mock("./locale-store", () => ({ useActiveLocale: () => state.locale }));
afterEach(cleanup);

it("synchronizes title, description and HTML language without changing page content", () => {
  const evidence = "Evidence 原文 / 사용자 입력";
  const view = render(
    <IntlRuntimeProvider>
      <p>{evidence}</p>
    </IntlRuntimeProvider>,
  );
  for (const locale of ["ko", "en", "zh", "es", "ja"] as const) {
    state.locale = locale;
    view.rerender(
      <IntlRuntimeProvider>
        <p>{evidence}</p>
      </IntlRuntimeProvider>,
    );
    expect(document.title).toBe(META_BY_LOCALE[locale].title);
    expect(document.head.querySelector('meta[name="description"]')?.getAttribute("content")).toBe(
      META_BY_LOCALE[locale].description,
    );
    expect(document.head.querySelectorAll("title")).toHaveLength(1);
    expect(document.head.querySelectorAll('meta[name="description"]')).toHaveLength(1);
    expect(document.documentElement.lang).toBe(locale);
    expect(view.container.textContent).toBe(evidence);
  }
});
