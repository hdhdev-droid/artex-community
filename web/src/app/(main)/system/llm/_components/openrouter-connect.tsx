"use client";

import { useState } from "react";

import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { http } from "@/lib/api";

export function OpenRouterConnect() {
  const t = useTranslations("modelOAuth");
  const [name, setName] = useState("OpenRouter");
  const [model, setModel] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState(false);
  return (
    <Card>
      <CardHeader>
        <CardTitle>OpenRouter</CardTitle>
      </CardHeader>
      <CardContent className="space-y-3">
        <p className="text-sm text-muted-foreground">{t("openrouterHelp")}</p>
        <div className="grid gap-3 sm:grid-cols-2">
          <div className="grid gap-2">
            <Label htmlFor="oauth-profile-name">{t("profileName")}</Label>
            <Input id="oauth-profile-name" maxLength={120} value={name} onChange={(e) => setName(e.target.value)} />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="oauth-model">{t("model")}</Label>
            <Input
              id="oauth-model"
              maxLength={200}
              placeholder="provider/model"
              value={model}
              onChange={(e) => setModel(e.target.value)}
            />
          </div>
        </div>
        <Button
          disabled={busy || !name.trim() || !model.trim()}
          onClick={async () => {
            setBusy(true);
            setError(false);
            try {
              const result = await http<{ url: string }>("/llm/oauth/openrouter/start", {
                method: "POST",
                body: JSON.stringify({ origin: window.location.origin, name, model }),
              });
              window.location.assign(result.url);
            } catch {
              setError(true);
              setBusy(false);
            }
          }}
        >
          {t(busy ? "connecting" : "openrouterConnect")}
        </Button>
        {error && (
          <p role="alert" className="text-sm text-destructive">
            {t("failed")}
          </p>
        )}
      </CardContent>
    </Card>
  );
}
