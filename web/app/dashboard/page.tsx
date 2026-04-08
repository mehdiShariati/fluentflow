"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { api } from "@/lib/api";
import { getToken } from "@/lib/auth";
import { AppNav } from "@/components/AppNav";

type Summary = {
  total_sessions: number;
  total_speaking_mins: number;
  completed_sessions: number;
  last_session_at: string | null;
  avg_score: number | null;
};

type SessionRow = {
  id: string;
  scenario_id: string;
  scenario_title?: string;
  status: string;
  started_at: string;
};

type LearningSnapshot = {
  id: number;
  captured_at: string;
  total_sessions: number;
  total_speaking_minutes: number;
  avg_session_score: number | null;
};

export default function DashboardPage() {
  const router = useRouter();
  const token = getToken();
  const [data, setData] = useState<Summary | null>(null);
  const [recent, setRecent] = useState<SessionRow[]>([]);
  const [snapshots, setSnapshots] = useState<LearningSnapshot[]>([]);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    if (!token) {
      router.replace("/login");
      return;
    }
    (async () => {
      try {
        const prof = await api<{ target_language: string }>("/v1/me/profile", {
          token,
        });
        if (!prof.target_language) {
          router.replace("/onboarding");
          return;
        }
        const [s, sess, snap] = await Promise.all([
          api<Summary>("/v1/dashboard/summary", { token }),
          api<{ sessions: SessionRow[] }>("/v1/sessions?limit=8", { token }),
          api<{ snapshots: LearningSnapshot[] }>(
            "/v1/me/learning-snapshots?limit=12",
            { token }
          ),
        ]);
        setData(s);
        setRecent(sess.sessions ?? []);
        setSnapshots(snap.snapshots ?? []);
      } catch {
        setErr("Could not load dashboard");
      }
    })();
  }, [token, router]);

  if (!token) return null;

  return (
    <>
      <AppNav />
      <main className="mx-auto max-w-2xl px-4 py-10 font-[family-name:var(--font-geist-sans)]">
        <h1 className="text-xl font-semibold text-white">Progress</h1>
        <p className="mt-2 text-sm text-zinc-400">
          Speaking-first metrics aligned with the PRD dashboard (§12.6).
        </p>
        {err && (
          <p className="mt-4 text-sm text-red-400" role="alert">
            {err}
          </p>
        )}
        {data && (
          <dl className="mt-8 grid grid-cols-2 gap-4 sm:grid-cols-3">
            <div className="rounded-xl border border-zinc-800 bg-zinc-900/40 p-4">
              <dt className="text-xs uppercase tracking-wide text-zinc-500">
                Sessions
              </dt>
              <dd className="mt-1 text-2xl font-semibold text-white">
                {data.total_sessions}
              </dd>
            </div>
            <div className="rounded-xl border border-zinc-800 bg-zinc-900/40 p-4">
              <dt className="text-xs uppercase tracking-wide text-zinc-500">
                Completed
              </dt>
              <dd className="mt-1 text-2xl font-semibold text-white">
                {data.completed_sessions}
              </dd>
            </div>
            <div className="rounded-xl border border-zinc-800 bg-zinc-900/40 p-4">
              <dt className="text-xs uppercase tracking-wide text-zinc-500">
                Speaking min
              </dt>
              <dd className="mt-1 text-2xl font-semibold text-white">
                {data.total_speaking_mins.toFixed(1)}
              </dd>
            </div>
            <div className="rounded-xl border border-zinc-800 bg-zinc-900/40 p-4 sm:col-span-2">
              <dt className="text-xs uppercase tracking-wide text-zinc-500">
                Avg session score
              </dt>
              <dd className="mt-1 text-2xl font-semibold text-white">
                {data.avg_score != null ? data.avg_score.toFixed(1) : "—"}
              </dd>
            </div>
            <div className="rounded-xl border border-zinc-800 bg-zinc-900/40 p-4 sm:col-span-3">
              <dt className="text-xs uppercase tracking-wide text-zinc-500">
                Last session
              </dt>
              <dd className="mt-1 text-sm text-zinc-300">
                {data.last_session_at
                  ? new Date(data.last_session_at).toLocaleString()
                  : "—"}
              </dd>
            </div>
          </dl>
        )}
        {snapshots.length > 0 && (
          <section className="mt-10">
            <h2 className="text-sm font-medium uppercase tracking-wide text-zinc-500">
              Progress snapshots
            </h2>
            <p className="mt-1 text-xs text-zinc-600">
              Captured when you complete a session (rollup for history / analytics).
            </p>
            <ul className="mt-3 max-h-48 space-y-2 overflow-y-auto text-sm">
              {snapshots.map((row) => (
                <li
                  key={row.id}
                  className="flex flex-wrap justify-between gap-2 rounded-lg border border-zinc-800 bg-zinc-900/30 px-3 py-2 text-zinc-300"
                >
                  <time
                    className="text-zinc-500"
                    dateTime={row.captured_at}
                  >
                    {new Date(row.captured_at).toLocaleString()}
                  </time>
                  <span>
                    {row.total_sessions} sess ·{" "}
                    {Number(row.total_speaking_minutes).toFixed(1)} min
                    {row.avg_session_score != null
                      ? ` · avg ${row.avg_session_score.toFixed(1)}`
                      : ""}
                  </span>
                </li>
              ))}
            </ul>
          </section>
        )}
        {recent.length > 0 && (
          <section className="mt-10">
            <h2 className="text-sm font-medium uppercase tracking-wide text-zinc-500">
              Recent sessions
            </h2>
            <ul className="mt-3 space-y-2">
              {recent.map((row) => (
                <li
                  key={row.id}
                  className="flex flex-wrap items-center justify-between gap-2 rounded-lg border border-zinc-800 bg-zinc-900/30 px-3 py-2 text-sm"
                >
                  <div>
                    <span className="font-medium text-zinc-200">
                      {row.scenario_title || row.scenario_id}
                    </span>
                    {row.scenario_title && (
                      <span className="ml-2 font-mono text-xs text-zinc-600">
                        {row.scenario_id}
                      </span>
                    )}
                    <span className="ml-2 text-zinc-500">{row.status}</span>
                  </div>
                  <div className="flex items-center gap-3 text-zinc-500">
                    <time dateTime={row.started_at}>
                      {new Date(row.started_at).toLocaleString()}
                    </time>
                    {row.status === "active" ? (
                      <Link
                        href={`/session/${row.id}`}
                        className="text-emerald-400 hover:underline"
                      >
                        Resume
                      </Link>
                    ) : (
                      <Link
                        href={`/feedback/${row.id}`}
                        className="text-emerald-400 hover:underline"
                      >
                        Feedback
                      </Link>
                    )}
                  </div>
                </li>
              ))}
            </ul>
          </section>
        )}
        <div className="mt-10 flex flex-wrap gap-3">
          <Link
            href="/scenarios"
            className="rounded-lg bg-emerald-600 px-4 py-2.5 text-sm font-medium text-white hover:bg-emerald-500"
          >
            Start a scenario
          </Link>
          <Link
            href="/onboarding"
            className="rounded-lg border border-zinc-700 px-4 py-2.5 text-sm text-zinc-200 hover:bg-zinc-900"
          >
            Edit profile
          </Link>
          <Link
            href="/settings"
            className="rounded-lg border border-zinc-700 px-4 py-2.5 text-sm text-zinc-200 hover:bg-zinc-900"
          >
            Account
          </Link>
        </div>
      </main>
    </>
  );
}
