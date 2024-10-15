import type { Paginated } from "@/types";
import type { ImageWithAuthor } from "@/types/images";
import { $api } from "../api-interceptor";

export async function getFavoriteImages(
  userId: string,
): Promise<Paginated<ImageWithAuthor>> {
  return $api
    .get(`images/favorites/${userId}`, {
      searchParams: { limit: 50, page: 1, sort: "popular" },
    })
    .then((res) => res.json<Paginated<ImageWithAuthor>>());
}
