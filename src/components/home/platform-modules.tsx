import {
  FlaskConical,
  Stethoscope,
  RadioTower,
  Package,
  CreditCard,
  BarChart3,
  ShoppingBag,
  User,
} from "lucide-react";

const modules = [
  {
    icon: FlaskConical,
    label: "Laboratory LIMS",
    desc: "Specimen tracking, analyzer integration, pathologist sign-off, connected to every referring clinic on the network",
    available: true,
  },
  {
    icon: Stethoscope,
    label: "Clinic EMR",
    desc: "Digital referrals, real-time result receipt, clinical notes, with no paper or phone calls",
    available: true,
  },
  {
    icon: RadioTower,
    label: "Radiology",
    desc: "Imaging orders, DICOM routing, radiologist reporting, coming soon to the network",
    available: false,
  },
  {
    icon: Package,
    label: "Inventory",
    desc: "Reagent management, expiry tracking, stock alerts across your organization",
    available: true,
  },
  {
    icon: CreditCard,
    label: "Billing & Settlements",
    desc: "Test pricing, invoicing, referral commissions, and payment reconciliation",
    available: true,
  },
  {
    icon: BarChart3,
    label: "Network Analytics",
    desc: "Test volumes, TAT trends, quality metrics across your entire network",
    available: true,
  },
  {
    icon: ShoppingBag,
    label: "Diagnostic Marketplace",
    desc: "Patients find labs, compare prices, book tests, and receive results digitally",
    available: true,
  },
  {
    icon: User,
    label: "Patient Vault",
    desc: "Secure health records, verified PDF reports, appointment history, owned by the patient",
    available: true,
  },
];

export function PlatformModules() {
  return (
    <section
      id="platform"
      className="section-padding bg-white dark:bg-[#0B1120]"
    >
      <div className="max-w-[1280px] mx-auto px-6">

        {/* Header */}
        <div className="max-w-2xl mb-12">
          <div className="accent-line mb-4" />
          <h2 className="text-section text-gray-900 dark:text-white mb-4">
            A Complete Operating Network.<br />Not Just Another App.
          </h2>
          <p className="text-body text-gray-500 dark:text-gray-400">
            Most competitors build hospital software, or lab software, or a patient portal, or a marketplace. Curexal builds all four as one connected ecosystem, because healthcare doesn't stop at the hospital gate.
          </p>
        </div>

        {/* 4×2 grid */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          {modules.map((mod) => {
            const Icon = mod.icon;
            return (
              <div
                key={mod.label}
                className="card-enterprise p-5 hover-lift group"
              >
                <div className="flex items-start justify-between mb-4">
                  <div className="w-9 h-9 rounded-[10px] bg-[#F0FDFA] dark:bg-[#0F766E]/10 flex items-center justify-center group-hover:bg-[#0F766E]/15 transition-colors">
                    <Icon className="h-4.5 w-4.5 text-[#0F766E]" strokeWidth={1.75} />
                  </div>
                  {!mod.available && (
                    <span className="text-[10px] font-semibold px-2 py-0.5 rounded-full bg-gray-100 dark:bg-[#1F2937] text-gray-400 dark:text-gray-500 border border-gray-200 dark:border-[#374151]">
                      Coming Soon
                    </span>
                  )}
                </div>
                <h3 className="text-[14px] font-semibold text-gray-900 dark:text-white mb-1.5">
                  {mod.label}
                </h3>
                <p className="text-[13px] text-gray-500 dark:text-gray-400 leading-relaxed">
                  {mod.desc}
                </p>
              </div>
            );
          })}
        </div>

      </div>
    </section>
  );
}
