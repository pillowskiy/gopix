import type { SafeResult } from "..";

export type SafeActionState<T extends object> = SafeResult<
  ActionSuccess<T>,
  ActionError
> | null;

export interface ActionError {
  message: string;
  errors?: Record<string, string | undefined>;
}

export type ActionSuccess<T extends object> = {
  data: T;
};
