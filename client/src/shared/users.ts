import type { User } from "@/types/users";
import { $sapi } from "./drill-interceptor";

export async function getMe(): Promise<User> {
  return $sapi.get("users/@me").json<User>();
}
