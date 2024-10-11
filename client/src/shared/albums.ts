import type { DetailedAlbum } from "@/types/albums";
import { $sapi } from "./drill-interceptor";

export async function getUserAlbums(userId: string): Promise<DetailedAlbum[]> {
  return $sapi
    .get(`albums/users/${userId}`)
    .then((res) => res.json<DetailedAlbum[]>());
}
