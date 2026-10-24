"use client";

import { useEffect, useRef, useState } from "react";

import Link from "next/link";

import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";
import { http } from "@/lib/api";

export default function ChatGPTCallback() {
  const t = useTranslations("chatgptOAuth");
  const started = useRef(false);
  const [account, setAccount] = useState<{ id: string; label: string } | null>(null);
  const [error, setError] = useState("");
  useEffect(() => {
    if (started.current) return;
    started.current = true;
    const params = new URLSearchParams(window.location.search);
    window.history.replaceState(null, "", window.location.pathname);
    const attempt = params.get("attempt");
    if (!attempt) {
      setError(t("failed"));
      return;
    }
    void http<{ id: string; label: string }>("/llm/oauth/chatgpt/complete", {
      method: "POST",
      body: JSON.stringify({ attempt }),
    })
      .then(setAccount)
      .catch((e: unknown) => setError(e instanceof Error ? e.message : t("failed")));
  }, [t]);
  return (
    <div className="mx-auto max-w-xl space-y-4 rounded-lg border p-6">
      <h1 className="text-xl font-semibold">{account ? t("welcome") : "ChatGPT"}</h1>
      {account ? (
        <>
          <p>{account.label}</p>
          <p>{t("welcomeBody")}</p>
        </>
      ) : (
        <p role="status">{error || t("working")}</p>
      )}
      <div className="flex flex-wrap gap-2">
        <Button asChild>
          <Link href={account ? `/system/llm?chatgpt_account=${encodeURIComponent(account.id)}` : "/system/llm"}>
            {t("profiles")}
          </Link>
        </Button>
        <Button variant="outline" asChild>
          <a href="https://chatgpt.com/settings/usage" target="_blank" rel="noopener noreferrer">
            {t("usage")}
          </a>
        </Button>
      </div>
    </div>
  );
}
