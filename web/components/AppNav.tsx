"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { clearToken, getToken } from "@/lib/auth";

const links = [
  { href: "/dashboard", label: "Dashboard" },
  { href: "/scenarios", label: "Scenarios" },
  { href: "/onboarding", label: "Profile" },
  { href: "/settings", label: "Account" },
];

export function AppNav() {
  const pathname = usePathname();
  const router = useRouter();

  if (!getToken()) return null;

  return (
    <header className="border-b border-zinc-800 bg-zinc-950/80 backdrop-blur">
      <div className="mx-auto flex max-w-3xl items-center justify-between gap-4 px-4 py-3 font-[family-name:var(--font-geist-sans)]">
        <Link
          href="/dashboard"
          className="text-sm font-semibold tracking-tight text-emerald-400"
        >
          FluentFlow
        </Link>
        <nav className="flex flex-wrap items-center gap-3 text-sm text-zinc-400">
          {links.map((l) => (
            <Link
              key={l.href}
              href={l.href}
              className={
                pathname === l.href ? "text-white" : "hover:text-zinc-200"
              }
            >
              {l.label}
            </Link>
          ))}
          <button
            type="button"
            onClick={() => {
              clearToken();
              router.replace("/login");
            }}
            className="rounded-md border border-zinc-700 px-2 py-1 text-zinc-300 hover:bg-zinc-900"
          >
            Sign out
          </button>
        </nav>
      </div>
    </header>
  );
}
