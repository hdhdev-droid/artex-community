"use client";

import * as React from "react";

import Link from "next/link";

import { useLocale, useTranslations } from "next-intl";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { api } from "@/lib/api";
import { useGeneratedMessage } from "@/lib/generated-message-i18n";
import {
  type ApiInputMode,
  type AssessmentCatalog,
  type AssessmentKind,
  type AssessmentRun,
  type CandidateReview,
  offensiveApi,
} from "@/lib/offensive-api";
import type { LLMProfile } from "@/lib/types";

const selectStyle = "h-9 w-full rounded-md border bg-background px-3 text-sm";
const active = (status?: string) => status === "queued" || status === "running";
const textError = (error: unknown, fallback: string) => (error instanceof Error ? error.message : fallback);

export default function OffensiveAssessmentsPage() {
  const t = useTranslations("offensiveAssessments");
  const generated = useGeneratedMessage();
  const [catalog, setCatalog] = React.useState<AssessmentCatalog | null>(null);
  const [runs, setRuns] = React.useState<AssessmentRun[]>([]);
  const [profiles, setProfiles] = React.useState<LLMProfile[]>([]);
  const [ready, setReady] = React.useState(false);
  const [connected, setConnected] = React.useState(false);
  const [kind, setKind] = React.useState<AssessmentKind>("api");
  const [mode, setMode] = React.useState<ApiInputMode>("openapi");
  const [targetID, setTargetID] = React.useState("");
  const [operationIDs, setOperationIDs] = React.useState<string[]>([]);
  const [name, setName] = React.useState("");
  const [baseURL, setBaseURL] = React.useState("");
  const [harOrigin, setHarOrigin] = React.useState("");
  const [file, setFile] = React.useState<File | null>(null);
  const [token, setToken] = React.useState("");
  const [selectedID, setSelectedID] = React.useState("");
  const [busy, setBusy] = React.useState(false);
  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState("");
  const [revision, setRevision] = React.useState(0);
  const mounted = React.useRef(false);
  const fileInput = React.useRef<HTMLInputElement>(null);
  const formID = React.useId();
  React.useEffect(() => {
    mounted.current = true;
    return () => {
      mounted.current = false;
    };
  }, []);

  // biome-ignore lint/correctness/useExhaustiveDependencies: Explicit refresh changes revision to restart polling immediately.
  React.useEffect(() => {
    let cancelled = false;
    let timer: ReturnType<typeof setTimeout> | undefined;
    async function refresh() {
      const results = await Promise.allSettled([
        offensiveApi.capabilities(),
        offensiveApi.catalog(),
        offensiveApi.list(),
        api.llmProfiles(),
      ]);
      if (cancelled) return;
      const [capabilities, targets, history, models] = results;
      setReady(capabilities.status === "fulfilled" && capabilities.value.review_storage_ready);
      setConnected(capabilities.status === "fulfilled" && capabilities.value.runner_configured);
      if (targets.status === "fulfilled") setCatalog(targets.value);
      if (history.status === "fulfilled") setRuns(history.value.runs ?? []);
      if (models.status === "fulfilled") setProfiles(models.value);
      const failed = results.find((result) => result.status === "rejected");
      setError(failed?.status === "rejected" ? textError(failed.reason, t("requestError")) : "");
      setLoading(false);
      if (failed || (history.status === "fulfilled" && (history.value.runs ?? []).some((run) => active(run.status)))) {
        timer = setTimeout(refresh, 5000);
      }
    }
    void refresh();
    return () => {
      cancelled = true;
      clearTimeout(timer);
    };
  }, [revision, t]);

  const targets = (catalog?.targets ?? []).filter((item) => item.kind === kind);
  const target = targets.find((item) => item.id === targetID);
  const operations = (target?.operations ?? []).filter((item) => ["GET", "HEAD"].includes(item.method.toUpperCase()));
  const maxMiB = kind === "source" ? 20 : mode === "har" ? 16 : 8;
  const needsFile = kind === "source" || mode !== "url";
  function resetTarget() {
    setTargetID("");
    setOperationIDs([]);
    setToken("");
  }
  async function upload(event: React.FormEvent) {
    event.preventDefault();
    if (busy) return;
    if (needsFile && (!file || file.size > maxMiB * 1024 * 1024)) {
      setError(t("fileLimit", { size: maxMiB }));
      return;
    }
    setBusy(true);
    setError("");
    try {
      const added = await offensiveApi.uploadTarget({
        kind,
        name: name.trim(),
        base_url: kind === "api" ? baseURL.trim() : undefined,
        input_mode: kind === "api" ? mode : undefined,
        har_origin: kind === "api" && mode === "har" ? harOrigin.trim() : undefined,
        file: needsFile ? (file ?? undefined) : undefined,
      });
      if (!mounted.current) return;
      setCatalog((previous) => ({
        ...previous,
        tools: previous?.tools ?? [],
        targets: [...(previous?.targets ?? []).filter((item) => item.id !== added.id), added],
      }));
      setTargetID(added.id);
      setOperationIDs([]);
      setToken("");
      setFile(null);
      if (fileInput.current) fileInput.current.value = "";
      setRevision((value) => value + 1);
    } catch (reason) {
      if (mounted.current) setError(textError(reason, t("requestError")));
    } finally {
      if (mounted.current) setBusy(false);
    }
  }
  async function start() {
    if (!target || busy) return;
    setBusy(true);
    setError("");
    try {
      const run = await offensiveApi.create(
        kind,
        target.id,
        kind === "api"
          ? {
              operation_ids: operationIDs,
              bearer_token: token || undefined,
            }
          : undefined,
      );
      if (!mounted.current) return;
      setSelectedID(run.id);
      setRevision((value) => value + 1);
    } catch (reason) {
      if (mounted.current) setError(textError(reason, t("requestError")));
    } finally {
      if (mounted.current) {
        setToken("");
        setBusy(false);
      }
    }
  }
  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="text-2xl font-semibold">{t("title")}</h1>
          <p className="text-sm text-muted-foreground">{t("description")}</p>
        </div>
        <Button variant="outline" disabled={busy || loading} onClick={() => setRevision((value) => value + 1)}>
          {t("refresh")}
        </Button>
      </div>
      {error && (
        <p role="alert" className="text-sm text-destructive">
          {generated(error)}
        </p>
      )}
      {!loading && !connected && (
        <p role="status" className="rounded-md border p-4 text-sm">
          {t("runnerMissing")}
        </p>
      )}
      {!loading && !ready && (
        <p role="status" className="rounded-md border p-4 text-sm">
          {t("storageMissing")}
        </p>
      )}
      <Card>
        <CardHeader>
          <CardTitle>{t("register")}</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={upload} className="space-y-4">
            <fieldset disabled={busy || !connected} className="space-y-4">
              <div className="grid gap-4 md:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor={`${formID}-kind`}>{t("kind")}</Label>
                  <select
                    id={`${formID}-kind`}
                    className={selectStyle}
                    value={kind}
                    onChange={(event) => {
                      setKind(event.target.value as AssessmentKind);
                      resetTarget();
                      setFile(null);
                      if (fileInput.current) fileInput.current.value = "";
                    }}
                  >
                    <option value="api">{t("api")}</option>
                    <option value="source">{t("source")}</option>
                  </select>
                </div>
                <div className="space-y-2">
                  <Label htmlFor={`${formID}-name`}>{t("name")}</Label>
                  <Input
                    id={`${formID}-name`}
                    value={name}
                    required
                    maxLength={120}
                    onChange={(event) => setName(event.target.value)}
                  />
                </div>
              </div>
              {kind === "api" && (
                <>
                  <div className="space-y-2">
                    <Label htmlFor={`${formID}-mode`}>{t("inputMode")}</Label>
                    <select
                      id={`${formID}-mode`}
                      className={selectStyle}
                      value={mode}
                      onChange={(event) => {
                        setMode(event.target.value as ApiInputMode);
                        setFile(null);
                        if (fileInput.current) fileInput.current.value = "";
                      }}
                    >
                      <option value="openapi">{t("openapi")}</option>
                      <option value="url">{t("url")}</option>
                      <option value="har">{t("har")}</option>
                    </select>
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor={`${formID}-url`}>{t("baseURL")}</Label>
                    <Input
                      id={`${formID}-url`}
                      type="url"
                      required
                      value={baseURL}
                      onChange={(event) => setBaseURL(event.target.value)}
                    />
                    <p className="text-xs text-muted-foreground">{t("addressHelp")}</p>
                  </div>
                  {mode === "har" && (
                    <div className="space-y-2">
                      <Label htmlFor={`${formID}-origin`}>{t("harOrigin")}</Label>
                      <Input
                        id={`${formID}-origin`}
                        type="url"
                        value={harOrigin}
                        onChange={(event) => setHarOrigin(event.target.value)}
                      />
                    </div>
                  )}
                </>
              )}
              {needsFile && (
                <div className="space-y-2">
                  <Label htmlFor={`${formID}-file`}>{t("fileLimit", { size: maxMiB })}</Label>
                  <Input
                    ref={fileInput}
                    id={`${formID}-file`}
                    type="file"
                    required
                    accept={kind === "source" ? ".zip" : mode === "har" ? ".har,.json" : ".json,.yaml,.yml"}
                    onChange={(event) => setFile(event.target.files?.[0] ?? null)}
                  />
                </div>
              )}
              <Button type="submit" disabled={busy || !connected}>
                {busy ? t("working") : t("register")}
              </Button>
            </fieldset>
          </form>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle>{t("diagnose")}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <Label htmlFor={`${formID}-target`}>{t("target")}</Label>
          <select
            id={`${formID}-target`}
            className={selectStyle}
            value={targetID}
            disabled={busy || !connected}
            onChange={(event) => {
              resetTarget();
              setTargetID(event.target.value);
            }}
          >
            <option value="">{t("selectTarget")}</option>
            {targets.map((item) => (
              <option key={item.id} value={item.id} disabled={!item.ready}>
                {item.name}
              </option>
            ))}
          </select>
          {target?.detail && <p className="text-sm text-muted-foreground">{generated(target.detail)}</p>}
          {target?.discovery?.summary && <p className="text-sm">{generated(target.discovery.summary)}</p>}
          {target?.discovery?.warnings?.map((warning, index) => (
            <p key={`${index}-${warning}`} className="text-sm text-muted-foreground">
              {generated(warning)}
            </p>
          ))}
          {kind === "api" && (
            <>
              <p className="text-sm text-muted-foreground">{t("operationHelp")}</p>
              <div className="max-h-64 space-y-2 overflow-auto">
                {operations.map((operation) => (
                  <label key={operation.id} className="flex items-center gap-2 text-sm">
                    <input
                      type="checkbox"
                      checked={operationIDs.includes(operation.id)}
                      disabled={busy || (!operationIDs.includes(operation.id) && operationIDs.length >= 40)}
                      onChange={(event) =>
                        setOperationIDs((previous) =>
                          event.target.checked
                            ? [...previous, operation.id]
                            : previous.filter((id) => id !== operation.id),
                        )
                      }
                    />
                    <span className="break-all font-mono">
                      {operation.method} {operation.path}
                    </span>
                  </label>
                ))}
              </div>
              <Label htmlFor={`${formID}-token`}>{t("bearer")}</Label>
              <Input
                id={`${formID}-token`}
                type="password"
                autoComplete="off"
                disabled={busy}
                value={token}
                onChange={(event) => setToken(event.target.value)}
              />
              <p className="text-xs text-muted-foreground">{t("tokenHelp")}</p>
            </>
          )}
          <Button
            onClick={() => void start()}
            disabled={busy || !connected || !target?.ready || (kind === "api" && operationIDs.length === 0)}
          >
            {busy ? t("working") : t("start")}
          </Button>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle>{t("history")}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-2">
          {loading ? (
            <p role="status">{t("loading")}</p>
          ) : (
            runs.length === 0 && <p className="text-sm text-muted-foreground">{t("empty")}</p>
          )}
          {runs.map((run) => (
            <Button
              key={run.id}
              variant={selectedID === run.id ? "secondary" : "outline"}
              className="h-auto w-full justify-start whitespace-normal text-left"
              onClick={() => setSelectedID(run.id)}
            >
              {run.target_snapshot?.name ?? run.target_id} · {t(`status.${run.status}`)} · {run.id}
            </Button>
          ))}
        </CardContent>
      </Card>
      {selectedID && <RunDetail key={selectedID} id={selectedID} profiles={profiles} storageReady={ready} />}
    </div>
  );
}

function RunDetail({ id, profiles, storageReady }: { id: string; profiles: LLMProfile[]; storageReady: boolean }) {
  const t = useTranslations("offensiveAssessments");
  const generated = useGeneratedMessage();
  const locale = useLocale();
  const [run, setRun] = React.useState<AssessmentRun | null>(null);
  const [review, setReview] = React.useState<CandidateReview | null>(null);
  const [profileID, setProfileID] = React.useState("");
  const [batch, setBatch] = React.useState(3);
  const [busy, setBusy] = React.useState(false);
  const [error, setError] = React.useState("");
  const [reviewError, setReviewError] = React.useState("");
  const [revision, setRevision] = React.useState(0);
  const mounted = React.useRef(false);
  const controlsID = React.useId();
  React.useEffect(() => {
    mounted.current = true;
    return () => {
      mounted.current = false;
    };
  }, []);
  // biome-ignore lint/correctness/useExhaustiveDependencies: Review actions change revision to restart polling immediately.
  React.useEffect(() => {
    let cancelled = false;
    let timer: ReturnType<typeof setTimeout> | undefined;
    async function refresh() {
      const results = await Promise.allSettled([
        offensiveApi.get(id),
        storageReady ? offensiveApi.candidateReview(id) : Promise.resolve({ review: null }),
      ]);
      if (cancelled) return;
      const [scan, checkpoint] = results;
      if (scan.status === "fulfilled") {
        setRun(scan.value);
        setError("");
      } else setError(textError(scan.reason, t("requestError")));
      if (checkpoint.status === "fulfilled") {
        setReview(checkpoint.value.review);
        setReviewError("");
      } else setReviewError(textError(checkpoint.reason, t("requestError")));
      // Poll only while either operation is active, and retry transient read failures.
      if (
        scan.status === "rejected" ||
        checkpoint.status === "rejected" ||
        active(scan.status === "fulfilled" ? scan.value.status : undefined) ||
        active(checkpoint.status === "fulfilled" ? checkpoint.value.review?.status : undefined)
      ) {
        timer = setTimeout(refresh, 3000);
      }
    }
    void refresh();
    return () => {
      cancelled = true;
      clearTimeout(timer);
    };
  }, [id, storageReady, revision, t]);

  async function reviewAction(cancel: boolean) {
    if (busy || !storageReady) return;
    setBusy(true);
    setReviewError("");
    try {
      const result = cancel
        ? await offensiveApi.cancelCandidateReview(id)
        : await offensiveApi.startCandidateReview(id, {
            llm_profile_id: profileID ? Number(profileID) : undefined,
            batch_size: batch,
          });
      if (!mounted.current) return;
      setReview(result.review);
      setRevision((value) => value + 1);
    } catch (reason) {
      if (mounted.current) setReviewError(textError(reason, t("requestError")));
    } finally {
      if (mounted.current) setBusy(false);
    }
  }
  async function download(name: string) {
    if (busy) return;
    setBusy(true);
    setError("");
    try {
      const blob = await offensiveApi.artifact(id, name);
      if (!mounted.current) return;
      const url = URL.createObjectURL(blob);
      const anchor = document.createElement("a");
      anchor.href = url;
      anchor.download = name;
      anchor.click();
      setTimeout(() => URL.revokeObjectURL(url), 1000);
    } catch (reason) {
      if (mounted.current) setError(textError(reason, t("requestError")));
    } finally {
      if (mounted.current) setBusy(false);
    }
  }
  const time = (value: string) =>
    new Intl.DateTimeFormat(locale, { dateStyle: "medium", timeStyle: "short" }).format(new Date(value));
  return (
    <Card>
      <CardHeader>
        <CardTitle>{t("results")}</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {error && (
          <p role="alert" className="text-sm text-destructive">
            {generated(error)}
          </p>
        )}
        <Button variant="outline" disabled={busy} onClick={() => setRevision((value) => value + 1)}>
          {t("refresh")}
        </Button>
        {!run ? (
          <p role="status">{t("loading")}</p>
        ) : (
          <>
            <Badge variant="secondary">{t(`status.${run.status}`)}</Badge>
            <p className="text-sm">{generated(run.summary ?? "")}</p>
            {run.error && (
              <p role="alert" className="text-sm text-destructive">
                {generated(run.error)}
              </p>
            )}
            <p className="text-sm text-muted-foreground">{t("candidateNotice")}</p>
            {(run.limitations ?? []).map((limitation, index) => (
              <p key={index} className="text-sm text-muted-foreground">
                {generated(limitation)}
              </p>
            ))}
            {(run.candidates ?? []).map((candidate) => (
              <details key={candidate.id} className="rounded-md border p-3">
                <summary className="cursor-pointer text-sm">
                  {candidate.title} · {candidate.severity}
                </summary>
                <p className="mt-2 whitespace-pre-wrap text-sm">{candidate.message}</p>
                {candidate.file && (
                  <p className="break-all font-mono text-xs">
                    {candidate.file}:{candidate.line}
                  </p>
                )}
                {candidate.evidence != null && (
                  <pre className="mt-2 max-h-64 overflow-auto whitespace-pre-wrap break-all text-xs">
                    {typeof candidate.evidence === "string"
                      ? candidate.evidence
                      : JSON.stringify(candidate.evidence, null, 2)}
                  </pre>
                )}
              </details>
            ))}
            <div className="flex flex-wrap gap-2">
              {(run.artifacts ?? []).map((artifact) => (
                <Button
                  key={artifact.name}
                  variant="outline"
                  disabled={busy}
                  onClick={() => void download(artifact.name)}
                >
                  {t("download")} · {artifact.name}
                </Button>
              ))}
            </div>
            <h2 className="text-lg font-semibold">{t("candidateReview")}</h2>
            <p className="text-sm text-muted-foreground">{t("reviewHelp")}</p>
            {!storageReady && (
              <p role="status" className="text-sm">
                {t("storageMissing")}
              </p>
            )}
            {reviewError && (
              <p role="alert" className="text-sm text-destructive">
                {reviewError}
              </p>
            )}
            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor={`${controlsID}-model`}>{t("model")}</Label>
                <select
                  id={`${controlsID}-model`}
                  className={selectStyle}
                  value={profileID}
                  disabled={busy || active(review?.status) || !storageReady}
                  onChange={(event) => setProfileID(event.target.value)}
                >
                  <option value="">{t("defaultModel")}</option>
                  {profiles.map((profile) => (
                    <option key={profile.id} value={profile.id}>
                      {profile.name}
                    </option>
                  ))}
                </select>
              </div>
              <div className="space-y-2">
                <Label htmlFor={`${controlsID}-batch`}>{t("batch")}</Label>
                <select
                  id={`${controlsID}-batch`}
                  className={selectStyle}
                  value={batch}
                  disabled={busy || active(review?.status) || !storageReady}
                  onChange={(event) => setBatch(Number(event.target.value))}
                >
                  {[1, 2, 3, 4, 5].map((value) => (
                    <option key={value} value={value}>
                      {value}
                    </option>
                  ))}
                </select>
              </div>
            </div>
            <div className="flex gap-2">
              <Button
                disabled={
                  busy ||
                  !!reviewError ||
                  !storageReady ||
                  run.status !== "completed" ||
                  active(review?.status) ||
                  review?.pending === 0
                }
                onClick={() => void reviewAction(false)}
              >
                {review ? t("resume") : t("startReview")}
              </Button>
              {active(review?.status) && (
                <Button variant="outline" disabled={busy} onClick={() => void reviewAction(true)}>
                  {t("stop")}
                </Button>
              )}
            </div>
            {review && (
              <>
                <p className="text-sm">
                  {t(`status.${review.status}`)} · {t("progress", { done: review.reviewed, total: review.total })}
                </p>
                <p className="text-sm">{review.reason}</p>
                {review.candidates.map((candidate) => (
                  <details key={candidate.id} className="rounded-md border p-3">
                    <summary className="cursor-pointer text-sm">
                      {t(`status.${candidate.status}`)} · {candidate.title}
                    </summary>
                    {candidate.result && (
                      <div className="mt-2 space-y-2 text-sm">
                        <p>{candidate.result.summary}</p>
                        <p className="whitespace-pre-wrap">{candidate.result.rationale}</p>
                        <p className="break-all font-mono text-xs">{candidate.result.evidence_refs?.join(", ")}</p>
                        {candidate.result.limitations?.map((item, index) => (
                          <p key={index} className="text-muted-foreground">
                            {item}
                          </p>
                        ))}
                      </div>
                    )}
                    {candidate.last_error && <p className="text-sm text-destructive">{candidate.last_error}</p>}
                  </details>
                ))}
                <h3 className="font-medium">{t("attempts")}</h3>
                {review.attempts.map((attempt) => (
                  <div key={attempt.id} className="rounded-md border p-3 text-sm">
                    <p>
                      {time(attempt.queued_at)} · {t(`status.${attempt.status}`)}
                    </p>
                    <p>
                      {t("metrics", {
                        calls: attempt.model_calls,
                        input: attempt.input_tokens,
                        output: attempt.output_tokens,
                        model: Math.round(attempt.model_elapsed_ms / 1000),
                        queue: Math.round(attempt.queue_wait_ms / 1000),
                      })}
                    </p>
                    <p>{attempt.reason}</p>
                  </div>
                ))}
              </>
            )}
            {(review?.legacy_task_id || run.task_id) && (
              <Link
                className="text-sm underline"
                href={`/function/tasks/detail?id=${encodeURIComponent(review?.legacy_task_id || run.task_id || "")}`}
              >
                {t("legacyTask")}
              </Link>
            )}
          </>
        )}
      </CardContent>
    </Card>
  );
}
