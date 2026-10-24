"use client";

import { useCallback, useEffect, useState } from "react";

import { useTranslations } from "next-intl";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { http } from "@/lib/api";

type Account = { id: string; label: string; connected: boolean };
type Model = { slug: string; display_name: string };
const root = "/llm/oauth/chatgpt";

export function ChatGPTConnect({ onSaved }: { onSaved: () => void }) {
  const t = useTranslations("chatgptOAuth");
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [accountID, setAccountID] = useState("new");
  const [models, setModels] = useState<Model[]>([]);
  const [model, setModel] = useState("");
  const [name, setName] = useState("ChatGPT");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  const [loaded, setLoaded] = useState(false);
  const selected = accounts.find((a) => a.id === accountID);
  const reload = useCallback(async () => {
    const result = await http<{ accounts: Account[] }>(`${root}/accounts`);
    setAccounts(result.accounts);
    setLoaded(true);
  }, []);

  useEffect(() => {
    // Only local connection metadata is read automatically, never the model API.
    void reload().catch(() => setError(t("unavailable")));
    const params = new URLSearchParams(window.location.search);
    const saved = params.get("chatgpt_account");
    if (saved) {
      setAccountID(saved);
      params.delete("chatgpt_account");
      window.history.replaceState(null, "", `${window.location.pathname}${params.size ? `?${params}` : ""}`);
    }
  }, [reload, t]);

  async function run(action: () => Promise<void>) {
    if (busy) return;
    setBusy(true);
    setError("");
    try {
      await action();
    } catch (e) {
      setError(e instanceof Error ? e.message : t("failed"));
    } finally {
      setBusy(false);
    }
  }

  function connect() {
    void run(async () => {
      const result = await http<{ url: string }>(`${root}/start`, {
        method: "POST",
        body: JSON.stringify({ origin: window.location.origin, account_id: selected?.id ?? "" }),
      });
      window.location.assign(result.url);
    });
  }
  function loadModels() {
    void run(async () => {
      setModels([]);
      setModel("");
      const result = await http<{ models: Model[] }>(`${root}/models`, {
        method: "POST",
        body: JSON.stringify({ account_id: accountID }),
      });
      setModels(result.models);
      if (!result.models.length) setError(t("noModels"));
    });
  }
  function save() {
    void run(async () => {
      await http(`${root}/profile`, {
        method: "POST",
        body: JSON.stringify({ account_id: accountID, name: name.trim(), model }),
      });
      toast.success(t("saved"));
      onSaved();
    });
  }
  function disconnect() {
    void run(async () => {
      const result = await http<{ revoked: boolean }>(`${root}/disconnect`, {
        method: "POST",
        body: JSON.stringify({ account_id: accountID }),
      });
      setModels([]);
      setModel("");
      await reload();
      onSaved();
      if (result.revoked) toast.success(t("disconnected"));
      else setError(t("revocationPending"));
    });
  }

  return (
    <Card>
      <CardContent className="space-y-4 pt-6">
        <div className="space-y-1">
          <h2 className="font-semibold">ChatGPT</h2>
          <p className="text-muted-foreground text-sm">{t("description")}</p>
          <p className="text-muted-foreground text-sm">{t("localOnly")}</p>
        </div>
        <div className="grid gap-2 sm:max-w-lg">
          <Label htmlFor="chatgpt-account">{t("account")}</Label>
          <Select
            value={accountID}
            disabled={busy || !loaded}
            onValueChange={(id) => {
              setAccountID(id);
              setModels([]);
              setModel("");
              setError("");
            }}
          >
            <SelectTrigger id="chatgpt-account">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="new">{t("addAccount")}</SelectItem>
              {accounts.map((a) => (
                <SelectItem key={a.id} value={a.id}>
                  {a.label} · {a.id.slice(0, 8)}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        {selected && <p className="text-sm">{selected.connected ? t("connected") : t("reconnect")}</p>}
        <div className="flex flex-wrap gap-2">
          <Button variant="outline" onClick={connect} disabled={busy || !loaded}>
            {t("continue")}
          </Button>
          {selected?.connected && (
            <Button variant="outline" onClick={loadModels} disabled={busy}>
              {t("loadModels")}
            </Button>
          )}
          {selected && (
            <Button variant="outline" onClick={disconnect} disabled={busy}>
              {t("disconnect")}
            </Button>
          )}
          {!loaded && (
            <Button variant="outline" disabled={busy} onClick={() => void run(reload)}>
              {t("retry")}
            </Button>
          )}
          <Button variant="link" asChild>
            <a href="https://chatgpt.com/settings/usage" target="_blank" rel="noopener noreferrer">
              {t("usage")}
            </a>
          </Button>
        </div>
        {models.length > 0 && (
          <div className="grid gap-3 sm:grid-cols-3">
            <div className="space-y-2">
              <Label htmlFor="chatgpt-name">{t("profileName")}</Label>
              <Input
                id="chatgpt-name"
                value={name}
                maxLength={120}
                disabled={busy}
                onChange={(e) => setName(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="chatgpt-model">{t("model")}</Label>
              <Select value={model} onValueChange={setModel} disabled={busy}>
                <SelectTrigger id="chatgpt-model">
                  <SelectValue placeholder={t("selectModel")} />
                </SelectTrigger>
                <SelectContent>
                  {models.map((m) => (
                    <SelectItem key={m.slug} value={m.slug}>
                      {m.display_name || m.slug}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <Button className="self-end" onClick={save} disabled={busy || !model || !name.trim()}>
              {t("save")}
            </Button>
          </div>
        )}
        {busy && (
          <p role="status" className="text-sm">
            {t("working")}
          </p>
        )}
        {error && (
          <p role="alert" className="text-destructive text-sm">
            {error}
          </p>
        )}
      </CardContent>
    </Card>
  );
}
