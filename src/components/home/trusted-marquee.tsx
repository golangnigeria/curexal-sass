import { Activity } from "lucide-react";

const trustedNames = [
  "LifeCare Diagnostics",
  "Everight Pathology",
  "Apex Laboratory Network",
  "Meridian Health Group",
  "Lumina Pathology Group",
  "MedTest Clinical Labs",
  "Synlab West Africa",
  "Diagnostics Global Inc",
];

export function TrustedMarquee() {
  return (
    <section className="py-12 border-y border-border bg-muted/20 overflow-hidden">
      <p className="text-center text-xs text-muted-foreground/80 font-semibold uppercase tracking-widest mb-8">
        Healthcare providers connected on the Curexal network
      </p>
      <div className="flex gap-12 animate-marquee whitespace-nowrap">
        {[...trustedNames, ...trustedNames].map((name, i) => (
          <span
            key={i}
            className="text-muted-foreground/75 font-semibold text-sm flex items-center gap-2"
          >
            <Activity className="h-3.5 w-3.5 text-primary/40" />
            {name}
          </span>
        ))}
      </div>
    </section>
  );
}
