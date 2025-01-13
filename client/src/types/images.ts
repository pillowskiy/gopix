import type { Tag } from "./tags";
import type { PublicUser } from "./users";

export interface Image {
  id: string;
  path: string;
  title: string;
  description: string;
  accessLevel: string;
  createdAt: string;
}

export interface ImageProperties {
  ext: string;
  width: number;
  height: number;
}

export interface ImageWithMeta extends Image {
  properties: ImageProperties;
  author: PublicUser;
}

export interface DetailedImage extends ImageWithMeta {
  likes: number;
  views: number;
  tags: Tag[];
}
