"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { api, ApiError } from "@/lib/api";
import { getToken, setToken } from "@/lib/auth";

type LoginResp = { token: string };

export default function LoginPage() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [err, setErr] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function continueAsGuest() {
    setErr(null);
    setLoading(true);
    try {
      const data = await api<LoginResp>("/v1/auth/guest", { method: "POST" });
      setToken(data.token);
      router.replace("/onboarding");
    } catch (ex) {
      setErr(ex instanceof ApiError ? ex.message : "Guest sign-in failed");
    } finally {
      setLoading(false);
    }
  }

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setErr(null);
    setLoading(true);
    try {
      const data = await api<LoginResp>("/v1/auth/login", {
        method: "POST",
        body: JSON.stringify({ email, password }),
      });
      setToken(data.token);
      const tok = getToken();
      const prof = await api<{
        target_language: string;
      }>("/v1/me/profile", { token: tok });
      if (!prof.target_language) router.replace("/onboarding");
      else router.replace("/dashboard");
    } catch (ex) {
      setErr(ex instanceof ApiError ? ex.message : "Login failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="mx-auto flex min-h-screen max-w-md flex-col justify-center px-6 font-[family-name:var(--font-geist-sans)]">
      <h1 className="text-2xl font-semibold tracking-tight text-white">
        Sign in
      </h1>
      <p className="mt-2 text-sm text-zinc-400">
        FluentFlow — practice speaking with clear goals and feedback.
      </p>
      <form onSubmit={onSubmit} className="mt-8 flex flex-col gap-4">
        <label className="block text-sm text-zinc-300">
          Email
          <input
            type="email"
            required
            autoComplete="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="mt-1 w-full rounded-lg border border-zinc-800 bg-zinc-900 px-3 py-2 text-white outline-none ring-emerald-500/40 focus:ring-2"
          />
        </label>
        <label className="block text-sm text-zinc-300">
          Password
          <input
            type="password"
            required
            autoComplete="current-password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="mt-1 w-full rounded-lg border border-zinc-800 bg-zinc-900 px-3 py-2 text-white outline-none ring-emerald-500/40 focus:ring-2"
          />
        </label>
        {err && (
          <p className="text-sm text-red-400" role="alert">
            {err}
          </p>
        )}
        <button
          type="submit"
          disabled={loading}
          className="rounded-lg bg-emerald-600 px-4 py-2.5 text-sm font-medium text-white transition hover:bg-emerald-500 disabled:opacity-50"
        >
          {loading ? "Signing in…" : "Continue"}
        </button>
      </form>
      <p className="mt-6 text-center text-sm text-zinc-500">
        No account?{" "}
        <Link href="/register" className="text-emerald-400 hover:underline">
          Register
        </Link>
      </p>
      <p className="mt-4 text-center">
        <button
          type="button"
          disabled={loading}
          onClick={() => void continueAsGuest()}
          className="text-sm text-zinc-400 underline decoration-zinc-600 hover:text-zinc-200 disabled:opacity-50"
        >
          Continue as guest
        </button>
      </p>
    </div>
  );
}
