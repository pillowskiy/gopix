import { NotificationStats } from "@/types/notifications";
import { $api } from "../api-interceptor";

export async function getNotificationsStats(): Promise<NotificationStats> {
  return $api.get("notifications/stats").json<NotificationStats>();
}
