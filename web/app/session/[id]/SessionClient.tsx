"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { ConnectionState, Room, RoomEvent, Track } from "livekit-client";
import { api, ApiError } from "@/lib/api";
import { getToken, liveKitStorageKey } from "@/lib/auth";

type SessionInfo = {
  session_id: string;
  scenario_id: string;
  scenario_title?: string;
  room_name: string;
  livekit_url: string;
  status: string;
};

type TokenResp = {
  livekit_token: string;
  livekit_url: string;
};

type TranscriptMessage = {
  id: string;
  segmentKey: string;
  speaker: "you" | "agent" | "system";
  text: string;
  isFinal?: boolean;
};

type LiveKitTranscriptSegment = {
  id?: string;
  text?: string;
  participantIdentity?: string;
  final?: boolean;
};

type ProfileResp = { target_language: string };
type TranslateResp = { translated_text: string; target_language: string; source: string };
type AnalyzeResp = { analysis: string; source: string };

export default function SessionClient({ sessionId }: { sessionId: string }) {
  const router = useRouter();
  const token = getToken();
  const roomRef = useRef<Room | null>(null);
  const audioElsRef = useRef<Map<string, HTMLAudioElement>>(new Map());
  /** Cumulative time spent in ConnectionState.Connected (ms), for dashboard metrics. */
  const speakingMsRef = useRef(0);
  const speakingSegmentStartRef = useRef<number | null>(null);
  const turnCountRef = useRef(0);
  const [session, setSession] = useState<SessionInfo | null>(null);
  const [conn, setConn] = useState<ConnectionState>(
    ConnectionState.Disconnected
  );
  const [remoteCount, setRemoteCount] = useState(0);
  const [msg, setMsg] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);
  const [messages, setMessages] = useState<TranscriptMessage[]>([]);
  const [agentPresent, setAgentPresent] = useState(false);
  const [targetLanguage, setTargetLanguage] = useState("English");
  const [translations, setTranslations] = useState<Record<string, string>>({});
  const [analyses, setAnalyses] = useState<Record<string, string>>({});
  const [translateBusy, setTranslateBusy] = useState<Record<string, boolean>>({});
  const [analyzeBusy, setAnalyzeBusy] = useState<Record<string, boolean>>({});
  const persistedTranscriptKeysRef = useRef<Set<string>>(new Set());
  const textHash = (s: string) =>
    String(
      Array.from(s).reduce((acc, ch) => ((acc << 5) - acc + ch.charCodeAt(0)) | 0, 0)
    );

  const upsertMessage = useCallback(
    (segmentKey: string, speaker: TranscriptMessage["speaker"], text: string, isFinal?: boolean) => {
      const trimmed = text.trim();
      if (!trimmed) return;
      setMessages((prev) => {
        const idx = prev.findIndex((m) => m.segmentKey === segmentKey);
        if (idx >= 0) {
          const next = [...prev];
          next[idx] = {
            ...next[idx],
            text: trimmed,
            isFinal: Boolean(isFinal) || next[idx].isFinal,
          };
          return next;
        }
        return [
          ...prev.slice(-119),
          {
            id: `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
            segmentKey,
            speaker,
            text: trimmed,
            isFinal,
          },
        ];
      });
    },
    []
  );

  const addSystem = useCallback(
    (text: string) => {
      upsertMessage(`system-${Date.now()}-${Math.random().toString(36).slice(2, 5)}`, "system", text, true);
    },
    [upsertMessage]
  );

  const persistTranscriptLines = useCallback(
    (entries: TranscriptMessage[]) => {
      if (!token || entries.length === 0) return;
      const fresh = entries.filter((e) => !persistedTranscriptKeysRef.current.has(e.segmentKey));
      if (fresh.length === 0) return;
      fresh.forEach((e) => persistedTranscriptKeysRef.current.add(e.segmentKey));
      void api(`/v1/sessions/${sessionId}/transcript`, {
        method: "POST",
        token,
        body: JSON.stringify({
          segments: fresh.map((e) => ({
            speaker: e.speaker === "agent" ? "assistant" : e.speaker,
            text: e.text,
          })),
        }),
      }).catch(() => {});
    },
    [sessionId, token]
  );

  const flushSpeakingSegment = useCallback(() => {
    if (speakingSegmentStartRef.current != null) {
      speakingMsRef.current += Date.now() - speakingSegmentStartRef.current;
      speakingSegmentStartRef.current = null;
    }
  }, []);

  const disconnect = useCallback(async () => {
    const r = roomRef.current;
    roomRef.current = null;
    if (r) {
      await r.disconnect();
    }
    audioElsRef.current.forEach((el) => {
      try {
        el.pause();
        el.srcObject = null;
        el.remove();
      } catch {
        // no-op cleanup
      }
    });
    audioElsRef.current.clear();
    flushSpeakingSegment();
    setConn(ConnectionState.Disconnected);
  }, [flushSpeakingSegment]);

  useEffect(() => {
    if (!token) return;
    void api<ProfileResp>("/v1/me/profile", { token })
      .then((p) => {
        if (p?.target_language?.trim()) setTargetLanguage(p.target_language);
      })
      .catch(() => {});
  }, [token]);

  useEffect(() => {
    if (!token) return;
    let cancelled = false;
    (async () => {
      try {
        const s = await api<SessionInfo>(`/v1/sessions/${sessionId}`, {
          token,
        });
        if (!cancelled) setSession(s);
      } catch {
        if (!cancelled) setMsg("Session not found");
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [token, sessionId]);

  useEffect(() => {
    return () => {
      void disconnect();
    };
  }, [disconnect]);

  async function ensureLiveKitToken(): Promise<TokenResp | null> {
    if (!token) return null;
    const key = liveKitStorageKey(sessionId);
    const cached = sessionStorage.getItem(key);
    if (cached && session?.livekit_url) {
      return { livekit_token: cached, livekit_url: session.livekit_url };
    }
    try {
      const t = await api<TokenResp>(
        `/v1/sessions/${sessionId}/livekit-token`,
        { method: "POST", token }
      );
      sessionStorage.setItem(key, t.livekit_token);
      return t;
    } catch (ex) {
      if (ex instanceof ApiError && ex.status === 503) {
        setMsg(
          "LiveKit is not configured on the API. Set LIVEKIT_* env vars to join a real room."
        );
        return null;
      }
      throw ex;
    }
  }

  async function connect() {
    if (!token || !session) return;
    setMsg(null);
    setBusy(true);
    try {
      await disconnect();
      const t = await ensureLiveKitToken();
      if (!t?.livekit_url || !t.livekit_token) {
        setBusy(false);
        return;
      }
      const room = new Room({
        adaptiveStream: true,
        dynacast: true,
      });
      roomRef.current = room;
      const refreshRemotes = () => {
        setRemoteCount(room.remoteParticipants.size);
        let hasAgent = false;
        room.remoteParticipants.forEach((rp) => {
          if (rp.identity.startsWith("agent-") || rp.identity.includes("tutor")) {
            hasAgent = true;
          }
        });
        setAgentPresent(hasAgent);
      };
      const attachAudioTrack = (participantId: string, track: Track) => {
        if (track.kind !== Track.Kind.Audio) return;
        const existing = audioElsRef.current.get(participantId);
        if (existing) {
          track.attach(existing);
          return;
        }
        const audioEl = document.createElement("audio");
        audioEl.autoplay = true;
        audioEl.dataset.participantId = participantId;
        document.body.appendChild(audioEl);
        track.attach(audioEl);
        audioElsRef.current.set(participantId, audioEl);
      };
      room.on(RoomEvent.ConnectionStateChanged, (s) => {
        setConn(s);
        if (s === ConnectionState.Connected) {
          speakingSegmentStartRef.current = Date.now();
        } else if (s === ConnectionState.Disconnected) {
          flushSpeakingSegment();
        }
      });
      room.on(RoomEvent.ParticipantConnected, refreshRemotes);
      room.on(RoomEvent.ParticipantDisconnected, refreshRemotes);
      room.on(RoomEvent.TrackSubscribed, (track, _pub, participant) => {
        attachAudioTrack(participant.sid, track);
      });
      room.on(RoomEvent.TrackUnsubscribed, (_track, _pub, participant) => {
        const el = audioElsRef.current.get(participant.sid);
        if (!el) return;
        try {
          el.pause();
          el.srcObject = null;
          el.remove();
        } catch {
          // no-op cleanup
        }
        audioElsRef.current.delete(participant.sid);
      });
      room.on(RoomEvent.TranscriptionReceived, (segments: LiveKitTranscriptSegment[]) => {
        const finalized: TranscriptMessage[] = [];
        for (const seg of segments ?? []) {
          const text = String(seg?.text ?? "");
          const pid = String(seg?.participantIdentity ?? "");
          const speaker: TranscriptMessage["speaker"] = pid.startsWith("agent-") ? "agent" : "you";
          const isFinal = Boolean(seg?.final);
          const key =
            seg.id && seg.id.trim()
              ? `${speaker}:${seg.id}`
              : `${speaker}:${pid || "unknown"}:${textHash(text.trim().toLowerCase())}`;
          upsertMessage(key, speaker, text, isFinal);
          if (isFinal) {
            finalized.push({
              id: key,
              segmentKey: key,
              speaker,
              text,
              isFinal: true,
            });
          }
        }
        persistTranscriptLines(finalized);
      });
      await room.connect(t.livekit_url, t.livekit_token);
      refreshRemotes();
      room.remoteParticipants.forEach((rp) => {
        rp.audioTrackPublications.forEach((pub) => {
          if (pub.track) {
            attachAudioTrack(rp.sid, pub.track);
          }
        });
      });
      await room.localParticipant.setMicrophoneEnabled(true);
      addSystem("Connected. Speak naturally; transcript appears below.");
      if (!agentPresent) {
        addSystem(
          "Connected, waiting for AI tutor to join. If this persists, ensure the agent container is running."
        );
      }
      await api(`/v1/sessions/${sessionId}/events`, {
        method: "POST",
        token,
        body: JSON.stringify({
          events: [
            { type: "session_joined", payload: { livekit: true } },
            { type: "turn_started", payload: { role: "user" } },
          ],
        }),
      });
      turnCountRef.current += 1;
    } catch (e) {
      setMsg(e instanceof Error ? e.message : "Connection failed");
      addSystem(
        e instanceof Error ? `Connection error: ${e.message}` : "Connection failed"
      );
    } finally {
      setBusy(false);
    }
  }

  async function translateMessage(m: TranscriptMessage) {
    if (!token || m.speaker === "system") return;
    setTranslateBusy((s) => ({ ...s, [m.id]: true }));
    try {
      const out = await api<TranslateResp>("/v1/ai/translate", {
        method: "POST",
        token,
        body: JSON.stringify({ text: m.text, target_language: targetLanguage }),
      });
      setTranslations((s) => ({ ...s, [m.id]: out.translated_text }));
    } catch {
      setTranslations((s) => ({ ...s, [m.id]: "Translation failed." }));
    } finally {
      setTranslateBusy((s) => ({ ...s, [m.id]: false }));
    }
  }

  async function analyzeMessage(m: TranscriptMessage) {
    if (!token || m.speaker === "system") return;
    setAnalyzeBusy((s) => ({ ...s, [m.id]: true }));
    try {
      const out = await api<AnalyzeResp>("/v1/ai/analyze", {
        method: "POST",
        token,
        body: JSON.stringify({ text: m.text, target_language: targetLanguage }),
      });
      setAnalyses((s) => ({ ...s, [m.id]: out.analysis }));
    } catch {
      setAnalyses((s) => ({ ...s, [m.id]: "Analysis failed." }));
    } finally {
      setAnalyzeBusy((s) => ({ ...s, [m.id]: false }));
    }
  }

  async function endSession() {
    if (!token) return;
    setBusy(true);
    try {
      void api(`/v1/sessions/${sessionId}/events`, {
        method: "POST",
        token,
        body: JSON.stringify({
          events: [{ type: "turn_completed", payload: { role: "user" } }],
        }),
      }).catch(() => {});
      await disconnect();
      const speakingSeconds = Math.max(
        0,
        Math.round(speakingMsRef.current / 1000)
      );
      await api(`/v1/sessions/${sessionId}/complete`, {
        method: "POST",
        token,
        body: JSON.stringify({
          speaking_seconds: speakingSeconds,
          turn_count: turnCountRef.current,
        }),
      });
      const completeLines = messages.filter((l) => l.speaker !== "system" && l.isFinal);
      if (completeLines.length > 0) {
        persistTranscriptLines(completeLines);
      } else {
        // Keep a fallback line for post-session feedback if no transcript events were emitted.
        void api(`/v1/sessions/${sessionId}/transcript`, {
          method: "POST",
          token,
          body: JSON.stringify({
            segments: [
              {
                speaker: "user",
                text: "[Session note] Learner completed the scenario.",
                offset_ms: 0,
              },
            ],
          }),
        }).catch(() => {});
      }
      // session_completed is recorded server-side in POST …/complete (Prometheus + DB).
      sessionStorage.removeItem(liveKitStorageKey(sessionId));
      router.push(`/feedback/${sessionId}`);
    } catch (e) {
      setMsg(e instanceof ApiError ? e.message : "Could not end session");
    } finally {
      setBusy(false);
    }
  }

  if (!token) {
    return (
      <p className="text-center text-zinc-500">
        <Link href="/login" className="text-emerald-400">
          Sign in
        </Link>
      </p>
    );
  }

  if (!session) {
    return (
      <p className="text-center text-zinc-500">
        {msg ?? "Loading session…"}
      </p>
    );
  }

  return (
    <div className="space-y-6 font-[family-name:var(--font-geist-sans)]">
      <div>
        <p className="text-xs uppercase tracking-wide text-zinc-500">
          Scenario
        </p>
        <h1 className="text-lg font-semibold text-white">
          {session.scenario_title || session.scenario_id}
        </h1>
        {session.scenario_title && (
          <p className="mt-0.5 font-mono text-xs text-zinc-500">
            {session.scenario_id}
          </p>
        )}
        <p className="mt-1 font-mono text-xs text-zinc-500">
          room: {session.room_name}
        </p>
      </div>
      <div className="rounded-xl border border-zinc-800 bg-zinc-900/50 p-4">
        <p className="text-sm text-zinc-300">
          Connection:{" "}
          <span className="font-medium text-emerald-400">{conn}</span>
        </p>
        <p className="mt-2 text-sm text-zinc-400">
          Other participants in room (e.g. AI agent):{" "}
          <span className="text-white">{remoteCount}</span>
        </p>
        <p className="mt-1 text-xs text-zinc-500">
          Agent status:{" "}
          <span className={agentPresent ? "text-emerald-400" : "text-amber-400"}>
            {agentPresent ? "joined" : "waiting"}
          </span>
        </p>
        <p className="mt-3 text-xs text-zinc-500">
          Microphone is enabled after connect so you can speak when an agent is
          in the room. Deploy a LiveKit agent (PRD §12.3) to hear the tutor.
        </p>
      </div>
      <section className="rounded-xl border border-zinc-800 bg-zinc-900/50 p-4">
        <h2 className="text-sm font-medium uppercase tracking-wide text-zinc-500">
          Live chat transcript
        </h2>
        {messages.length === 0 ? (
          <p className="mt-3 text-sm text-zinc-500">
            Connect to room to see what you and the agent say.
          </p>
        ) : (
          <ul className="mt-3 max-h-72 space-y-3 overflow-y-auto text-sm">
            {messages.map((l) => (
              <li
                key={l.id}
                className={`rounded-xl border px-3 py-2 ${
                  l.speaker === "you"
                    ? "ml-8 border-sky-800/70 bg-sky-950/30"
                    : l.speaker === "agent"
                    ? "mr-8 border-emerald-800/70 bg-emerald-950/20"
                    : "border-zinc-800 bg-zinc-950/60"
                }`}
              >
                <span
                  className={
                    l.speaker === "agent"
                      ? "font-medium text-emerald-400"
                      : l.speaker === "you"
                      ? "font-medium text-sky-400"
                      : "font-medium text-zinc-400"
                  }
                >
                  {l.speaker === "agent"
                    ? "Agent"
                    : l.speaker === "you"
                    ? "You"
                    : "System"}
                  {l.isFinal === false ? " (live)" : ""}
                </span>
                <p className="mt-1 text-zinc-200">{l.text}</p>
                {l.speaker !== "system" && (
                  <div className="mt-2 flex flex-wrap gap-2">
                    <button
                      type="button"
                      disabled={Boolean(translateBusy[l.id])}
                      onClick={() => void translateMessage(l)}
                      className="rounded-md border border-zinc-700 px-2 py-1 text-xs text-zinc-300 hover:bg-zinc-900 disabled:opacity-50"
                    >
                      {translateBusy[l.id] ? "Translating..." : `Translate (${targetLanguage})`}
                    </button>
                    <button
                      type="button"
                      disabled={Boolean(analyzeBusy[l.id])}
                      onClick={() => void analyzeMessage(l)}
                      className="rounded-md border border-zinc-700 px-2 py-1 text-xs text-zinc-300 hover:bg-zinc-900 disabled:opacity-50"
                    >
                      {analyzeBusy[l.id] ? "Analyzing..." : "Analyze"}
                    </button>
                  </div>
                )}
                {translations[l.id] && (
                  <p className="mt-2 rounded-md border border-indigo-900/60 bg-indigo-950/30 px-2 py-1 text-xs text-indigo-200">
                    {translations[l.id]}
                  </p>
                )}
                {analyses[l.id] && (
                  <pre className="mt-2 whitespace-pre-wrap rounded-md border border-amber-900/60 bg-amber-950/30 px-2 py-1 text-xs text-amber-100">
                    {analyses[l.id]}
                  </pre>
                )}
              </li>
            ))}
          </ul>
        )}
      </section>
      {msg && (
        <p className="rounded-lg border border-amber-900/50 bg-amber-950/30 px-3 py-2 text-sm text-amber-200">
          {msg}
        </p>
      )}
      <div className="flex flex-wrap gap-3">
        <button
          type="button"
          disabled={busy || session.status !== "active"}
          onClick={() => void connect()}
          className="rounded-lg bg-emerald-600 px-4 py-2 text-sm font-medium text-white hover:bg-emerald-500 disabled:opacity-50"
        >
          {busy ? "Working…" : "Connect to room"}
        </button>
        <button
          type="button"
          disabled={busy}
          onClick={() => void disconnect()}
          className="rounded-lg border border-zinc-700 px-4 py-2 text-sm text-zinc-200 hover:bg-zinc-900"
        >
          Disconnect
        </button>
        <button
          type="button"
          disabled={busy || session.status !== "active"}
          onClick={() => void endSession()}
          className="rounded-lg bg-red-900/80 px-4 py-2 text-sm font-medium text-red-100 hover:bg-red-800 disabled:opacity-50"
        >
          End &amp; feedback
        </button>
      </div>
    </div>
  );
}
