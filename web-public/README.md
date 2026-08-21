# CUREXAL_WEB-PUBLIC

Enterprise Healthcare Platform, Diagnostic Laboratory Operating Network & Healthcare Marketplace.

## 🚀 Overview

`CUREXAL_WEB-PUBLIC` is the public web platform and customer research engine for Curexal. It enables healthcare leads, laboratory system directors, clinic managers, medical suppliers, and patients to discover services, join the early access waitlist, vote on roadmap feature priorities, and browse accredited diagnostic facilities.

---

## 🛠️ Technology Stack

- **Framework**: React 19 + TypeScript + Vite
- **Styling**: TailwindCSS + Vanilla CSS Design System
- **Database & Lead Storage**: Supabase (`@supabase/supabase-js`)
- **State & Data Fetching**: TanStack Query (React Query)
- **UI Components**: Shadcn UI + Radix Primitives + Lucide Icons + Sonner Toasts
- **Routing**: React Router v7

---

## 💻 Local Development

### 1. Install Dependencies
```bash
bun install
```

### 2. Environment Variables
Copy `.env.example` to `.env`:
```bash
cp .env.example .env
```

Ensure the following variables are configured:
```env
VITE_SUPABASE_URL=https://xrqwupliqiuotzbkcrja.supabase.co
VITE_SUPABASE_PUBLISHABLE_KEY=sb_publishable_-yewa_hvqkBJQ8fuJ9nSCQ_sc0J-_iC
VITE_PATIENT_PORTAL_URL=https://patient.curexal.com
VITE_PORTAL_URL=https://app.curexal.com
```

### 3. Run Development Server
```bash
bun run dev
```

### 4. Typecheck & Build
```bash
bun run typecheck
bun run build
```

---

## 📄 License
Copyright © 2026 Curexal Inc. All rights reserved.
