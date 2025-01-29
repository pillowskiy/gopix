"use client";

import { Button } from "../ui/button";
import UserDropdown from "../user/user-dropdown";
import { CursorArrowRaysIcon } from "@heroicons/react/20/solid";
import styles from "./header.module.scss";
import { useUserStore } from "@/providers/auth/store";
import Link from "next/link";
import NotificationDropdown from "../notification/notification-dropdown";

export default function HeaderActions() {
  const { data: user } = useUserStore();

  return (
    <div className={styles.headerSection}>
      <div className={styles.headerActions}>
        <Button size="icon" variant="ghost">
          <CursorArrowRaysIcon />
        </Button>
        {user && <NotificationDropdown />}
      </div>

      {user ? (
        <UserDropdown user={user} />
      ) : (
        <Button size="small" href="/login" as={Link}>
          Log in
        </Button>
      )}
    </div>
  );
}
