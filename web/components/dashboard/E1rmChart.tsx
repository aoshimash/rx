'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { useLogs } from '@/lib/hooks/useLogs';
import { calculateE1rm } from '@/lib/utils/e1rm';
import type { Log } from '@/types/api';
import { format, parseISO, subWeeks } from 'date-fns';
import { useMemo, useState } from 'react';
import {
  CartesianGrid,
  Line,
  LineChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts';

type ChartPoint = { date: string; e1rm: number };

function getLogPoints(log: Log, exercise: string, cutoff: Date): ChartPoint[] {
  const logDate = parseISO(log.performed_at);
  if (logDate < cutoff) return [];

  return log.entries.flatMap((entry) => {
    if (entry.exercise_name !== exercise || !entry.load_kg || !entry.reps) return [];
    const e1rm = calculateE1rm(entry.load_kg, entry.reps, entry.rpe ?? undefined);
    if (e1rm === null) return [];
    return [{ date: format(logDate, 'MM/dd'), e1rm: Math.round(e1rm * 10) / 10 }];
  });
}

function buildChartData(logs: Log[], exercise: string): ChartPoint[] {
  const cutoff = subWeeks(new Date(), 12);
  const points = logs.flatMap((log) => getLogPoints(log, exercise, cutoff));
  return points.sort((a, b) => a.date.localeCompare(b.date));
}

export function E1rmChart() {
  const { data } = useLogs();
  const [selectedExercise, setSelectedExercise] = useState<string>('');

  const exerciseNames = useMemo(() => {
    if (!data?.data) return [];
    const names = new Set<string>();
    for (const log of data.data) {
      for (const entry of log.entries) {
        names.add(entry.exercise_name);
      }
    }
    return Array.from(names).sort();
  }, [data?.data]);

  const chartData = useMemo(
    () => (data?.data && selectedExercise ? buildChartData(data.data, selectedExercise) : []),
    [data?.data, selectedExercise]
  );

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-base">Estimated 1RM</CardTitle>
        <Select value={selectedExercise} onValueChange={setSelectedExercise}>
          <SelectTrigger className="w-48 h-8 text-sm">
            <SelectValue placeholder="Select exercise" />
          </SelectTrigger>
          <SelectContent>
            {exerciseNames.map((name) => (
              <SelectItem key={name} value={name}>
                {name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </CardHeader>
      <CardContent>
        {!selectedExercise && (
          <p className="text-sm text-muted-foreground py-8 text-center">
            Select an exercise to view e1RM trend.
          </p>
        )}
        {selectedExercise && chartData.length === 0 && (
          <p className="text-sm text-muted-foreground py-8 text-center">
            No data for {selectedExercise} in the last 12 weeks.
          </p>
        )}
        {selectedExercise && chartData.length > 0 && (
          <ResponsiveContainer width="100%" height={240}>
            <LineChart data={chartData}>
              <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
              <XAxis dataKey="date" tick={{ fontSize: 12 }} />
              <YAxis tick={{ fontSize: 12 }} unit="kg" />
              <Tooltip formatter={(value) => [`${value} kg`, 'e1RM']} />
              <Line
                type="monotone"
                dataKey="e1rm"
                stroke="hsl(var(--primary))"
                strokeWidth={2}
                dot={{ r: 3 }}
                activeDot={{ r: 5 }}
              />
            </LineChart>
          </ResponsiveContainer>
        )}
      </CardContent>
    </Card>
  );
}
