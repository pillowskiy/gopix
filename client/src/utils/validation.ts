import { ZodError } from "zod";

export function prettifyZodErrors(errors: ZodError): Record<string, string> {
  const fieldErrors = errors.flatten().fieldErrors;

  const prettifiedErrors: Record<string, string> = {};
  for (const [key, error] of Object.entries(fieldErrors)) {
    if (!error) continue;
    prettifiedErrors[key] = error.toString();
  }

  return prettifiedErrors;
}
