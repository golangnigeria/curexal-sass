import React from "react";
import { Link } from "react-router-dom";
import { useCapabilities } from "@/api/hooks/use-capabilities";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Sparkles, Lock, ArrowRight, ShieldCheck } from "lucide-react";

interface CapabilityGateProps {
  capability: string;
  moduleCode?: string;
  fallback?: React.ReactNode;
  title?: string;
  description?: string;
  requiredPlan?: string;
  children: React.ReactNode;
}

export function CapabilityGate({
  capability,
  moduleCode,
  fallback,
  title = "Commercial Feature Locked",
  description = "This advanced clinical or diagnostic capability requires an active add-on subscription or plan upgrade.",
  requiredPlan = "Optimize or Pro",
  children,
}: CapabilityGateProps) {
  const { hasCapability, isModuleLicensed, isLoading } = useCapabilities();

  if (isLoading) {
    return (
      <div className="p-12 text-center text-xs text-muted-foreground animate-pulse">
        Verifying capability licenses...
      </div>
    );
  }

  const isEntitled = hasCapability(capability) || (moduleCode ? isModuleLicensed(moduleCode) : false);

  if (isEntitled) {
    return <>{children}</>;
  }

  if (fallback) {
    return <>{fallback}</>;
  }

  return (
    <Card className="border-border shadow-md bg-gradient-to-b from-card to-secondary/30 p-8 max-w-xl mx-auto my-8 text-center space-y-4">
      <div className="mx-auto w-12 h-12 rounded-2xl bg-primary/10 text-primary flex items-center justify-center shadow-inner">
        <Lock className="w-6 h-6" />
      </div>

      <div className="space-y-1.5">
        <div className="flex items-center justify-center gap-2">
          <CardTitle className="text-base font-bold text-foreground">{title}</CardTitle>
          <Badge variant="outline" className="border-primary/40 text-primary bg-primary/10 text-[9px] font-mono uppercase">
            {requiredPlan}
          </Badge>
        </div>
        <CardDescription className="text-xs max-w-md mx-auto">
          {description}
        </CardDescription>
      </div>

      <div className="pt-2 flex flex-col sm:flex-row items-center justify-center gap-3">
        <Button asChild size="sm" className="text-xs h-9 gap-1.5 bg-primary text-primary-foreground shadow">
          <Link to="/organization/billing">
            <Sparkles className="w-3.5 h-3.5" />
            Upgrade Plan or Unlock Add-On
          </Link>
        </Button>
      </div>
    </Card>
  );
}
