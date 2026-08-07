const stats = [
  { value: "100k+", label: "Specimens Logged" },
  { value: "99.9%", label: "Result Accuracy" },
  { value: "< 4 Hours", label: "Average Turnaround" },
  { value: "CLIA", label: "Compliant Audit Trails" },
];

export function Stats() {
  return (
    <section className="py-20 px-6">
      <div className="max-w-4xl mx-auto grid grid-cols-2 md:grid-cols-4 gap-8 text-center">
        {stats.map((stat) => (
          <div key={stat.label}>
            <div className="text-4xl font-black brand-gradient mb-2">{stat.value}</div>
            <div className="text-sm text-muted-foreground">{stat.label}</div>
          </div>
        ))}
      </div>
    </section>
  );
}
