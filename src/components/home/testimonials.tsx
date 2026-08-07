const testimonials = [
  {
    quote: "Before Curexal, we called clinicians on the phone to deliver results. Now, every referral, result, and commission is tracked automatically. Our turnaround time dropped 40%, and we haven't printed a single paper report in months.",
    name: "Dr. Adaeze Nwosu",
    title: "Medical Director, LifeCare Diagnostic Centre",
    org: "Private pathology network · Abuja, FCT",
  },
  {
    quote: "We joined the Curexal network because we were tired of fragmented systems. Now our clinicians receive lab results the moment they're signed off, with no chasing, no delays, and no lost reports. It's how healthcare should work.",
    name: "Mr. Emeka Okafor",
    title: "Chief Operations Officer, Meridian Health Group",
    org: "Multi-branch hospital group · Lagos State",
  },
  {
    quote: "The diagnostic marketplace changed everything for us. Patients find us, book tests, and receive verified PDF reports directly, all without a single phone call. We're connected to referring clinics we never had access to before.",
    name: "Mrs. Fatima Ibrahim",
    title: "CEO, Synlab West Africa",
    org: "Regional diagnostic laboratory · Pan-Nigeria",
  },
];

export function Testimonials() {
  return (
    <section
      id="testimonials"
      className="section-padding bg-[#F8FAFC] dark:bg-[#0B1120] border-y border-gray-100 dark:border-[#1F2937]"
    >
      <div className="max-w-[1280px] mx-auto px-6">

        <div className="max-w-xl mb-12">
          <div className="accent-line mb-4" />
          <h2 className="text-section text-gray-900 dark:text-white mb-4">
            Trusted by Healthcare Providers<br />Across the Network.
          </h2>
          <p className="text-body text-gray-500 dark:text-gray-400">
            From single-site diagnostic labs to national healthcare networks, providers join Curexal because connected healthcare works.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-5">
          {testimonials.map((t) => (
            <div key={t.name} className="card-enterprise p-6 hover-lift flex flex-col gap-5">
              {/* Quote mark */}
              <div className="text-[40px] leading-none text-[#0F766E]/30 font-serif select-none">&ldquo;</div>

              <blockquote className="text-[15px] text-gray-700 dark:text-gray-300 leading-relaxed flex-1">
                {t.quote}
              </blockquote>

              <div className="pt-4 border-t border-gray-100 dark:border-[#374151]">
                <p className="text-sm font-semibold text-gray-900 dark:text-white">{t.name}</p>
                <p className="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{t.title}</p>
                <p className="text-xs text-gray-400 dark:text-gray-500 mt-0.5">{t.org}</p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
