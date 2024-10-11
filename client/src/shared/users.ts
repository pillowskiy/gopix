import type { DetailedUser, User } from "@/types/users";
import { $sapi } from "./drill-interceptor";

export async function getMe(): Promise<User> {
  return $sapi.get("users/@me").json<User>();
}

export async function getByUsername(username: string): Promise<DetailedUser> {
  return $sapi.get(`users/${username}`).json<DetailedUser>();
}
