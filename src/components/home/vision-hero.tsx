import React, { useState } from "react";
import { Link } from "react-router-dom";
import {
  Sparkles,
  ShieldCheck,
  Search,
  MapPin,
  Clock,
  Building2,
  CheckCircle2,
  ArrowRight,
  FileCheck,
  Stethoscope,
  TestTube,
  Lock,
  Layers,
  Zap,
  Globe2,
  Calendar,
  ChevronRight,
  Home,
  Star,
  Users,
} from "lucide-react";

export function VisionHero() {
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedCity, setSelectedCity] = useState("all");
  const [selectedCategory, setSelectedCategory] = useState("all");
  const [bookedLab, setBookedLab] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<"marketplace" | "vision" | "ecosystem">("marketplace");

  const labPartners = [
    {
      id: "lab_apollo",
      name: "Apollo Diagnostics Central Lab",
      accreditation: "ISO 15189 Accredited",
      location: "Victoria Island, Lagos",
      rating: "4.9 ★★★★★ (240 reviews)",
      tests: [
        { name: "Full Blood Count (FBC)", price: 7500, tat: "4 Hours", category: "hematology" },
        { name: "Lipid Panel Profile", price: 12500, tat: "6 Hours", category: "biochemistry" },
        { name: "Comprehensive Metabolic Panel (CMP)", price: 18000, tat: "8 Hours", category: "biochemistry" },
        { name: "HbA1c Diabetes Profile", price: 9500, tat: "5 Hours", category: "biochemistry" },
      ],
      homeCollectionAvailable: true,
    },
    {
      id: "lab_everight",
      name: "Everight Pathology Diagnostic Network",
      accreditation: "MLSCN Certified",
      location: "Ikeja GRA, Lagos",
      rating: "4.8 ★★★★★ (185 reviews)",
      tests: [
        { name: "Full Blood Count (FBC)", price: 6800, tat: "5 Hours", category: "hematology" },
        { name: "Liver Function Test (LFT)", price: 14000, tat: "6 Hours", category: "biochemistry" },
        { name: "Thyroid Stimulating Hormone (TSH)", price: 15500, tat: "12 Hours", category: "endocrinology" },
        { name: "Prostate Specific Antigen (PSA)", price: 16000, tat: "24 Hours", category: "oncology" },
      ],
      homeCollectionAvailable: true,
    },
    {
      id: "lab_lifecare",
      name: "LifeCare Specialist Diagnostic Centre",
      accreditation: "ISO 15189 Accredited",
      location: "Abuja Central Business District",
      rating: "4.9 ★★★★★ (310 reviews)",
      tests: [
        { name: "Full Blood Count (FBC)", price: 8000, tat: "3 Hours", category: "hematology" },
        { name: "Full Lipid & Cardiac Risk Profile", price: 19500, tat: "6 Hours", category: "biochemistry" },
        { name: "Vitamin D (25-OH) Assay", price: 28000, tat: "24 Hours", category: "endocrinology" },
        { name: "Full Sex Hormone Panel", price: 32000, tat: "24 Hours", category: "endocrinology" },
      ],
      homeCollectionAvailable: false,
    },
  ];

  const filteredLabs = labPartners.filter((lab) => {
    const matchesSearch =
      lab.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      lab.location.toLowerCase().includes(searchQuery.toLowerCase()) ||
      lab.tests.some((t) => t.name.toLowerCase().includes(searchQuery.toLowerCase()));

    const matchesCity =
      selectedCity === "all" ||
      (selectedCity === "lagos" && lab.location.includes("Lagos")) ||
      (selectedCity === "abuja" && lab.location.includes("Abuja"));

    const matchesCategory =
      selectedCategory === "all" ||
      lab.tests.some((t) => t.category === selectedCategory);

    return matchesSearch && matchesCity && matchesCategory;
  });

  return (
    <section className="relative pt-24 pb-20 bg-slate-950 text-slate-100 overflow-hidden font-sans border-b border-slate-800">
      {/* Background Lighting Elements */}
      <div className="absolute inset-0 pointer-events-none">
        <div className="absolute -top-40 left-1/2 -translate-x-1/2 w-[1000px] h-[500px] rounded-full bg-emerald-500/10 blur-[140px]" />
        <div className="absolute top-1/3 left-10 w-96 h-96 bg-teal-500/5 blur-[120px] rounded-full" />
        <div className="absolute bottom-10 right-10 w-80 h-80 bg-cyan-500/5 blur-[100px] rounded-full" />
      </div>

      <div className="relative z-10 max-w-7xl mx-auto px-6 space-y-12">
        {/* Top Hero Heading & Grand Vision Banner */}
        <div className="text-center max-w-4xl mx-auto space-y-6">
          <div className="inline-flex items-center gap-2 px-3.5 py-1.5 rounded-full text-xs font-extrabold bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 shadow-lg shadow-emerald-500/10">
            <Sparkles className="w-3.5 h-3.5 text-emerald-400" />
            The Unified Infrastructure for Pathology Labs, Hospitals & Patients
          </div>

          <h1 className="text-4xl sm:text-6xl font-black tracking-tight text-white leading-[1.08]">
            Connecting Diagnostic Labs, <br />
            Clinicians & Patients <br />
            <span className="bg-gradient-to-r from-emerald-400 via-teal-300 to-cyan-400 bg-clip-text text-transparent">
              In One Interoperable Health Cloud
            </span>
          </h1>

          <p className="text-base sm:text-lg text-slate-300 max-w-3xl mx-auto leading-relaxed">
            Curexal delivers end-to-end multi-tenant Laboratory Information Systems (LIS) for diagnostic networks, paired with a public digital diagnostic marketplace for patients to search accredited labs, compare test prices, and access electronically signed PDF reports in their private Patient Vault.
          </p>

          {/* Interactive Mode Selector */}
          <div className="inline-flex p-1.5 rounded-2xl bg-slate-900/90 border border-slate-800 shadow-2xl">
            <button
              onClick={() => setActiveTab("marketplace")}
              className={`px-5 py-2.5 rounded-xl text-xs font-bold transition-all cursor-pointer ${
                activeTab === "marketplace"
                  ? "bg-emerald-500 text-slate-950 shadow-md"
                  : "text-slate-400 hover:text-slate-200"
              }`}
            >
              🌐 Patient Diagnostic Marketplace
            </button>
            <button
              onClick={() => setActiveTab("vision")}
              className={`px-5 py-2.5 rounded-xl text-xs font-bold transition-all cursor-pointer ${
                activeTab === "vision"
                  ? "bg-emerald-500 text-slate-950 shadow-md"
                  : "text-slate-400 hover:text-slate-200"
              }`}
            >
              🔬 The Grand Vision & LIS Pillars
            </button>
            <button
              onClick={() => setActiveTab("ecosystem")}
              className={`px-5 py-2.5 rounded-xl text-xs font-bold transition-all cursor-pointer ${
                activeTab === "ecosystem"
                  ? "bg-emerald-500 text-slate-950 shadow-md"
                  : "text-slate-400 hover:text-slate-200"
              }`}
            >
              🏥 Healthcare Ecosystem Network
            </button>
          </div>
        </div>

        {/* TAB 1: PATIENT DIAGNOSTIC MARKETPLACE */}
        {activeTab === "marketplace" && (
          <div className="space-y-8 animate-fade-in">
            {/* Search Box */}
            <div className="max-w-4xl mx-auto bg-slate-900/90 border border-slate-800 p-4 rounded-2xl shadow-2xl space-y-4">
              <div className="flex flex-col md:flex-row gap-3">
                <div className="relative flex-1">
                  <Search className="absolute left-4 top-3.5 w-5 h-5 text-slate-400" />
                  <input
                    type="text"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    placeholder="Search diagnostic test (e.g. Full Blood Count, Lipid Profile, HbA1c, PSA)..."
                    className="w-full pl-12 pr-4 py-3 bg-slate-950/80 border border-slate-800 rounded-xl text-sm text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 font-medium"
                  />
                </div>

                <div className="flex gap-2">
                  <select
                    value={selectedCity}
                    onChange={(e) => setSelectedCity(e.target.value)}
                    className="px-4 py-3 bg-slate-950/80 border border-slate-800 rounded-xl text-xs font-bold text-slate-200 focus:outline-none focus:ring-2 focus:ring-emerald-500"
                  >
                    <option value="all">All Locations</option>
                    <option value="lagos">Lagos State</option>
                    <option value="abuja">Abuja FCT</option>
                  </select>

                  <select
                    value={selectedCategory}
                    onChange={(e) => setSelectedCategory(e.target.value)}
                    className="px-4 py-3 bg-slate-950/80 border border-slate-800 rounded-xl text-xs font-bold text-slate-200 focus:outline-none focus:ring-2 focus:ring-emerald-500"
                  >
                    <option value="all">All Specialities</option>
                    <option value="hematology">Hematology</option>
                    <option value="biochemistry">Biochemistry</option>
                    <option value="endocrinology">Endocrinology</option>
                    <option value="oncology">Oncology</option>
                  </select>
                </div>
              </div>

              <div className="flex flex-wrap items-center justify-between gap-3 pt-2 text-xs text-slate-400">
                <div className="flex items-center gap-2">
                  <ShieldCheck className="w-4 h-4 text-emerald-400" /> Verified Laboratory Accreditation & Compliance Tracking
                </div>
                <div className="flex items-center gap-4">
                  <span className="text-slate-300 font-semibold">Showing {filteredLabs.length} Partner Labs</span>
                </div>
              </div>
            </div>

            {/* Lab Directory Grid */}
            <div className="grid grid-cols-1 gap-6">
              {filteredLabs.map((lab) => (
                <div
                  key={lab.id}
                  className="bg-slate-900/80 border border-slate-800 rounded-2xl p-6 shadow-xl hover:border-slate-700 transition-all space-y-6"
                >
                  <div className="flex flex-col md:flex-row justify-between gap-4 border-b border-slate-800 pb-4">
                    <div>
                      <div className="flex items-center gap-2">
                        <h3 className="text-lg font-bold text-white">{lab.name}</h3>
                        <span className="px-2.5 py-0.5 rounded-full text-xs font-semibold bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
                          {lab.accreditation}
                        </span>
                      </div>
                      <div className="flex items-center gap-4 text-xs text-slate-400 mt-2">
                        <span className="flex items-center gap-1">
                          <MapPin className="w-3.5 h-3.5 text-slate-500" /> {lab.location}
                        </span>
                        <span className="text-amber-400 font-medium">{lab.rating}</span>
                      </div>
                    </div>

                    {lab.homeCollectionAvailable && (
                      <span className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-teal-500/10 text-teal-300 border border-teal-500/20 text-xs font-medium self-start">
                        <Home className="w-3.5 h-3.5 text-teal-400" /> Home Sample Phlebotomy Available
                      </span>
                    )}
                  </div>

                  {/* Test Cards */}
                  <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                    {lab.tests.map((t, idx) => (
                      <div
                        key={idx}
                        className="p-4 rounded-xl bg-slate-950/80 border border-slate-800 space-y-3 flex flex-col justify-between"
                      >
                        <div>
                          <h4 className="text-xs font-bold text-slate-200 line-clamp-1">{t.name}</h4>
                          <div className="flex items-center gap-1 text-[11px] text-slate-400 mt-1">
                            <Clock className="w-3 h-3 text-emerald-400" /> TAT: {t.tat}
                          </div>
                        </div>

                        <div className="pt-2 border-t border-slate-800 flex justify-between items-baseline">
                          <span className="text-[11px] text-slate-400">Price</span>
                          <span className="text-sm font-extrabold text-emerald-400">₦{t.price.toLocaleString()}</span>
                        </div>
                      </div>
                    ))}
                  </div>

                  {/* Requisition Button */}
                  <div className="flex flex-col sm:flex-row justify-between items-center gap-4 pt-2">
                    <div className="text-xs text-slate-400 flex items-center gap-2">
                      <FileCheck className="w-4 h-4 text-emerald-400" /> Digitally signed PDF report delivered to your Patient Vault
                    </div>

                    <button
                      onClick={() => setBookedLab(lab.name)}
                      className="w-full sm:w-auto px-6 py-2.5 rounded-xl bg-emerald-500 hover:bg-emerald-400 text-slate-950 font-bold text-xs transition-all shadow-lg shadow-emerald-500/10 flex items-center justify-center gap-2 cursor-pointer"
                    >
                      <Calendar className="w-4 h-4" /> Book Lab Appointment
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* TAB 2: GRAND VISION & LIS PILLARS */}
        {activeTab === "vision" && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 animate-fade-in">
            <div className="p-6 rounded-2xl bg-slate-900/80 border border-slate-800 space-y-4 hover:border-emerald-500/40 transition-colors">
              <div className="w-12 h-12 rounded-xl bg-emerald-500/10 text-emerald-400 flex items-center justify-center border border-emerald-500/20">
                <TestTube className="w-6 h-6" />
              </div>
              <h3 className="text-base font-bold text-white">Pathology & LIMS Core</h3>
              <p className="text-xs text-slate-400 leading-relaxed">
                Automated specimen accessioning, barcode label generation, analyzer auto-ingestion, and multi-tier pathologist verification.
              </p>
            </div>

            <div className="p-6 rounded-2xl bg-slate-900/80 border border-slate-800 space-y-4 hover:border-emerald-500/40 transition-colors">
              <div className="w-12 h-12 rounded-xl bg-teal-500/10 text-teal-400 flex items-center justify-center border border-teal-500/20">
                <Stethoscope className="w-6 h-6" />
              </div>
              <h3 className="text-base font-bold text-white">Clinical EMR Interoperability</h3>
              <p className="text-xs text-slate-400 leading-relaxed">
                Connect hospital electronic health record systems with diagnostic labs for seamless electronic lab orders and instant result dispatch.
              </p>
            </div>

            <div className="p-6 rounded-2xl bg-slate-900/80 border border-slate-800 space-y-4 hover:border-emerald-500/40 transition-colors">
              <div className="w-12 h-12 rounded-xl bg-cyan-500/10 text-cyan-400 flex items-center justify-center border border-cyan-500/20">
                <Globe2 className="w-6 h-6" />
              </div>
              <h3 className="text-base font-bold text-white">Patient Health Vault</h3>
              <p className="text-xs text-slate-400 leading-relaxed">
                Empower patients with portable electronic health records, test history analytics, and verified PDF reports accessible anytime.
              </p>
            </div>

            <div className="p-6 rounded-2xl bg-slate-900/80 border border-slate-800 space-y-4 hover:border-emerald-500/40 transition-colors">
              <div className="w-12 h-12 rounded-xl bg-indigo-500/10 text-indigo-400 flex items-center justify-center border border-indigo-500/20">
                <Lock className="w-6 h-6" />
              </div>
              <h3 className="text-base font-bold text-white">ISO & CLIA Compliance</h3>
              <p className="text-xs text-slate-400 leading-relaxed">
                Cryptographic digital signatures, audit logs, role-based access controls, and complete specimen chain of custody tracking.
              </p>
            </div>
          </div>
        )}

        {/* TAB 3: ECOSYSTEM NETWORK */}
        {activeTab === "ecosystem" && (
          <div className="bg-slate-900/80 border border-slate-800 rounded-3xl p-8 space-y-8 animate-fade-in">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
              <div className="space-y-4">
                <span className="px-3 py-1 rounded-full text-[11px] font-extrabold bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
                  For Pathology Laboratories
                </span>
                <h3 className="text-lg font-bold text-white">Automate Lab Throughput</h3>
                <ul className="space-y-2 text-xs text-slate-400">
                  <li className="flex items-center gap-2"><CheckCircle2 className="w-4 h-4 text-emerald-400" /> Direct interface with automated analyzers</li>
                  <li className="flex items-center gap-2"><CheckCircle2 className="w-4 h-4 text-emerald-400" /> Phlebotomy & accessioning worklists</li>
                  <li className="flex items-center gap-2"><CheckCircle2 className="w-4 h-4 text-emerald-400" /> Electronic pathologist stamp & sign-off</li>
                </ul>
              </div>

              <div className="space-y-4">
                <span className="px-3 py-1 rounded-full text-[11px] font-extrabold bg-teal-500/10 text-teal-400 border border-teal-500/20">
                  For Clinics & Hospitals
                </span>
                <h3 className="text-lg font-bold text-white">Streamline Patient Orders</h3>
                <ul className="space-y-2 text-xs text-slate-400">
                  <li className="flex items-center gap-2"><CheckCircle2 className="w-4 h-4 text-teal-400" /> Electronic lab test ordering</li>
                  <li className="flex items-center gap-2"><CheckCircle2 className="w-4 h-4 text-teal-400" /> Real-time status tracking & notifications</li>
                  <li className="flex items-center gap-2"><CheckCircle2 className="w-4 h-4 text-teal-400" /> Auto-population of EMR clinical charts</li>
                </ul>
              </div>

              <div className="space-y-4">
                <span className="px-3 py-1 rounded-full text-[11px] font-extrabold bg-cyan-500/10 text-cyan-400 border border-cyan-500/20">
                  For Patients & Families
                </span>
                <h3 className="text-lg font-bold text-white">Take Control of Your Health</h3>
                <ul className="space-y-2 text-xs text-slate-400">
                  <li className="flex items-center gap-2"><CheckCircle2 className="w-4 h-4 text-cyan-400" /> Compare accredited lab prices</li>
                  <li className="flex items-center gap-2"><CheckCircle2 className="w-4 h-4 text-cyan-400" /> Book home sample phlebotomy visits</li>
                  <li className="flex items-center gap-2"><CheckCircle2 className="w-4 h-4 text-cyan-400" /> Download verified PDF lab reports</li>
                </ul>
              </div>
            </div>
          </div>
        )}

        {/* Live Metrics Row */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-6 pt-6 border-t border-slate-800/80">
          <div className="text-center space-y-1">
            <span className="text-2xl sm:text-3xl font-black text-emerald-400">4.2M+</span>
            <p className="text-xs text-slate-400 font-medium">Specimen Accessionings</p>
          </div>
          <div className="text-center space-y-1">
            <span className="text-2xl sm:text-3xl font-black text-teal-400">120+</span>
            <p className="text-xs text-slate-400 font-medium">Accredited Pathology Labs</p>
          </div>
          <div className="text-center space-y-1">
            <span className="text-2xl sm:text-3xl font-black text-cyan-400">&lt;4 Hours</span>
            <p className="text-xs text-slate-400 font-medium">Average Report TAT</p>
          </div>
          <div className="text-center space-y-1">
            <span className="text-2xl sm:text-3xl font-black text-indigo-400">99.98%</span>
            <p className="text-xs text-slate-400 font-medium">Platform Reliability</p>
          </div>
        </div>
      </div>

      {/* Appointment Modal Confirmation */}
      {bookedLab && (
        <div className="fixed inset-0 bg-slate-950/80 backdrop-blur-md flex items-center justify-center p-4 z-50">
          <div className="bg-slate-900 border border-slate-800 rounded-3xl p-6 max-w-md w-full shadow-2xl space-y-5 animate-fade-in">
            <div className="w-12 h-12 rounded-2xl bg-emerald-500/20 text-emerald-400 flex items-center justify-center border border-emerald-500/30">
              <CheckCircle2 className="w-6 h-6" />
            </div>
            <div className="space-y-2">
              <h3 className="text-lg font-bold text-white">Diagnostic Visit Requested</h3>
              <p className="text-xs text-slate-300 leading-relaxed">
                You have requested a diagnostic test appointment with <strong className="text-emerald-400">{bookedLab}</strong>. Our clinical coordinator will reach out shortly to confirm your appointment time and home phlebotomy details.
              </p>
            </div>
            <button
              onClick={() => setBookedLab(null)}
              className="w-full py-3 rounded-xl bg-emerald-500 hover:bg-emerald-400 text-slate-950 font-bold text-xs transition-colors cursor-pointer"
            >
              Close Confirmation
            </button>
          </div>
        </div>
      )}
    </section>
  );
}
