import { HTTPError, type ResponsePromise } from "ky";
import type { ApiException } from "@/types/api";
import type { ReadableObject, SafeResult } from "@/types";

export async function safeSerialize<T extends ReadableObject>(
  res: ResponsePromise<T>,
): Promise<SafeResult<{ data: T }, { error: ApiException }>> {
  try {
    const body = await res.json<T>();
    return { success: true, data: body };
  } catch (err) {
    const fetchErr = await FetchError.fromKy(err);
    return { success: false, error: fetchErr };
  }
}

export class FetchError implements ApiException {
  error: string;
  status: number;
  statusText: string;

  static async fromKy(err: unknown): Promise<FetchError> {
    try {
      if (!(err instanceof HTTPError)) throw err;
      const body = await err.response.json<object>();

      const errorMessage =
        "error" in body && typeof body.error === "string"
          ? body.error
          : "Something went wrong, please wait a moment and try again.";

      const status =
        "status" in body && typeof body.status === "number"
          ? body.status
          : err.response.status;

      const statusText =
        "statusText" in body && typeof body.statusText === "string"
          ? body.statusText
          : err.response.statusText;

      return new FetchError(errorMessage, status, statusText);
    } catch (_) {
      return FetchError.unhandled();
    }
  }

  static unhandled() {
    return new FetchError(
      "Unhandled error Occurred",
      500,
      "Internal Server Error",
    );
  }

  constructor(error: string, status: number, statusText: string) {
    this.error = error;
    this.status = status;
    this.statusText = statusText;
  }
}
