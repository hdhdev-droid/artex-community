"use client";

import { useEffect, useRef, useState } from "react";

import Link from "next/link";

import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";
import { http } from "@/lib/api";

export default function OAuthCallback() {
  const t = useTranslations("modelOAuth");
  const [status, setStatus] = useState("connecting");
  const started = useRef(false);
  useEffect(() => {
    if (started.current) return;
    started.current = true;
    const query = new URLSearchParams(window.location.search);
    window.history.replaceState(null, "", window.location.pathname);
    const code = query.get("code");
    const state = query.get("state");
    if (query.has("error") || !code || !state) {
      setStatus("failed");
      return;
    }
    void http("/llm/oauth/openrouter/complete", { method: "POST", body: JSON.stringify({ code, state }) })
      .then(() => setStatus("connected"))
      .catch(() => setStatus("failed"));
  }, []);
  return (
    <div className="mx-auto max-w-xl space-y-4 rounded-lg border p-6">
      <h1 className="text-xl font-semibold">OpenRouter</h1>
      <p role="status">{t(status)}</p>
      <div className="flex flex-wrap gap-2">
        <Button asChild>
          <Link href="/system/llm?from=onboarding">{t("profiles")}</Link>
        </Button>
        <Button variant="outline" asChild>
          <Link href="/onboarding">{t("onboarding")}</Link>
        </Button>
      </div>
    </div>
  );
}
