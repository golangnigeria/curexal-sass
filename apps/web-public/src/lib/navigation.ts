export function isPatientSubdomain(): boolean {
  if (typeof window === "undefined") return false;
  return window.location.hostname.startsWith("patient");
}

export function ensurePatientDomain(targetPath?: string): boolean {
  if (typeof window === "undefined") return false;
  if (isPatientSubdomain()) return false;
  if (targetPath) {
    window.location.href = targetPath;
    return true;
  }
  return false;
}

export function ensurePortalDomain(targetPath?: string): boolean {
  if (typeof window === "undefined") return false;
  if (targetPath) {
    window.location.href = targetPath;
    return true;
  }
  return false;
}
