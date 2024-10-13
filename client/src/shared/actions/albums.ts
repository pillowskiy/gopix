import type { DetailedAlbum } from "@/types/albums";
import { $api } from "../api-interceptor";

export async function getUserAlbums(userId: string): Promise<DetailedAlbum[]> {
  return $api
    .get(`albums/users/${userId}`)
    .then((res) => res.json<DetailedAlbum[]>());
}
