export function isPrimitive(
  value: unknown,
): value is string | number | boolean {
  return (
    (typeof value !== "object" && typeof value !== "function") || value === null
  );
}
