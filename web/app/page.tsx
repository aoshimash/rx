'use client';

import { ActivePlans } from '@/components/dashboard/ActivePlans';
import { RecentLogs } from '@/components/dashboard/RecentLogs';

export default function Home() {
  return (
    <main className="container mx-auto p-6 space-y-6">
      <h1 className="text-3xl font-bold">Dashboard</h1>

      <div className="grid gap-6 md:grid-cols-2">
        <RecentLogs />
        <ActivePlans />
      </div>
    </main>
  );
}
