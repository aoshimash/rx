import type { SessionGroup } from '@/components/plan-editor/types';
import type { Log } from '@/types/api';

export interface SessionStatus {
  session: SessionGroup;
  isNext: boolean;
  completedCount: number;
  lastPerformedAt?: string;
}

/**
 * Detect the next session based on logs for a given plan.
 * Looks at the most recent log's metadata.session_name to determine
 * which session was done last, then advances to the next one cyclically.
 */
export function detectNextSession(
  sessions: SessionGroup[],
  logs: Log[],
  planId: string
): SessionStatus[] {
  if (sessions.length === 0) return [];

  const planLogs = logs
    .filter((log) => log.plan_id === planId)
    .sort((a, b) => new Date(b.performed_at).getTime() - new Date(a.performed_at).getTime());

  // Count completions and find last performed date per session
  const sessionStats = new Map<string, { count: number; lastAt?: string }>();
  for (const log of planLogs) {
    const sessionName = log.metadata?.session_name as string | undefined;
    if (!sessionName) continue;
    const existing = sessionStats.get(sessionName) || { count: 0 };
    existing.count++;
    if (!existing.lastAt) {
      existing.lastAt = log.performed_at;
    }
    sessionStats.set(sessionName, existing);
  }

  // Determine next session index
  let nextIndex = 0;
  const lastLog = planLogs[0];
  if (lastLog) {
    const lastSessionName = lastLog.metadata?.session_name as string | undefined;
    if (lastSessionName) {
      const lastIndex = sessions.findIndex((s) => s.name === lastSessionName);
      if (lastIndex >= 0) {
        nextIndex = (lastIndex + 1) % sessions.length;
      }
    }
  }

  return sessions.map((session, index) => {
    const stats = sessionStats.get(session.name);
    return {
      session,
      isNext: index === nextIndex,
      completedCount: stats?.count ?? 0,
      lastPerformedAt: stats?.lastAt,
    };
  });
}

/**
 * Rotate session statuses so that the "next" session appears first.
 */
export function sortSessionsByNext(statuses: SessionStatus[]): SessionStatus[] {
  if (statuses.length === 0) return [];
  const nextIndex = statuses.findIndex((s) => s.isNext);
  if (nextIndex <= 0) return statuses;
  return [...statuses.slice(nextIndex), ...statuses.slice(0, nextIndex)];
}
