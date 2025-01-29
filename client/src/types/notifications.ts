export interface Notification {
  title: string;
  message?: string;
  hidden: boolean;
  read: boolean;
  sentAt: string;
}

export interface NotificationStats {
  unread: number;
}
