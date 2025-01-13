import type { Paginated } from "@/types";
import type { ImageWithMeta } from "@/types/images";
import { $api } from "../api-interceptor";

export enum ImageSort {
  Popular = "popular",
  Newest = "newest",
  Oldest = "oldest",
  MostViewed = "mostViewed",
}

export async function getFavoriteImages(
  userId: string,
): Promise<Paginated<ImageWithMeta>> {
  return $api
    .get(`images/favorites/${userId}`, {
      searchParams: { limit: 50, page: 1, sort: "popular" },
    })
    .then((res) => res.json<Paginated<ImageWithMeta>>());
}

export async function getDiscoverImages(
  page: number,
  limit: number,
  sort: ImageSort,
): Promise<Paginated<ImageWithMeta>> {
  return $api
    .get(`images/`, { searchParams: { limit, page, sort } })
    .then((res) => res.json<Paginated<ImageWithMeta>>());
}
