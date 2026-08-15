import {
  FlaskConical,
  ClipboardCheck,
  Users,
  FileText,
  Lock,
  History,
  Check,
} from "lucide-react";

const features = [
  {
    icon: FlaskConical,
    color: "text-sky-500 dark:text-sky-400",
    bg: "bg-sky-500/10",
    title: "Specimen Logging & Tracking",
    desc: "Register incoming specimens, print unique barcode identifiers, track physical locations, and maintain an unbroken chain of custody from draw to analysis.",
    items: [
      "Barcode label generation",
      "Chain of custody logging",
      "Multi-specimen batch intake",
      "Specimen source & stability tracking"
    ],
  },
  {
    icon: ClipboardCheck,
    color: "text-violet-500 dark:text-violet-400",
    bg: "bg-violet-500/10",
    title: "Analyte & Reference Ranges",
    desc: "Define clinical test profiles with custom analyte parameters. Set custom biological reference intervals with automatic flagging of abnormal results.",
    items: [
      "Biological reference intervals",
      "Critical value flags (H/L)",
      "Hematology, Biochemistry & Immunology",
      "Custom testing profile designer"
    ],
  },
  {
    icon: Users,
    color: "text-emerald-500 dark:text-emerald-400",
    bg: "bg-emerald-500/10",
    title: "Pathology Validation Queues",
    desc: "Enforce a secure diagnostic workflow. Technicians run panels and record results, while pathologists verify, validate, and sign off reports.",
    items: [
      "Technician worklists",
      "Pathologist validation queue",
      "Dual-verification workflow",
      "Secure digital sign-off signatures"
    ],
  },
  {
    icon: FileText,
    color: "text-amber-500 dark:text-amber-400",
    bg: "bg-amber-500/10",
    title: "Secure Diagnostic Reports",
    desc: "Generate professional, print-ready PDF diagnostic reports complete with laboratory letterheads, barcode markings, and pathologist validation signatures.",
    items: [
      "Print-ready PDF exports",
      "Custom lab branding & headers",
      "Barcode identification tags",
      "Doctor & clinic carbon copy (CC) logs"
    ],
  },
  {
    icon: Lock,
    color: "text-pink-500 dark:text-pink-400",
    bg: "bg-pink-500/10",
    title: "HIPAA-Compliant Access Control",
    desc: "Enforce standard clinical roles (Laboratory Director, Manager, Pathologist, Technician, Phlebotomist) with strict segregation of system access.",
    items: [
      "Granular user permissions",
      "Workspace segregation",
      "Session expiration controls",
      "Audit trail identity validation"
    ],
  },
  {
    icon: History,
    color: "text-cyan-500 dark:text-cyan-400",
    bg: "bg-cyan-500/10",
    title: "Regulatory Audit Logging",
    desc: "Maintain absolute compliance. Every reference range change, manual result override, and diagnostic release is logged in a tamper-proof audit trail.",
    items: [
      "CAP/CLIA compliant logs",
      "Analyte modification history",
      "User activity trails",
      "Compliance report exporting"
    ],
  },
];

export function Features() {
  return (
    <section id="features" className="py-24 px-6 border-t border-border bg-muted/5">
      <div className="max-w-7xl mx-auto">
        <div className="text-center max-w-2xl mx-auto mb-16">
          <h2 className="text-4xl font-black mb-4 text-foreground">
            Every workflow aligned with{" "}
            <span className="brand-gradient">clinical guidelines</span>
          </h2>
          <p className="text-muted-foreground text-lg">
            We built Curexal to support pure laboratory operations, focusing entirely on diagnostic precision, regulatory audit paths, and staff workspace isolation.
          </p>
        </div>

        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          {features.map((f, i) => {
            const Icon = f.icon;
            return (
              <div
                key={f.title}
                className="rounded-2xl border border-border bg-card p-6 glow-card animate-fade-in-up"
                style={{ animationDelay: `${i * 80}ms` }}
              >
                <div
                  className={`w-10 h-10 rounded-xl ${f.bg} flex items-center justify-center mb-4`}
                >
                  <Icon className={`h-5 w-5 ${f.color}`} />
                </div>
                <h3 className="font-bold text-foreground mb-2">{f.title}</h3>
                <p className="text-sm text-muted-foreground mb-4 leading-relaxed">{f.desc}</p>
                <ul className="space-y-1.5">
                  {f.items.map((item) => (
                    <li key={item} className="flex items-center gap-2 text-xs text-muted-foreground/80">
                      <Check className="h-3 w-3 text-primary flex-shrink-0" />
                      {item}
                    </li>
                  ))}
                </ul>
              </div>
            );
          })}
        </div>
      </div>
    </section>
  );
}
