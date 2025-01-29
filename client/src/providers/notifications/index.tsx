import { getNotificationsStats } from "@/shared/actions/notifications";
import NotificationStatsProvider from "./notifications-provider";

export default async function NotificationStatsClientWrapper({
  children,
}: React.PropsWithChildren) {
  const stats = await getNotificationsStats().catch(() => null);
  console.log("stats", stats);

  return (
    <NotificationStatsProvider stats={stats} dirty>
      {children}
    </NotificationStatsProvider>
  );
}
