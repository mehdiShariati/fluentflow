"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { api, ApiError } from "@/lib/api";
import { getToken } from "@/lib/auth";
import { AppNav } from "@/components/AppNav";

type Feedback = {
  strengths: string[];
  top_mistakes: string[];
  suggestions: string[];
  recommended_scenario: string | null;
  recommended_scenario_title?: string;
  score: number | null;
  raw_notes: string | null;
  transcript_summary?: string | null;
  generation_source?: string | null;
};

export default function FeedbackPage() {
  const params = useParams();
  const router = useRouter();
  const id = typeof params?.id === "string" ? params.id : "";
  const token = getToken();
  const [fb, setFb] = useState<Feedback | null>(null);
  const [err, setErr] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [recBusy, setRecBusy] = useState(false);

  async function trackRecommendationAndGo(scenarioId: string) {
    if (!token || !id) return;
    setRecBusy(true);
    try {
      await api(`/v1/sessions/${id}/recommendation-click`, {
        method: "POST",
        token,
        body: JSON.stringify({
          scenario_id: scenarioId,
          source: "feedback_ui",
        }),
      });
      router.push(
        `/scenarios?highlight=${encodeURIComponent(scenarioId)}`
      );
    } catch {
      router.push(
        `/scenarios?highlight=${encodeURIComponent(scenarioId)}`
      );
    } finally {
      setRecBusy(false);
    }
  }

  useEffect(() => {
    if (!token) {
      router.replace("/login");
      return;
    }
    if (!id) {
      setLoading(false);
      return;
    }
    (async () => {
      try {
        await api(`/v1/sessions/${id}/feedback/generate`, {
          method: "POST",
          token,
        });
        const data = await api<Feedback>(`/v1/sessions/${id}/feedback`, {
          token,
        });
        setFb(data);
        await api(`/v1/sessions/${id}/feedback/viewed`, {
          method: "POST",
          token,
          body: JSON.stringify({}),
        });
      } catch (ex) {
        if (ex instanceof ApiError && ex.message === "complete_session_first") {
          setErr("Finish the session first, then open feedback again.");
        } else {
          setErr(ex instanceof ApiError ? ex.message : "Could not load feedback");
        }
      } finally {
        setLoading(false);
      }
    })();
  }, [token, id, router]);

  if (!token || !id) return null;

  return (
    <>
      <AppNav />
      <main className="mx-auto max-w-xl px-4 py-10 font-[family-name:var(--font-geist-sans)]">
        <h1 className="text-xl font-semibold text-white">Session feedback</h1>
        <p className="mt-2 text-sm text-zinc-400">
          {fb
            ? fb.generation_source === "openai"
              ? "Generated with OpenAI (gpt-4o-mini) from your scenario and transcript."
              : fb.generation_source === "stub"
                ? "Deterministic stub on the API (set OPENAI_API_KEY for LLM coaching)."
                : "Post-session summary from the API."
            : "Structured coaching after your session (OpenAI when the API key is set, otherwise stub)."}
        </p>
        {loading && (
          <p className="mt-8 text-zinc-500">Generating feedback…</p>
        )}
        {err && (
          <p className="mt-8 text-sm text-red-400" role="alert">
            {err}
          </p>
        )}
        {fb && (
          <div className="mt-8 space-y-6">
            {fb.transcript_summary && (
              <section>
                <h2 className="text-sm font-medium uppercase tracking-wide text-zinc-500">
                  Session recap
                </h2>
                <p className="mt-2 text-sm leading-relaxed text-zinc-300">
                  {fb.transcript_summary}
                </p>
              </section>
            )}
            {fb.score != null && (
              <p className="text-3xl font-semibold text-emerald-400">
                Score: {fb.score.toFixed(1)} / 10
              </p>
            )}
            <section>
              <h2 className="text-sm font-medium uppercase tracking-wide text-zinc-500">
                What went well
              </h2>
              <ul className="mt-2 list-inside list-disc text-sm text-zinc-300">
                {fb.strengths.map((s) => (
                  <li key={s}>{s}</li>
                ))}
              </ul>
            </section>
            <section>
              <h2 className="text-sm font-medium uppercase tracking-wide text-zinc-500">
                Top fixes
              </h2>
              <ul className="mt-2 list-inside list-disc text-sm text-zinc-300">
                {fb.top_mistakes.map((s) => (
                  <li key={s}>{s}</li>
                ))}
              </ul>
            </section>
            <section>
              <h2 className="text-sm font-medium uppercase tracking-wide text-zinc-500">
                Next practice
              </h2>
              <ul className="mt-2 list-inside list-disc text-sm text-zinc-300">
                {fb.suggestions.map((s) => (
                  <li key={s}>{s}</li>
                ))}
              </ul>
            </section>
            {fb.recommended_scenario && (
              <div className="rounded-lg border border-zinc-800 bg-zinc-900/40 p-4">
                <p className="text-sm text-zinc-400">
                  Suggested next scenario:{" "}
                  <span className="font-medium text-emerald-400">
                    {fb.recommended_scenario_title || fb.recommended_scenario}
                  </span>
                  {fb.recommended_scenario_title && (
                    <span className="ml-2 font-mono text-xs text-zinc-500">
                      ({fb.recommended_scenario})
                    </span>
                  )}
                </p>
                <button
                  type="button"
                  disabled={recBusy}
                  onClick={() =>
                    void trackRecommendationAndGo(fb.recommended_scenario!)
                  }
                  className="mt-3 rounded-lg bg-emerald-700 px-4 py-2 text-sm font-medium text-white hover:bg-emerald-600 disabled:opacity-50"
                >
                  {recBusy ? "Saving…" : "Start recommended scenario"}
                </button>
                <p className="mt-2 text-xs text-zinc-600">
                  Records a product analytics event before opening scenarios.
                </p>
              </div>
            )}
            {fb.raw_notes && (
              <p className="text-xs text-zinc-600">{fb.raw_notes}</p>
            )}
          </div>
        )}
        <div className="mt-10 flex flex-wrap gap-3">
          <Link
            href={
              fb?.recommended_scenario
                ? `/scenarios?highlight=${encodeURIComponent(fb.recommended_scenario)}`
                : "/scenarios"
            }
            className="rounded-lg bg-emerald-600 px-4 py-2 text-sm font-medium text-white hover:bg-emerald-500"
          >
            Another session
          </Link>
          <Link
            href="/dashboard"
            className="rounded-lg border border-zinc-700 px-4 py-2 text-sm text-zinc-200 hover:bg-zinc-900"
          >
            Dashboard
          </Link>
        </div>
      </main>
    </>
  );
}
