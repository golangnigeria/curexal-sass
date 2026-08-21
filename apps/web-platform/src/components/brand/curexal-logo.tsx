import React from "react";
import { cn } from "@/lib/utils";

interface LogoProps extends React.SVGProps<SVGSVGElement> {
  className?: string;
  showText?: boolean;
}

export function CurexalLogoSymbol({ className, ...props }: React.SVGProps<SVGSVGElement>) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 500 500"
      fill="none"
      className={cn("w-8 h-8 shrink-0", className)}
      {...props}
    >
      {/* Light Mint / Cyan Connecting Network Struts */}
      <g stroke="#56E3B8" strokeWidth="14" strokeLinecap="round" strokeLinejoin="round">
        <line x1="250" y1="105" x2="375" y2="323" />
        <line x1="375" y1="323" x2="142" y2="198" />
        <line x1="125" y1="323" x2="352" y2="192" />
        <line x1="125" y1="323" x2="238" y2="368" />
      </g>

      {/* Dark Hexagonal Perimeter Struts */}
      <g stroke="#1E2538" strokeWidth="26" strokeLinecap="butt">
        <line x1="222" y1="120" x2="148" y2="163" />
        <line x1="278" y1="120" x2="352" y2="163" />
        <line x1="125" y1="205" x2="125" y2="295" />
        <line x1="375" y1="205" x2="375" y2="295" />
        <line x1="148" y1="337" x2="222" y2="380" />
        <line x1="352" y1="337" x2="278" y2="380" />
      </g>

      {/* Central Iris / Precision Core */}
      <g>
        <circle cx="250" cy="250" r="70" stroke="#00DF9B" strokeWidth="14" fill="#0F172A" />
        <circle cx="250" cy="250" r="42" fill="#2563EB" />
        <circle cx="250" cy="250" r="17" fill="#FFFFFF" />
      </g>

      {/* 3 Turquoise Spherical Nodes */}
      <g fill="#00DF9B">
        <circle cx="250" cy="105" r="30" />
        <circle cx="375" cy="323" r="30" />
        <circle cx="125" cy="323" r="30" />
      </g>
    </svg>
  );
}

export function CurexalLogo({ className, showText = true }: LogoProps) {
  return (
    <div className={cn("flex items-center gap-3 select-none", className)}>
      <div className="relative flex items-center justify-center p-1 rounded-xl bg-slate-900 shadow-md">
        <CurexalLogoSymbol className="w-8 h-8" />
      </div>
      {showText && (
        <div className="flex flex-col">
          <span className="font-extrabold tracking-tight text-foreground text-base leading-tight">
            CUREXAL
          </span>
          <span className="text-[10px] tracking-wider font-semibold text-primary uppercase">
            Platform Console
          </span>
        </div>
      )}
    </div>
  );
}
