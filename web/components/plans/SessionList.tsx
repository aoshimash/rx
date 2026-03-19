import type { SessionStatus } from '@/lib/utils/next-session';
import { SessionCard } from './SessionCard';

interface SessionListProps {
  statuses: SessionStatus[];
  onRecordLog: (sessionName: string) => void;
}

export function SessionList({ statuses, onRecordLog }: SessionListProps) {
  return (
    <div className="space-y-3">
      {statuses.map((status) => (
        <SessionCard
          key={status.session.name}
          status={status}
          onRecordLog={() => onRecordLog(status.session.name)}
        />
      ))}
    </div>
  );
}
