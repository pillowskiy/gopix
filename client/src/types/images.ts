import type { Tag } from "./tags";
import type { PublicUser } from "./users";

export interface Image {
  id: string;
  path: string;
  title: string;
  description: string;
  accessLevel: string;
  mime: string;
  ext: string;
  url: string;
  createdAt: string;
}

export interface ImageWithAuthor extends Image {
  author: PublicUser;
}

export interface DetailedImage extends ImageWithAuthor {
  likes: number;
  views: number;
  tags: Tag[];
}
