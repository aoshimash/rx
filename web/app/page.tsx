'use client';

import { ActivePlans } from '@/components/dashboard/ActivePlans';
import { E1rmChart } from '@/components/dashboard/E1rmChart';
import { RecentLogs } from '@/components/dashboard/RecentLogs';

export default function Home() {
  return (
    <main className="container mx-auto p-6 space-y-6">
      <h1 className="text-3xl font-bold">Dashboard</h1>

      <div className="grid gap-6 md:grid-cols-2">
        <E1rmChart />
        <div className="space-y-6">
          <RecentLogs />
          <ActivePlans />
        </div>
      </div>
    </main>
  );
}
