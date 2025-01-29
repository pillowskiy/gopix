"use client";

import { Button } from "@/components/ui/button";
import { BellIcon } from "@heroicons/react/24/outline";
import styles from "./notification-dropdown.module.scss";
import { useNotificationStats } from "@/providers/notifications/store";

export default function NotificationDropdown() {
  const { stats } = useNotificationStats();

  if (!stats) {
    return null;
  }

  return (
    <Button className={styles.trigger} size="icon" variant="ghost">
      <BellIcon />
      <span className={styles.triggerUnreadCounter}>{stats.unread}</span>
    </Button>
  );
}
