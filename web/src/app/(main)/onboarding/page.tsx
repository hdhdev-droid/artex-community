"use client";

import * as React from "react";

import Link from "next/link";

import { useTranslations } from "next-intl";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { isLocale, LOCALE_LABELS, LOCALES, type Locale } from "@/i18n/config";
import { useActiveLocale } from "@/i18n/locale-store";
import { saveServerLocale } from "@/i18n/server-locale";
import { api } from "@/lib/api";
import { offensiveApi } from "@/lib/offensive-api";
import { securityApi } from "@/lib/security-api";

type Status = "missing" | "configured" | "verified" | "failed" | "unknown";
type Check = { status: Status; at?: string; count?: number };
type Purpose = "ai" | "assessment" | "defense";
const unknown: Check = { status: "unknown" };

export default function OnboardingPage() {
  const t = useTranslations("onboarding");
  const locale = useActiveLocale();
  const [language, setLanguage] = React.useState<Locale>(locale);
  const [purpose, setPurpose] = React.useState<Purpose>("ai");
  const [checks, setChecks] = React.useState<Record<string, Check>>({
    model: unknown,
    runner: unknown,
    mcp: unknown,
    splunk: unknown,
  });
  const [loading, setLoading] = React.useState(true);
  const [busy, setBusy] = React.useState("");
  const [notice, setNotice] = React.useState("");
  const [servers, setServers] = React.useState<{ id: number; name: string }[]>([]);
  const [sources, setSources] = React.useState<{ id: string; name: string }[]>([]);
  const [serverID, setServerID] = React.useState("");
  const [sourceID, setSourceID] = React.useState("");
  const generation = React.useRef(0);

  React.useEffect(() => {
    try {
      const saved = sessionStorage.getItem("artex-onboarding-purpose");
      if (saved === "ai" || saved === "assessment" || saved === "defense") setPurpose(saved);
      const pending = sessionStorage.getItem("artex-onboarding-language");
      setLanguage(isLocale(pending) ? pending : locale);
    } catch {
      setLanguage(locale);
    }
  }, [locale]);

  const refresh = React.useCallback(async () => {
    const current = ++generation.current;
    setLoading(true);
    setChecks({ model: unknown, runner: unknown, mcp: unknown, splunk: unknown });
    const [model, runner, mcp, splunk] = await Promise.allSettled([
      api.getLLM(),
      offensiveApi.capabilities(),
      api.mcpSummary(),
      securityApi.sources(),
    ]);
    if (current !== generation.current) return;
    const enabledServers =
      mcp.status === "fulfilled" ? mcp.value.filter((s) => s.enabled).map(({ id, name }) => ({ id, name })) : [];
    const enabledSources =
      splunk.status === "fulfilled"
        ? splunk.value.sources.filter((s) => s.enabled && s.credential_set).map(({ id, name }) => ({ id, name }))
        : [];
    setServers(enabledServers);
    setSources(enabledSources);
    setServerID(enabledServers[0] ? String(enabledServers[0].id) : "");
    setSourceID(enabledSources[0]?.id ?? "");
    setChecks({
      model: {
        status:
          model.status === "rejected"
            ? "unknown"
            : model.value.configured && model.value.key_set
              ? "configured"
              : "missing",
      },
      runner: {
        status: runner.status === "rejected" ? "unknown" : runner.value.runner_configured ? "configured" : "missing",
      },
      mcp: { status: mcp.status === "rejected" ? "unknown" : enabledServers.length ? "configured" : "missing" },
      splunk: { status: splunk.status === "rejected" ? "unknown" : enabledSources.length ? "configured" : "missing" },
    });
    setLoading(false);
  }, []);

  React.useEffect(() => {
    void refresh();
    const focus = () => {
      if (document.visibilityState === "visible") void refresh();
    };
    window.addEventListener("focus", focus);
    return () => {
      ++generation.current;
      window.removeEventListener("focus", focus);
    };
  }, [refresh]);

  async function verify(kind: string) {
    const current = generation.current;
    setBusy(kind);
    let count: number | undefined;
    try {
      if (kind === "model") {
        const [config, profiles] = await Promise.all([api.getLLM(), api.llmProfiles()]);
        const active = profiles.find((p) => p.is_default);
        const result = await api.testLLM(
          config.provider,
          config.model,
          config.base_url,
          "",
          config.proxy ?? "",
          config.thinking_type ?? "",
          config.reasoning_effort ?? "",
          active ? Number(active.id) : undefined,
          active?.streaming ?? true,
          active?.session_header_key ?? "",
        );
        if (!result.ok) throw new Error("connection_failed");
      } else if (kind === "runner") {
        const catalog = await offensiveApi.catalog();
        count = catalog.tools.filter((tool) => tool.available).length;
        if (!count) throw new Error("no_tools");
      } else if (kind === "mcp") {
        count = (await api.refreshMcpServer(Number(serverID))).length;
      } else {
        const result = await securityApi.testSource(sourceID);
        if (!result.ok) throw new Error("connection_failed");
      }
      if (current === generation.current)
        setChecks((prev) => ({ ...prev, [kind]: { status: "verified", at: new Date().toLocaleTimeString(), count } }));
    } catch {
      // Provider error bodies may contain credentials. Show a local, translated summary.
      if (current === generation.current) setChecks((prev) => ({ ...prev, [kind]: { status: "failed" } }));
    } finally {
      setBusy("");
    }
  }

  const cards = [
    { id: "model", path: "/system/llm", required: purpose === "ai" },
    ...(purpose === "assessment" ? [{ id: "runner", path: "/offensive/assessments", required: true }] : []),
    ...(purpose === "defense" ? [{ id: "splunk", path: "/function/purple/sources", required: false }] : []),
    { id: "mcp", path: "/system/mcp", required: false },
  ];
  const destination =
    purpose === "ai" ? "/function/tasks" : purpose === "assessment" ? "/offensive/assessments" : "/defensive";
  return (
    <div className="mx-auto flex max-w-4xl flex-col gap-6">
      <div className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <h1 className="text-xl font-semibold">{t("title")}</h1>
          <p className="mt-1 text-sm text-muted-foreground">{t("description")}</p>
        </div>
        <Button variant="outline" onClick={() => void refresh()} disabled={loading || !!busy}>
          {t("refresh")}
        </Button>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>{t("language")}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <p className="text-sm text-muted-foreground">{t("languageHelp")}</p>
          <div className="flex flex-wrap gap-2">
            <select
              aria-label={t("language")}
              value={language}
              className="h-9 rounded-md border bg-background px-3 text-sm"
              onChange={(e) => {
                if (isLocale(e.target.value)) setLanguage(e.target.value);
              }}
            >
              {LOCALES.map((value) => (
                <option key={value} value={value}>
                  {LOCALE_LABELS[value]}
                </option>
              ))}
            </select>
            <Button
              disabled={!!busy}
              onClick={async () => {
                setBusy("language");
                setNotice("");
                try {
                  await saveServerLocale(language);
                  try {
                    sessionStorage.removeItem("artex-onboarding-language");
                  } catch {
                    // A blocked session store must not undo the saved server language.
                  }
                  setNotice("languageSaved");
                } catch {
                  setNotice("saveFailed");
                } finally {
                  setBusy("");
                }
              }}
            >
              {t("saveLanguage")}
            </Button>
          </div>
          <p role="status" className="text-sm">
            {notice ? t(notice) : ""}
          </p>
        </CardContent>
      </Card>
      <fieldset className="space-y-3">
        <legend className="font-medium">{t("purpose")}</legend>
        <div className="flex flex-wrap gap-4">
          {(["ai", "assessment", "defense"] as const).map((value) => (
            <Label key={value} className="flex items-center gap-2">
              <input
                type="radio"
                name="purpose"
                value={value}
                checked={purpose === value}
                onChange={() => {
                  setPurpose(value);
                  try {
                    sessionStorage.setItem("artex-onboarding-purpose", value);
                  } catch {
                    // The selection remains usable when browser storage is unavailable.
                  }
                }}
              />
              {t(`purposes.${value}`)}
            </Label>
          ))}
        </div>
        <p className="text-sm text-muted-foreground">{t(`purposeHelp.${purpose}`)}</p>
      </fieldset>
      <p className="rounded-md border p-3 text-sm text-muted-foreground">{t("externalNotice")}</p>
      <div className="grid gap-4 md:grid-cols-2">
        {cards.map(({ id, path, required }) => (
          <Card key={id}>
            <CardHeader>
              <CardTitle className="flex flex-wrap items-center justify-between gap-2">
                {t(`cards.${id}.title`)}
                <Badge variant="outline">{t(required ? "required" : "optional")}</Badge>
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <p className="text-sm text-muted-foreground">{t(`cards.${id}.help`)}</p>
              <p role="status" className="text-sm">
                <Badge variant="secondary">{t(`status.${checks[id].status}`)}</Badge>
                {checks[id].at && <span className="ml-2">{t("checkedAt", { time: checks[id].at })}</span>}
              </p>
              {checks[id].count !== undefined && (
                <p className="text-sm">{t("toolCount", { count: checks[id].count })}</p>
              )}
              {checks[id].status === "failed" && <p className="text-sm text-destructive">{t("checkFailed")}</p>}
              {(id === "mcp" || id === "splunk") && (
                <select
                  aria-label={t(`cards.${id}.title`)}
                  value={id === "mcp" ? serverID : sourceID}
                  disabled={!!busy}
                  className="h-9 w-full rounded-md border bg-background px-2 text-sm"
                  onChange={(e) => {
                    if (id === "mcp") setServerID(e.target.value);
                    else setSourceID(e.target.value);
                    setChecks((prev) => ({ ...prev, [id]: { status: "configured" } }));
                  }}
                >
                  <option value="" disabled>
                    {t("chooseConnection")}
                  </option>
                  {(id === "mcp" ? servers : sources).map((item) => (
                    <option key={item.id} value={item.id}>
                      {item.name}
                    </option>
                  ))}
                </select>
              )}
              <div className="flex flex-wrap gap-2">
                <Button variant="outline" asChild>
                  <Link href={`${path}?from=onboarding`}>{t("configure")}</Link>
                </Button>
                <Button
                  disabled={
                    loading ||
                    !!busy ||
                    checks[id].status === "missing" ||
                    checks[id].status === "unknown" ||
                    (id === "mcp" && !serverID) ||
                    (id === "splunk" && !sourceID)
                  }
                  onClick={() => void verify(id)}
                >
                  {t(busy === id ? "checking" : "verify")}
                </Button>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
      <Card>
        <CardHeader>
          <CardTitle>{t("next")}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <p className="text-sm text-muted-foreground">{t("nextHelp")}</p>
          <div className="flex flex-wrap gap-2">
            <Button asChild>
              <Link href={destination}>{t("continue")}</Link>
            </Button>
            <Button variant="outline" asChild>
              <Link href="/system/settings?from=onboarding">{t("optionalSettings")}</Link>
            </Button>
            <Button variant="ghost" asChild>
              <Link href="/function/tasks">{t("later")}</Link>
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
