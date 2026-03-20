import type { PlanStatus } from '@/lib/utils/next-session';
import { PlanCard } from './PlanCard';

interface PlanListProps {
  statuses: PlanStatus[];
  onRecordLog: (planId: string) => void;
}

export function PlanList({ statuses, onRecordLog }: PlanListProps) {
  return (
    <div className="space-y-3">
      {statuses.map((status) => (
        <PlanCard
          key={status.plan.id}
          status={status}
          onRecordLog={() => onRecordLog(status.plan.id)}
        />
      ))}
    </div>
  );
}
