import { MarketingNavbar } from "@/components/layouts/marketing-navbar";
import { MarketingFooter } from "@/components/layouts/marketing-footer";
import { Link } from "react-router-dom";
import { Code2, Webhook, Key, Layers, ArrowRight } from "lucide-react";

const categories = [
  {
    icon: Code2,
    label: "REST API Reference",
    desc: "Full API documentation for all Curexal platform endpoints. Authentication, request/response schemas, error handling, and pagination.",
    href: "/book-demo",
    cta: "Request API Access",
  },
  {
    icon: Key,
    label: "Authentication",
    desc: "JWT-based session authentication, API key management, OAuth 2.0 integration guide, and role scoped token issuance.",
    href: "/book-demo",
    cta: "Auth Specs",
  },
  {
    icon: Webhook,
    label: "Webhooks",
    desc: "Subscribe to real-time platform events (specimen received, result validated, report dispatched) via HTTPS webhooks with HMAC signature verification.",
    href: "/book-demo",
    cta: "Webhook Reference",
  },
  {
    icon: Layers,
    label: "SDKs & Integrations",
    desc: "Official TypeScript and Python SDKs for the Curexal API. Integration guides for EMR systems, HL7 FHIR, and analyzer interfaces.",
    href: "/book-demo",
    cta: "SDK Specs",
  },
];

const codeSnippet = `// Fetch lab results for a patient
const response = await curexal.results.list({
  patientId: "pat_01HXYZ",
  status: "validated",
  limit: 20,
});

console.log(response.data);
// [{ id, testName, result, validatedAt, signedPdfUrl }]`;

export function DevelopersPage() {
  return (
    <div className="min-h-screen bg-white dark:bg-[#0B1120] font-sans">
      <MarketingNavbar />

      {/* Header */}
      <div className="pt-28 pb-16 bg-[#F8FAFC] dark:bg-[#0B1120] border-b border-gray-100 dark:border-[#1F2937]">
        <div className="max-w-[1280px] mx-auto px-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
            <div>
              <div className="accent-line mb-4" />
              <h1 className="text-hero text-gray-900 dark:text-white mb-5">
                Build on<br />Curexal.
              </h1>
              <p className="text-body text-gray-500 dark:text-gray-400 mb-8">
                RESTful APIs, real-time webhooks, and SDKs to integrate diagnostic laboratory workflows into your own applications.
              </p>
              <Link to="/book-demo">
                <button
                  className="flex items-center gap-2 px-5 py-2.5 rounded-[10px] bg-[#0F766E] hover:bg-[#115E59] text-white text-sm font-semibold transition-colors cursor-pointer border-0"
                  style={{ boxShadow: "0 1px 2px rgba(0,0,0,0.1)" }}
                >
                  Request API Access
                  <ArrowRight className="h-4 w-4" />
                </button>
              </Link>
            </div>

            {/* Code block */}
            <div
              className="rounded-[14px] bg-[#111827] border border-[#374151] overflow-hidden"
              style={{ boxShadow: "0 1px 2px rgba(0,0,0,0.12), 0 8px 32px rgba(0,0,0,0.15)" }}
            >
              <div className="flex items-center gap-2 px-5 py-3 border-b border-[#374151]">
                <div className="flex gap-1.5">
                  <div className="w-3 h-3 rounded-full bg-[#374151]" />
                  <div className="w-3 h-3 rounded-full bg-[#374151]" />
                  <div className="w-3 h-3 rounded-full bg-[#374151]" />
                </div>
                <span className="text-xs text-[#6B7280] font-mono ml-2">results.ts</span>
              </div>
              <pre className="px-6 py-5 text-[13px] font-mono text-[#9CA3AF] leading-relaxed overflow-x-auto">
                <code>{codeSnippet}</code>
              </pre>
            </div>
          </div>
        </div>
      </div>

      {/* API categories */}
      <main className="max-w-[1280px] mx-auto px-6 py-16">
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-5">
          {categories.map((cat) => {
            const Icon = cat.icon;
            return (
              <Link key={cat.label} to={cat.href} className="block group">
                <div className="card-enterprise p-6 hover-lift h-full">
                  <div className="w-10 h-10 rounded-[10px] bg-[#F0FDFA] dark:bg-[#0F766E]/10 flex items-center justify-center mb-5 group-hover:bg-[#0F766E]/15 transition-colors">
                    <Icon className="h-5 w-5 text-[#0F766E]" strokeWidth={1.75} />
                  </div>
                  <h3 className="text-[15px] font-semibold text-gray-900 dark:text-white mb-2">{cat.label}</h3>
                  <p className="text-sm text-gray-500 dark:text-gray-400 leading-relaxed mb-5">{cat.desc}</p>
                  <span className="flex items-center gap-1.5 text-sm font-semibold text-[#0F766E] group-hover:gap-2.5 transition-all">
                    {cat.cta}
                    <ArrowRight className="h-4 w-4" />
                  </span>
                </div>
              </Link>
            );
          })}
        </div>
      </main>

      <MarketingFooter />
    </div>
  );
}
