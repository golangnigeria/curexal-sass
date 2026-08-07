import "./index.css";
import { LandingPage } from "@/pages/landing-page";
import { BookDemoPage } from "@/pages/book-demo";
import { PatientMarketplacePage } from "@/pages/patient-marketplace-page";
import { SolutionsPage } from "@/pages/solutions-page";
import { PricingPage } from "@/pages/pricing-page";
import { AboutPage } from "@/pages/about-page";
import { ResourcesPage } from "@/pages/resources-page";
import { DevelopersPage } from "@/pages/developers-page";
import { WaitlistPage } from "@/pages/waitlist-page";
import {
  Route,
  Navigate,
  createBrowserRouter,
  createRoutesFromElements,
} from "react-router-dom";
import { env } from "@/config/env";

const patientPortalUrl = env.VITE_PATIENT_PORTAL_URL || "https://patient.curexal.com";
const orgPortalUrl = env.VITE_PORTAL_URL || "https://app.curexal.com";

function ExternalRedirect({ targetUrl }: { targetUrl: string }) {
  window.location.href = targetUrl;
  return null;
}

const routes = createRoutesFromElements(
  <>
    {/* ── Public Landing ─────────────────────────────── */}
    <Route path="/" element={<LandingPage />} />

    {/* ── Public Information Architecture ───────────── */}
    <Route path="/solutions" element={<SolutionsPage />} />
    <Route path="/pricing" element={<PricingPage />} />
    <Route path="/about" element={<AboutPage />} />
    <Route path="/resources" element={<ResourcesPage />} />
    <Route path="/developers" element={<DevelopersPage />} />
    <Route path="/book-demo" element={<BookDemoPage />} />
    <Route path="/waitlist" element={<WaitlistPage />} />

    {/* ── Healthcare & Medical Marketplace ───────────── */}
    <Route path="/marketplace" element={<PatientMarketplacePage />} />

    {/* ── External App Redirects (Strict Public Constitution) ──── */}
    <Route path="/login" element={<ExternalRedirect targetUrl={`${patientPortalUrl}/login`} />} />
    <Route path="/register" element={<ExternalRedirect targetUrl={`${patientPortalUrl}/register`} />} />
    <Route path="/dashboard" element={<ExternalRedirect targetUrl={`${patientPortalUrl}/dashboard`} />} />

    {/* Catch-all fallback */}
    <Route path="*" element={<Navigate to="/" replace />} />
  </>
);

export const router = createBrowserRouter(routes);
