import { Link } from "react-router-dom";
import { Activity, Twitter, Linkedin, Github, Mail } from "lucide-react";

const footerNav = {
  Solutions: [
    { label: "Laboratory LIMS", href: "/solutions#lims" },
    { label: "Clinic EMR", href: "/solutions#emr" },
    { label: "Patient Portal", href: "/solutions#patient" },
    { label: "Diagnostic Marketplace", href: "/marketplace" },
  ],
  Platform: [
    { label: "Specimen Tracking", href: "/solutions#lims" },
    { label: "Pathologist Validation", href: "/solutions#lims" },
    { label: "Audit Logs & RBAC", href: "/solutions#security" },
    { label: "Multi-Tenant Architecture", href: "/solutions#platform" },
  ],
  Resources: [
    { label: "Documentation", href: "/developers" },
    { label: "API Reference", href: "/developers" },
    { label: "Release Notes", href: "/resources#releases" },
    { label: "Compliance Center", href: "/resources#compliance" },
  ],
  Company: [
    { label: "About", href: "/about" },
    { label: "Careers", href: "/about#careers" },
    { label: "Blog", href: "/resources#blog" },
    { label: "Contact", href: "/about#contact" },
  ],
  Legal: [
    { label: "Privacy Policy", href: "#" },
    { label: "Terms of Service", href: "#" },
    { label: "Data Processing", href: "#" },
    { label: "HIPAA Notice", href: "#" },
  ],
};

const complianceBadges = [
  { label: "ISO 15189 Aligned", color: "#0F766E" },
  { label: "HIPAA Privacy Controls", color: "#0F766E" },
  { label: "NDPR Data Protection", color: "#0F766E" },
  { label: "AES-256 Encrypted", color: "#6B7280" },
];

export function MarketingFooter() {
  return (
    <footer className="bg-[#F8FAFC] dark:bg-[#0B1120] border-t border-gray-200 dark:border-[#1F2937]">
      <div className="max-w-[1280px] mx-auto px-6 pt-16 pb-10">

        {/* Main Grid */}
        <div className="grid grid-cols-2 md:grid-cols-7 gap-10 pb-12 border-b border-gray-200 dark:border-[#1F2937]">

          {/* Brand: spans 2 columns */}
          <div className="col-span-2">
            <Link to="/" className="flex items-center gap-2 mb-4 group">
              <img src="/logo-symbol.svg" alt="Curexal Logo" className="w-8 h-8 rounded-xl shadow-sm group-hover:scale-105 transition-transform" />
              <span className="font-bold text-[15px] tracking-tight text-gray-900 dark:text-white">
                Curexal
              </span>
            </Link>
            <p className="text-sm text-gray-500 dark:text-gray-400 leading-relaxed max-w-[240px]">
              The connected healthcare operating network for African laboratories, clinics, pharmacies, and patients.
            </p>
            <div className="flex items-center gap-2.5 mt-6">
              {[
                { icon: Twitter, href: "#", label: "Twitter" },
                { icon: Linkedin, href: "#", label: "LinkedIn" },
                { icon: Github, href: "#", label: "GitHub" },
                { icon: Mail, href: "mailto:hello@curexal.com", label: "Email" },
              ].map(({ icon: Icon, href, label }) => (
                <a
                  key={label}
                  href={href}
                  aria-label={label}
                  className="w-8 h-8 rounded-[8px] border border-gray-200 dark:border-[#374151] flex items-center justify-center text-gray-400 hover:text-gray-900 dark:hover:text-gray-100 hover:border-gray-300 dark:hover:border-[#4B5563] transition-colors"
                >
                  <Icon className="h-3.5 w-3.5" />
                </a>
              ))}
            </div>
          </div>

          {/* Nav columns */}
          {Object.entries(footerNav).map(([group, links]) => (
            <div key={group}>
              <h4 className="text-xs font-semibold text-gray-400 dark:text-gray-500 uppercase tracking-widest mb-4">
                {group}
              </h4>
              <ul className="space-y-2.5">
                {links.map((link) => (
                  <li key={link.label}>
                    <Link
                      to={link.href}
                      className="text-sm text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100 transition-colors"
                    >
                      {link.label}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>

        {/* Bottom bar */}
        <div className="pt-8 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-5">
          <p className="text-xs text-gray-400 dark:text-gray-500">
            © {new Date().getFullYear()} Curexal Health Technologies. All rights reserved.
          </p>

          {/* Compliance badges */}
          <div className="flex items-center gap-2 flex-wrap">
            {complianceBadges.map((badge) => (
              <span
                key={badge.label}
                className="inline-flex items-center px-2.5 py-1 rounded-[6px] text-[11px] font-semibold text-[#0F766E] bg-[#F0FDFA] dark:bg-[#0F766E]/10 border border-[#0F766E]/20"
              >
                {badge.label}
              </span>
            ))}
            <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-[6px] text-[11px] font-medium text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-[#1F2937] border border-gray-200 dark:border-[#374151]">
              <span className="w-1.5 h-1.5 rounded-full bg-green-500" />
              All systems operational
            </span>
          </div>
        </div>

      </div>
    </footer>
  );
}
