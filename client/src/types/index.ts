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
