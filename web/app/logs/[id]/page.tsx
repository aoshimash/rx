'use client';

import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { videosApi } from '@/lib/api/videos';
import { useLog } from '@/lib/hooks/useLogs';
import type { LogEntry, LogSet } from '@/types/api';
import { ArrowLeft, Pencil } from 'lucide-react';
import Link from 'next/link';
import { useParams } from 'next/navigation';
import { useEffect, useState } from 'react';

type LogExerciseGroup = { name: string; entries: LogEntry[] };

function groupByExercise(entries: LogEntry[]): LogExerciseGroup[] {
  const groups: LogExerciseGroup[] = [];
  const map = new Map<string, LogExerciseGroup>();
  for (const entry of [...entries].sort((a, b) => a.order - b.order)) {
    if (!map.has(entry.exercise_name)) {
      const g: LogExerciseGroup = { name: entry.exercise_name, entries: [] };
      groups.push(g);
      map.set(entry.exercise_name, g);
    }
    map.get(entry.exercise_name)?.entries.push(entry);
  }
  return groups;
}

// ============================================================================
// VideoPlayer — fetches a pre-signed download URL and renders a video element
// ============================================================================

function VideoPlayer({ objectKey }: { objectKey: string }) {
  const [url, setUrl] = useState<string | null>(null);
  const [error, setError] = useState(false);

  useEffect(() => {
    let cancelled = false;
    videosApi
      .getDownloadUrl({ object_key: objectKey })
      .then(({ download_url }) => {
        if (!cancelled) setUrl(download_url);
      })
      .catch(() => {
        if (!cancelled) setError(true);
      });
    return () => {
      cancelled = true;
    };
  }, [objectKey]);

  if (error) return <p className="text-xs text-destructive">Video unavailable</p>;
  if (!url) return <p className="text-xs text-muted-foreground">Loading video…</p>;

  return (
    // biome-ignore lint/a11y/useMediaCaption: captions not available for user-uploaded training videos
    <video src={url} controls className="max-w-full rounded mt-1" style={{ maxHeight: '200px' }} />
  );
}

// ============================================================================
// SetRow — displays a single LogSet with dynamic fields and optional video
// ============================================================================

function SetRow({ set }: { set: LogSet }) {
  const fieldEntries = Object.entries(set.fields ?? {}).filter(([, v]) => v != null);

  return (
    <div className="space-y-1">
      <div className="flex flex-wrap gap-3 text-sm text-muted-foreground">
        <span className="text-xs font-medium">Set {set.set_number}</span>
        {fieldEntries.map(([key, val]) => (
          <span key={key} className="tabular-nums">
            {key}: {String(val)}
          </span>
        ))}
        {set.notes && <span className="text-xs italic">{set.notes}</span>}
      </div>
      {set.video_object_key && <VideoPlayer objectKey={set.video_object_key} />}
    </div>
  );
}

// ============================================================================
// LegacyEntryRow — renders a single LogEntry in the legacy fields table
// ============================================================================

function LegacyEntryRow({
  entry,
  hasLabel,
  hasLoad,
}: {
  entry: LogEntry;
  hasLabel: boolean;
  hasLoad: boolean;
}) {
  const label = entry.fields?.label as string | undefined;
  const loadKg = entry.fields?.load_kg as number | null | undefined;
  const reps = entry.fields?.reps as number | null | undefined;
  const sets = entry.fields?.sets as number | null | undefined;

  const hasStandardFields = reps != null || sets != null || loadKg != null;
  if (!hasStandardFields) {
    const kvPairs = Object.entries(entry.fields ?? {}).filter(([, v]) => v != null);
    if (kvPairs.length === 0) return null;
    return (
      <tr key={entry.id}>
        <td colSpan={3} className="py-0.5">
          <div className="flex flex-wrap gap-3 text-sm text-muted-foreground">
            {kvPairs.map(([k, v]) => (
              <span key={k} className="tabular-nums">
                {k}: {String(v)}
              </span>
            ))}
          </div>
        </td>
      </tr>
    );
  }

  return (
    <tr key={entry.id} className="text-muted-foreground">
      {hasLabel && <td className="text-xs pr-3 py-0.5">{label ?? ''}</td>}
      {hasLoad && (
        <td className="text-right tabular-nums pr-4 py-0.5">
          {loadKg != null ? `${loadKg}kg` : '—'}
        </td>
      )}
      <td className="text-right tabular-nums pr-4 py-0.5">{reps ?? '—'}</td>
      <td className="text-right tabular-nums pr-4 py-0.5">{sets ?? '—'}</td>
    </tr>
  );
}

// ============================================================================
// EntryCard — displays one exercise group
// ============================================================================

function EntryCard({ group }: { group: LogExerciseGroup }) {
  const hasSets = group.entries.some((e) => e.sets && e.sets.length > 0);
  const hasLabel = group.entries.some((e) => e.fields?.label);
  const hasLoad = group.entries.some((e) => e.fields?.load_kg != null);

  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-base">{group.name}</CardTitle>
      </CardHeader>
      <CardContent>
        {hasSets ? (
          // New model: render sets with dynamic fields and video
          <div className="space-y-3">
            {group.entries.map((entry) =>
              (entry.sets ?? []).map((set) => <SetRow key={set.id} set={set} />)
            )}
          </div>
        ) : (
          // Legacy model: render entry.fields as key-value pairs
          <table className="w-full text-sm">
            <thead>
              <tr className="text-xs text-muted-foreground">
                {hasLabel && <th className="text-left font-normal pb-1 w-16" />}
                {hasLoad && <th className="text-right font-normal pb-1 pr-4">Load</th>}
                <th className="text-right font-normal pb-1 pr-4">Reps</th>
                <th className="text-right font-normal pb-1 pr-4">Sets</th>
              </tr>
            </thead>
            <tbody>
              {group.entries.map((entry) => (
                <LegacyEntryRow
                  key={entry.id}
                  entry={entry}
                  hasLabel={hasLabel}
                  hasLoad={hasLoad}
                />
              ))}
            </tbody>
          </table>
        )}
      </CardContent>
    </Card>
  );
}

// ============================================================================
// LogDetailPage
// ============================================================================

export default function LogDetailPage() {
  const params = useParams();
  const logId = params.id as string;

  const { data: log, isLoading: logLoading, error: logError } = useLog(logId);

  const performedDate = log
    ? new Date(log.performed_at).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
      })
    : '';

  const formatTime = (iso: string) =>
    new Date(iso).toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      hour12: false,
    });

  if (logLoading) {
    return (
      <main className="container mx-auto p-6 space-y-4">
        <Skeleton className="h-8 w-[200px]" />
        <Skeleton className="h-12 w-[400px]" />
        <Skeleton className="h-[400px] w-full" />
      </main>
    );
  }

  if (logError || !log) {
    return (
      <main className="container mx-auto p-6">
        <p className="text-destructive">Failed to load log. Please try again later.</p>
        <Link href="/logs" className="text-primary hover:underline mt-4 inline-block">
          Back to Logs
        </Link>
      </main>
    );
  }

  const groups = groupByExercise(log.entries);

  return (
    <main className="container mx-auto p-6">
      <Link
        href="/logs"
        className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground mb-4"
      >
        <ArrowLeft className="h-4 w-4" />
        Back to Logs
      </Link>

      <div className="mb-6">
        <h1 className="text-3xl font-bold">{performedDate}</h1>
        {log.session_name && (
          <p className="text-muted-foreground mt-1">Session: {log.session_name}</p>
        )}
        {(log.started_at || log.finished_at) && (
          <p className="text-muted-foreground mt-1 text-sm">
            {log.started_at && <>Start: {formatTime(log.started_at)}</>}
            {log.started_at && log.finished_at && <span className="mx-2">·</span>}
            {log.finished_at && <>End: {formatTime(log.finished_at)}</>}
            {log.started_at && log.finished_at && (
              <>
                <span className="mx-2">·</span>
                {Math.round(
                  (new Date(log.finished_at).getTime() - new Date(log.started_at).getTime()) / 60000
                )}{' '}
                min
              </>
            )}
          </p>
        )}
        {log.notes && <p className="text-muted-foreground mt-2">{log.notes}</p>}
        <div className="flex items-center gap-2 mt-2">
          <Button variant="outline" size="sm" asChild>
            <Link href={`/logs/${logId}/edit`}>
              <Pencil className="h-4 w-4 mr-1" />
              Edit
            </Link>
          </Button>
        </div>
      </div>

      <div className="space-y-4">
        {groups.map((group) => (
          <EntryCard key={group.name} group={group} />
        ))}
      </div>
    </main>
  );
}
