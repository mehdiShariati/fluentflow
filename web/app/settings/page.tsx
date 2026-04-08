"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { api, ApiError } from "@/lib/api";
import { clearToken, getToken } from "@/lib/auth";
import { AppNav } from "@/components/AppNav";

type Me = {
  user_id: string;
  email: string;
  is_guest: boolean;
};

export default function SettingsPage() {
  const router = useRouter();
  const token = getToken();
  const [me, setMe] = useState<Me | null>(null);
  const [password, setPassword] = useState("");
  const [err, setErr] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);

  useEffect(() => {
    if (!token) {
      router.replace("/login");
      return;
    }
    (async () => {
      try {
        const u = await api<Me>("/v1/me", { token });
        setMe(u);
      } catch {
        setErr("Could not load account");
      }
    })();
  }, [token, router]);

  async function deleteAccount() {
    if (!token || !me) return;
    setBusy(true);
    setErr(null);
    try {
      await api("/v1/me/account", {
        method: "DELETE",
        token,
        body: JSON.stringify({ password }),
      });
      clearToken();
      router.replace("/login");
    } catch (ex) {
      setErr(
        ex instanceof ApiError ? ex.message : "Could not delete account"
      );
    } finally {
      setBusy(false);
    }
  }

  if (!token) return null;

  return (
    <>
      <AppNav />
      <main className="mx-auto max-w-lg px-4 py-10 font-[family-name:var(--font-geist-sans)]">
        <h1 className="text-xl font-semibold text-white">Account</h1>
        <p className="mt-2 text-sm text-zinc-400">
          Delete your account and all associated sessions, transcripts, and
          feedback (cascaded in the database).
        </p>
        {me && (
          <p className="mt-4 text-sm text-zinc-500">
            Signed in as{" "}
            <span className="font-mono text-zinc-300">{me.email}</span>
            {me.is_guest ? " (guest)" : ""}
          </p>
        )}
        {err && (
          <p className="mt-4 text-sm text-red-400" role="alert">
            {err}
          </p>
        )}
        {!me?.is_guest && (
          <label className="mt-6 block text-sm text-zinc-400">
            Confirm password
            <input
              type="password"
              autoComplete="current-password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="mt-1 w-full rounded-lg border border-zinc-700 bg-zinc-900 px-3 py-2 text-white"
            />
          </label>
        )}
        <div className="mt-8 flex flex-wrap gap-3">
          <button
            type="button"
            disabled={busy || (!me?.is_guest && !password.trim())}
            onClick={() => void deleteAccount()}
            className="rounded-lg bg-red-900/80 px-4 py-2 text-sm font-medium text-red-100 hover:bg-red-800 disabled:opacity-50"
          >
            {busy ? "Deleting…" : "Delete account"}
          </button>
          <Link
            href="/dashboard"
            className="rounded-lg border border-zinc-700 px-4 py-2 text-sm text-zinc-200 hover:bg-zinc-900"
          >
            Cancel
          </Link>
        </div>
      </main>
    </>
  );
}
