import type { Paginated } from "@/types";
import type { ImageWithAuthor } from "@/types/images";
import { $api } from "../api-interceptor";

export async function getFavoriteImages(
  userId: string,
): Promise<Paginated<ImageWithAuthor>> {
  return $api
    .get(`/api/images/favorites/${userId}`)
    .then((res) => res.json<Paginated<ImageWithAuthor>>());
}
