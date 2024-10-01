export interface PublicUser {
  id: string;
  username: string;
  avatarURL: string;
}

export interface User extends PublicUser {
  email: string;
  permissions: string;
  createdAt: string;
  updatedAt: string;
}

export interface DetailedUser extends User {
  subscription: Subscription;
}

export interface Subscription {
  followers: number;
  following: number;
  isFollowing: boolean;
}
