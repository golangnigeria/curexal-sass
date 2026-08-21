import { useState, useEffect } from "react";
import { Link, useLocation } from "react-router-dom";
import { cn } from "@/lib/utils";
import {
  Activity,
  Menu,
  X,
  ChevronDown,
  ChevronRight,
  Sparkles,
  Calendar,
} from "lucide-react";
import { WaitlistModal } from "@/components/waitlist-modal";

interface NavDropdownItem {
  label: string;
  href: string;
  description: string;
}

const solutionsItems: NavDropdownItem[] = [
  { label: "Laboratory LIMS", href: "/solutions#lims", description: "Specimen tracking & automated results delivery" },
  { label: "Clinic EMR", href: "/solutions#emr", description: "Clinical ordering & patient electronic records" },
  { label: "Healthcare Marketplace", href: "/marketplace", description: "Public search for labs, clinics, pharmacies & supply vendors" },
];

export function MarketingNavbar() {
  const [scrolled, setScrolled] = useState(false);
  const [mobileOpen, setMobileOpen] = useState(false);
  const [solutionsOpen, setSolutionsOpen] = useState(false);
  const [waitlistOpen, setWaitlistOpen] = useState(false);
  const location = useLocation();
  const pathname = location.pathname;

  useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 8);
    window.addEventListener("scroll", onScroll, { passive: true });
    return () => window.removeEventListener("scroll", onScroll);
  }, []);

  // Close mobile menu on route change
  useEffect(() => {
    setMobileOpen(false);
    setSolutionsOpen(false);
  }, [pathname]);

  // Lock body scroll when mobile menu is open
  useEffect(() => {
    if (mobileOpen) {
      document.body.style.overflow = "hidden";
    } else {
      document.body.style.overflow = "";
    }
    return () => {
      document.body.style.overflow = "";
    };
  }, [mobileOpen]);

  return (
    <>
      <header
        className={cn(
          "fixed top-0 left-0 right-0 z-[9999] transition-all duration-200",
          scrolled
            ? "bg-white/98 dark:bg-[#0B1120]/98 backdrop-blur-xl border-b border-slate-200/80 dark:border-slate-800 shadow-md"
            : "bg-white/90 dark:bg-[#0B1120]/90 backdrop-blur-md border-b border-slate-100 dark:border-slate-800/50"
        )}
      >
        <div className="max-w-[1280px] mx-auto px-3 sm:px-6 h-14 sm:h-16 flex items-center justify-between gap-2">
          {/* Brand Logo */}
          <div className="flex items-center gap-2 flex-shrink-0">
            <Link
              to="/"
              className="flex items-center gap-2 group flex-shrink-0 cursor-pointer"
              id="marketing-logo"
            >
              <img src="/logo-symbol.svg" alt="Curexal Logo" className="w-8 h-8 rounded-xl shadow-sm group-hover:scale-105 transition-transform" />
              <span className="font-extrabold text-base tracking-tight text-slate-900 dark:text-white">
                Curexal
              </span>
            </Link>

          </div>

          {/* Desktop Navigation Links */}
          <nav className="hidden lg:flex items-center gap-6 xl:gap-8">
            <Link
              to="/"
              className={cn(
                "text-sm font-medium transition-colors hover:text-[#0F766E]",
                pathname === "/" ? "text-[#0F766E] font-bold" : "text-slate-600 dark:text-slate-300"
              )}
            >
              Home
            </Link>

            {/* Desktop Solutions Dropdown */}
            <div
              className="relative"
              onMouseEnter={() => setSolutionsOpen(true)}
              onMouseLeave={() => setSolutionsOpen(false)}
            >
              <Link
                to="/solutions"
                className={cn(
                  "text-sm font-medium transition-colors flex items-center gap-1 py-2 hover:text-[#0F766E]",
                  pathname.startsWith("/solutions") ? "text-[#0F766E] font-bold" : "text-slate-600 dark:text-slate-300"
                )}
              >
                <span>Solutions</span>
                <ChevronDown className={cn("w-3.5 h-3.5 transition-transform duration-200", solutionsOpen && "rotate-180")} />
              </Link>

              {solutionsOpen && (
                <div className="absolute top-full left-0 w-80 p-2 rounded-2xl bg-white dark:bg-[#0B1120] border border-slate-200/80 dark:border-slate-800 shadow-2xl animate-in fade-in slide-in-from-top-2 duration-150 z-[99999]">
                  {solutionsItems.map((item) => (
                    <Link
                      key={item.href}
                      to={item.href}
                      onClick={() => setSolutionsOpen(false)}
                      className="block p-3 rounded-xl hover:bg-teal-50/60 dark:hover:bg-teal-950/40 transition-colors group"
                    >
                      <p className="text-sm font-bold text-slate-900 dark:text-white group-hover:text-[#0F766E]">
                        {item.label}
                      </p>
                      <p className="text-xs text-slate-500 dark:text-slate-400 mt-0.5">
                        {item.description}
                      </p>
                    </Link>
                  ))}
                </div>
              )}
            </div>

            <Link
              to="/marketplace"
              className={cn(
                "text-sm font-medium transition-colors hover:text-[#0F766E]",
                pathname.startsWith("/marketplace") ? "text-[#0F766E] font-bold" : "text-slate-600 dark:text-slate-300"
              )}
            >
              Marketplace
            </Link>

            <Link
              to="/pricing"
              className={cn(
                "text-sm font-medium transition-colors hover:text-[#0F766E]",
                pathname === "/pricing" ? "text-[#0F766E] font-bold" : "text-slate-600 dark:text-slate-300"
              )}
            >
              Pricing
            </Link>

            <Link
              to="/about"
              className={cn(
                "text-sm font-medium transition-colors hover:text-[#0F766E]",
                pathname === "/about" ? "text-[#0F766E] font-bold" : "text-slate-600 dark:text-slate-300"
              )}
            >
              About
            </Link>
          </nav>

          {/* Desktop Right Actions */}
          <div className="hidden lg:flex items-center gap-3">
            <Link to="/waitlist">
              <button
                className="px-4 py-2 rounded-xl bg-teal-50 dark:bg-teal-950/40 text-[#0F766E] dark:text-teal-300 hover:bg-teal-100 dark:hover:bg-teal-900/60 border border-teal-200/80 dark:border-teal-800/80 text-xs font-bold transition-all cursor-pointer shadow-xs"
              >
                Join Waitlist
              </button>
            </Link>

            <Link to="/book-demo" id="nav-book-demo">
              <button className="px-5 py-2 rounded-xl bg-[#0F766E] hover:bg-[#115E59] text-white text-xs font-bold transition-all shadow-sm cursor-pointer border-0">
                Book Demo
              </button>
            </Link>
          </div>

          {/* Mobile Header Controls */}
          <div className="flex items-center gap-1.5 lg:hidden flex-shrink-0">
            <Link to="/waitlist">
              <button
                className="px-2.5 py-1 rounded-lg bg-teal-50 dark:bg-teal-950/60 text-[#0F766E] dark:text-teal-300 border border-teal-200/80 dark:border-teal-800/80 text-[11px] font-bold transition-all cursor-pointer flex items-center gap-1"
              >
                <Sparkles className="w-3 h-3 text-[#0F766E]" />
                <span>Waitlist</span>
              </button>
            </Link>

            <Link to="/book-demo">
              <button className="px-2.5 py-1 rounded-lg bg-[#0F766E] text-white text-[11px] font-bold border-0 cursor-pointer">
                Demo
              </button>
            </Link>

            <button
              className="p-1.5 rounded-lg text-slate-800 dark:text-slate-100 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors border-0 cursor-pointer flex-shrink-0"
              onClick={() => setMobileOpen(!mobileOpen)}
              aria-label="Toggle Navigation Menu"
            >
              {mobileOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
            </button>
          </div>
        </div>
      </header>

      {/* COMPACT MOBILE MENU OVERLAY (No Scrolling Needed!) */}
      {mobileOpen && (
        <div className="lg:hidden fixed top-14 sm:top-16 inset-x-0 bg-white dark:bg-[#0B1120] z-[99999] p-4 shadow-xl border-b border-slate-200 dark:border-slate-800 animate-in fade-in slide-in-from-top-2 duration-150">
          <div className="flex flex-col gap-1">
            {[
              { to: "/", label: "Home" },
              { to: "/solutions", label: "Solutions" },
              { to: "/marketplace", label: "Marketplace" },
              { to: "/pricing", label: "Pricing" },
              { to: "/about", label: "About" },
            ].map((link) => {
              const active = pathname === link.to || (link.to !== "/" && pathname.startsWith(link.to));
              return (
                <Link
                  key={link.to}
                  to={link.to}
                  onClick={() => setMobileOpen(false)}
                  className={cn(
                    "flex items-center justify-between py-2 px-3 rounded-lg text-sm font-semibold transition-all",
                    active
                      ? "bg-teal-50 dark:bg-teal-950/60 text-[#0F766E] dark:text-teal-400 font-bold"
                      : "text-slate-800 dark:text-slate-200 hover:bg-slate-50 dark:hover:bg-slate-900"
                  )}
                >
                  <span>{link.label}</span>
                  <ChevronRight className="w-4 h-4 opacity-30" />
                </Link>
              );
            })}
          </div>

          {/* Compact Bottom Action Buttons */}
          <div className="pt-3 mt-2 border-t border-slate-100 dark:border-slate-800 grid grid-cols-2 gap-2">
            <button
              onClick={() => {
                setMobileOpen(false);
                setWaitlistOpen(true);
              }}
              className="py-2.5 rounded-lg bg-teal-50 dark:bg-teal-950/60 text-[#0F766E] dark:text-teal-300 font-bold border border-teal-200 dark:border-teal-800 text-xs flex items-center justify-center gap-1 cursor-pointer"
            >
              <Sparkles className="w-3.5 h-3.5" />
              <span>Waitlist</span>
            </button>

            <Link to="/book-demo" onClick={() => setMobileOpen(false)}>
              <button className="w-full py-2.5 rounded-lg bg-[#0F766E] text-white text-xs font-bold border-0 flex items-center justify-center gap-1 cursor-pointer">
                <Calendar className="w-3.5 h-3.5" />
                <span>Book Demo</span>
              </button>
            </Link>
          </div>
        </div>
      )}

      {/* Waitlist Modal */}
      <WaitlistModal open={waitlistOpen} onOpenChange={setWaitlistOpen} />
    </>
  );
}
