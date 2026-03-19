import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import type { Plan } from '@/types/api';

interface PlanSelectorProps {
  plans: Plan[];
  selectedPlanId: string | null;
  onSelect: (planId: string) => void;
}

export function PlanSelector({ plans, selectedPlanId, onSelect }: PlanSelectorProps) {
  if (plans.length <= 1) return null;

  return (
    <Select value={selectedPlanId ?? undefined} onValueChange={onSelect}>
      <SelectTrigger className="w-[240px]">
        <SelectValue placeholder="Select a plan" />
      </SelectTrigger>
      <SelectContent>
        {plans.map((plan) => (
          <SelectItem key={plan.id} value={plan.id}>
            {plan.name}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
