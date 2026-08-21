# Curexal Enterprise Design System Specification

> **Purpose**: Authoritative design system specification defining default Curexal theme tokens, typography, grid, components, and accessible color palettes.  
> **Owner**: Lead UI/UX Designer  
> **Status**: APPROVED / PRODUCTION READY  
> **Last Updated**: 2026-07-27  
> **Verification Criteria**: Implemented in `@curexal/design-tokens` and `@curexal/ui-core`.

---

## 1. Global Default Brand Color Palette

```text
Primary Color:     #266210  (Forest Deep Green)
Secondary Accent:  #90B800  (Lime Emerald Accent)
Status Accent:     #00E1E1  (Cyan Medical Success)
Background Dark:   #063B00  (Deep Emerald Black)
Surface Dark:      #042B00  (Dark Card Background)
```

> **Note on Tenant Branding Overrides**: Tenant customer organizations may override primary colors via workspace branding settings, but the above values represent Curexal's global default theme.

---

## 2. Spacing & Typography System

- **Font Family**: `Inter, system-ui, -apple-system, sans-serif`
- **Monospace Font**: `JetBrains Mono, Fira Code, monospace` (for DDL schemas, IPs, ULIDs)
- **Baseline Spacing System**: 8px grid (`4px`, `8px`, `12px`, `16px`, `24px`, `32px`, `48px`)
- **Border Radii**: 12px for metric cards, 16px for view containers, 24px for modals.
- **Shadows**: Soft dark elevation shadows (`shadow-lg`, `shadow-2xl`).

---

## 3. Required Component States

Every component across `@curexal/ui-core` and `@curexal/ui-healthcare` must explicitly support:
1. **Loading State**: Smooth skeleton placeholders matching exact card/table dimensions to prevent Cumulative Layout Shifts (CLS).
2. **Empty State**: Clean iconography and actionable call-to-action prompts when zero records are returned.
3. **Error State**: RFC 7807 problem details viewer with a "Retry Connection" action button.
4. **Success State**: Clear visual confirmation toasts or status badges.
5. **Permission Denied State**: Lock icon with Casbin role requirements text.
