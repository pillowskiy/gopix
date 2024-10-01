import type { Image } from "./images";
import type { PublicUser } from "./users";

export interface Album {
  id: string;
  name: string;
  description: string;
  createdAt: string;
  updatedAt: string;
}

export interface AlbumWithAuthor extends Album {
  author: PublicUser;
}

export interface DetailedAlbum extends AlbumWithAuthor {
  cover: Image[];
}
