/**
 * Unified API Error Taxonomy for Curexal
 */

export class ApiError extends Error {
  public statusCode: number;
  public code?: string;
  public details?: any;
  public requestId?: string;

  constructor(
    message: string,
    statusCode: number = 500,
    options?: { code?: string; details?: any; requestId?: string }
  ) {
    super(message);
    this.name = "ApiError";
    this.statusCode = statusCode;
    this.code = options?.code;
    this.details = options?.details;
    this.requestId = options?.requestId;
  }
}

export function formatErrorMessage(
  err: unknown,
  fallback: string = "An unexpected error occurred"
): string {
  if (err instanceof ApiError) return err.message;
  if (err instanceof Error) return err.message;
  if (typeof err === "string") return err;
  return fallback;
}
