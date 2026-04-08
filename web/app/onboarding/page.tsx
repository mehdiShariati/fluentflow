"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { api, ApiError } from "@/lib/api";
import { getToken } from "@/lib/auth";
import { AppNav } from "@/components/AppNav";

type Profile = {
  source_language: string;
  target_language: string;
  proficiency_level: string;
  learning_goal: string;
  tutor_style: string;
};

const levels = ["A1", "A2", "B1", "B2", "C1", "C2"];
const goals = [
  { id: "travel", label: "Travel" },
  { id: "business", label: "Business" },
  { id: "exam", label: "Exam prep" },
  { id: "daily_life", label: "Daily life" },
];
const styles = [
  { id: "gentle", label: "Gentle corrections" },
  { id: "strict", label: "Direct / strict" },
];

export default function OnboardingPage() {
  const router = useRouter();
  const token = getToken();
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [err, setErr] = useState<string | null>(null);
  const [form, setForm] = useState<Profile>({
    source_language: "en",
    target_language: "",
    proficiency_level: "A2",
    learning_goal: "daily_life",
    tutor_style: "gentle",
  });

  useEffect(() => {
    if (!token) {
      router.replace("/login");
      return;
    }
    (async () => {
      try {
        const p = await api<Profile>("/v1/me/profile", { token });
        setForm((f) => ({
          ...f,
          ...p,
          target_language: p.target_language || "",
        }));
      } catch {
        setErr("Could not load profile");
      } finally {
        setLoading(false);
      }
    })();
  }, [token, router]);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!token) return;
    setErr(null);
    setSaving(true);
    try {
      await api("/v1/me/profile", {
        method: "PUT",
        token,
        body: JSON.stringify(form),
      });
      router.replace("/scenarios");
    } catch (ex) {
      setErr(ex instanceof ApiError ? ex.message : "Save failed");
    } finally {
      setSaving(false);
    }
  }

  if (!token) return null;
  if (loading) {
    return (
      <>
        <AppNav />
        <p className="p-8 text-center text-zinc-500">Loading profile…</p>
      </>
    );
  }

  return (
    <>
      <AppNav />
      <main className="mx-auto max-w-lg px-4 py-10 font-[family-name:var(--font-geist-sans)]">
        <h1 className="text-xl font-semibold text-white">Your learning profile</h1>
        <p className="mt-2 text-sm text-zinc-400">
          Used to personalize scenarios and tutor behavior (PRD §12.1).
        </p>
        <form onSubmit={onSubmit} className="mt-8 flex flex-col gap-4">
          <label className="text-sm text-zinc-300">
            I speak (source)
            <input
              className="mt-1 w-full rounded-lg border border-zinc-800 bg-zinc-900 px-3 py-2 text-white outline-none ring-emerald-500/40 focus:ring-2"
              value={form.source_language}
              onChange={(e) =>
                setForm({ ...form, source_language: e.target.value })
              }
              placeholder="en"
            />
          </label>
          <label className="text-sm text-zinc-300">
            I am learning (target) *
            <input
              required
              className="mt-1 w-full rounded-lg border border-zinc-800 bg-zinc-900 px-3 py-2 text-white outline-none ring-emerald-500/40 focus:ring-2"
              value={form.target_language}
              onChange={(e) =>
                setForm({ ...form, target_language: e.target.value })
              }
              placeholder="es, fr, de…"
            />
          </label>
          <label className="text-sm text-zinc-300">
            Level
            <select
              className="mt-1 w-full rounded-lg border border-zinc-800 bg-zinc-900 px-3 py-2 text-white outline-none ring-emerald-500/40 focus:ring-2"
              value={form.proficiency_level}
              onChange={(e) =>
                setForm({ ...form, proficiency_level: e.target.value })
              }
            >
              {levels.map((l) => (
                <option key={l} value={l}>
                  {l}
                </option>
              ))}
            </select>
          </label>
          <label className="text-sm text-zinc-300">
            Goal
            <select
              className="mt-1 w-full rounded-lg border border-zinc-800 bg-zinc-900 px-3 py-2 text-white outline-none ring-emerald-500/40 focus:ring-2"
              value={form.learning_goal}
              onChange={(e) =>
                setForm({ ...form, learning_goal: e.target.value })
              }
            >
              {goals.map((g) => (
                <option key={g.id} value={g.id}>
                  {g.label}
                </option>
              ))}
            </select>
          </label>
          <label className="text-sm text-zinc-300">
            Tutor style
            <select
              className="mt-1 w-full rounded-lg border border-zinc-800 bg-zinc-900 px-3 py-2 text-white outline-none ring-emerald-500/40 focus:ring-2"
              value={form.tutor_style}
              onChange={(e) =>
                setForm({ ...form, tutor_style: e.target.value })
              }
            >
              {styles.map((s) => (
                <option key={s.id} value={s.id}>
                  {s.label}
                </option>
              ))}
            </select>
          </label>
          {err && (
            <p className="text-sm text-red-400" role="alert">
              {err}
            </p>
          )}
          <button
            type="submit"
            disabled={saving}
            className="rounded-lg bg-emerald-600 px-4 py-2.5 text-sm font-medium text-white hover:bg-emerald-500 disabled:opacity-50"
          >
            {saving ? "Saving…" : "Save and continue"}
          </button>
        </form>
      </main>
    </>
  );
}
