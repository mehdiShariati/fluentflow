"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { api, ApiError } from "@/lib/api";
import { getToken, liveKitStorageKey } from "@/lib/auth";
import { AppNav } from "@/components/AppNav";

type Scenario = {
  id: string;
  title: string;
  description: string;
  level: string;
};

type CreateSessionResp = {
  session_id: string;
  livekit_token?: string;
  note?: string;
};

export default function ScenariosView({
  highlight,
}: {
  highlight?: string;
}) {
  const router = useRouter();
  const token = getToken();
  const [scenarios, setScenarios] = useState<Scenario[]>([]);
  const [assignments, setAssignments] = useState<Record<string, string>>({});
  const [err, setErr] = useState<string | null>(null);
  const [starting, setStarting] = useState<string | null>(null);

  useEffect(() => {
    if (!token) {
      router.replace("/login");
      return;
    }
    (async () => {
      try {
        const [sc, ex] = await Promise.all([
          api<{ scenarios: Scenario[] }>("/v1/scenarios", { token }),
          api<{ assignments: Record<string, string> }>("/v1/experiments", {
            token,
          }),
        ]);
        setScenarios(sc.scenarios);
        setAssignments(ex.assignments);
      } catch {
        setErr("Could not load scenarios");
      }
    })();
  }, [token, router]);

  useEffect(() => {
    if (!highlight || scenarios.length === 0) return;
    const t = window.setTimeout(() => {
      document
        .getElementById(`scenario-${highlight}`)
        ?.scrollIntoView({ behavior: "smooth", block: "center" });
    }, 100);
    return () => window.clearTimeout(t);
  }, [highlight, scenarios]);

  async function start(id: string) {
    if (!token) return;
    setStarting(id);
    setErr(null);
    try {
      const data = await api<CreateSessionResp>("/v1/sessions", {
        method: "POST",
        token,
        body: JSON.stringify({ scenario_id: id }),
      });
      if (data.livekit_token) {
        sessionStorage.setItem(
          liveKitStorageKey(data.session_id),
          data.livekit_token
        );
      }
      await api(`/v1/sessions/${data.session_id}/events`, {
        method: "POST",
        token,
        body: JSON.stringify({
          events: [{ type: "session_joined", payload: { scenario_id: id } }],
        }),
      });
      router.push(`/session/${data.session_id}`);
    } catch (ex) {
      setErr(ex instanceof ApiError ? ex.message : "Could not start session");
    } finally {
      setStarting(null);
    }
  }

  if (!token) return null;

  return (
    <>
      <AppNav />
      <main className="mx-auto max-w-2xl px-4 py-10 font-[family-name:var(--font-geist-sans)]">
        <h1 className="text-xl font-semibold text-white">Scenarios</h1>
        <p className="mt-2 text-sm text-zinc-400">
          Pick a situation to practice. Experiments:{" "}
          <span className="font-mono text-xs text-zinc-500">
            {Object.entries(assignments)
              .map(([k, v]) => `${k}=${v}`)
              .join(", ") || "—"}
          </span>
        </p>
        {err && (
          <p className="mt-4 text-sm text-red-400" role="alert">
            {err}
          </p>
        )}
        <ul className="mt-8 flex flex-col gap-4">
          {scenarios.map((s) => (
            <li
              key={s.id}
              id={`scenario-${s.id}`}
              className={`rounded-xl border bg-zinc-900/50 p-4 ${
                highlight === s.id
                  ? "border-emerald-500/60 ring-2 ring-emerald-500/30"
                  : "border-zinc-800"
              }`}
            >
              <div className="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
                <div>
                  <h2 className="font-medium text-white">{s.title}</h2>
                  <p className="mt-1 text-sm text-zinc-400">{s.description}</p>
                  <p className="mt-2 text-xs text-zinc-500">Level {s.level}</p>
                </div>
                <button
                  type="button"
                  disabled={starting !== null}
                  onClick={() => start(s.id)}
                  className="shrink-0 rounded-lg bg-emerald-600 px-4 py-2 text-sm font-medium text-white hover:bg-emerald-500 disabled:opacity-50"
                >
                  {starting === s.id ? "Starting…" : "Start session"}
                </button>
              </div>
            </li>
          ))}
        </ul>
        <p className="mt-8 text-center text-sm text-zinc-500">
          <Link href="/dashboard" className="text-emerald-400 hover:underline">
            Back to dashboard
          </Link>
        </p>
      </main>
    </>
  );
}
