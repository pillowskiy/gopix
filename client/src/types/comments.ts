import type { PublicUser } from "./users";

export interface Comment {
  id: string;
  text: string;
  createdAt: string;
  updatedAt: string;
}

export interface CommentWithAuthor extends Comment {
  author: PublicUser;
}

export interface DetailedComment extends CommentWithAuthor {
  stats: CommentStats;
}

export interface CommentStats {
  repliesCount: number;
  likes: number;
  liked: boolean;
}
