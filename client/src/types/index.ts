export interface PaginationInput {
  limit: number;
  page: number;
}

export interface Pagination extends PaginationInput {
  total: number;
}

export interface Paginated<T extends object> extends Pagination {
  data: T[];
}

export type SafeResult<S extends ReadableObject, E extends ReadableObject> =
  | Success<S>
  | Failure<E>;

export type Failure<T extends ReadableObject> = { success: false } & T;
export type Success<T extends ReadableObject> = { success: true } & T;

// biome-ignore lint/suspicious/noExplicitAny: <explanation>
export type ReadableObject = Record<string, any>;
