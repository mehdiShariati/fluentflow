"use client";

import { useParams } from "next/navigation";
import Link from "next/link";
import { AppNav } from "@/components/AppNav";
import SessionClient from "./SessionClient";

export default function SessionPage() {
  const params = useParams();
  const id = typeof params?.id === "string" ? params.id : "";

  if (!id) {
    return (
      <>
        <AppNav />
        <p className="p-8 text-center text-zinc-500">Invalid session.</p>
      </>
    );
  }

  return (
    <>
      <AppNav />
      <main className="mx-auto max-w-lg px-4 py-8">
        <SessionClient sessionId={id} />
        <p className="mt-8 text-center text-sm text-zinc-500">
          <Link href="/scenarios" className="text-emerald-400 hover:underline">
            All scenarios
          </Link>
        </p>
      </main>
    </>
  );
}
