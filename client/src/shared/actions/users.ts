import type { DetailedUser, User } from "@/types/users";
import { $api } from "../api-interceptor";

export async function getMe(): Promise<User> {
  return $api.get("users/@me").json<User>();
}

export async function getByUsername(username: string): Promise<DetailedUser> {
  return $api.get(`users/${username}`).json<DetailedUser>();
}
