"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { api, ApiError } from "@/lib/api";
import { setToken } from "@/lib/auth";

type RegResp = { token: string };

export default function RegisterPage() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [err, setErr] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setErr(null);
    setLoading(true);
    try {
      const data = await api<RegResp>("/v1/auth/register", {
        method: "POST",
        body: JSON.stringify({ email, password }),
      });
      setToken(data.token);
      router.replace("/onboarding");
    } catch (ex) {
      setErr(ex instanceof ApiError ? ex.message : "Could not register");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="mx-auto flex min-h-screen max-w-md flex-col justify-center px-6 font-[family-name:var(--font-geist-sans)]">
      <h1 className="text-2xl font-semibold tracking-tight text-white">
        Create account
      </h1>
      <p className="mt-2 text-sm text-zinc-400">
        Password must be at least 8 characters.
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
            minLength={8}
            autoComplete="new-password"
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
          {loading ? "Creating…" : "Register"}
        </button>
      </form>
      <p className="mt-6 text-center text-sm text-zinc-500">
        Already have an account?{" "}
        <Link href="/login" className="text-emerald-400 hover:underline">
          Sign in
        </Link>
      </p>
    </div>
  );
}
