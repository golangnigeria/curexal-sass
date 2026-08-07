import { useEffect } from "react";

interface SEOHeadProps {
  title: string;
  description?: string;
  ogType?: string;
  canonicalUrl?: string;
}

export function SEOHead({
  title,
  description = "Curexal — Enterprise Healthcare Platform & Public Marketplace. Seamlessly connect laboratories, clinics, hospitals, pharmacies, and medical vendors.",
  ogType = "website",
  canonicalUrl,
}: SEOHeadProps) {
  useEffect(() => {
    // Update Title
    const fullTitle = title.includes("Curexal") ? title : `${title} | Curexal Healthcare Platform`;
    document.title = fullTitle;

    // Update Meta Description
    let metaDescription = document.querySelector('meta[name="description"]');
    if (!metaDescription) {
      metaDescription = document.createElement("meta");
      metaDescription.setAttribute("name", "description");
      document.head.appendChild(metaDescription);
    }
    metaDescription.setAttribute("content", description);

    // Update OG Title
    let ogTitle = document.querySelector('meta[property="og:title"]');
    if (!ogTitle) {
      ogTitle = document.createElement("meta");
      ogTitle.setAttribute("property", "og:title");
      document.head.appendChild(ogTitle);
    }
    ogTitle.setAttribute("content", fullTitle);

    // Update OG Description
    let ogDesc = document.querySelector('meta[property="og:description"]');
    if (!ogDesc) {
      ogDesc = document.createElement("meta");
      ogDesc.setAttribute("property", "og:description");
      document.head.appendChild(ogDesc);
    }
    ogDesc.setAttribute("content", description);

    // Update Canonical URL
    if (canonicalUrl) {
      let canonical = document.querySelector('link[rel="canonical"]');
      if (!canonical) {
        canonical = document.createElement("link");
        canonical.setAttribute("rel", "canonical");
        document.head.appendChild(canonical);
      }
      canonical.setAttribute("href", canonicalUrl);
    }
  }, [title, description, ogType, canonicalUrl]);

  return null;
}
