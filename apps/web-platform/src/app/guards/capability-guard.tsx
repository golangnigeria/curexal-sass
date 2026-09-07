import React from "react";
import { CapabilityGate } from "@/components/design-system/capability-gate";

export interface CapabilityGuardProps {
  capability: string;
  moduleCode?: string;
  title?: string;
  description?: string;
  requiredPlan?: string;
  children: React.ReactNode;
}

export const CapabilityGuard: React.FC<CapabilityGuardProps> = ({
  capability,
  moduleCode,
  title,
  description,
  requiredPlan,
  children,
}) => {
  return (
    <CapabilityGate
      capability={capability}
      moduleCode={moduleCode}
      title={title}
      description={description}
      requiredPlan={requiredPlan}
    >
      {children}
    </CapabilityGate>
  );
};
