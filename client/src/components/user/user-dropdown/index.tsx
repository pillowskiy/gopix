"use client";

import {
  DropdownMenu,
  DropdownMenuItem,
  DropdownMenuItems,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown";
import {
  ArrowRightStartOnRectangleIcon,
  Cog6ToothIcon,
  FolderIcon,
  HeartIcon,
  SunIcon,
  UserIcon,
} from "@heroicons/react/16/solid";
import Link from "next/link";
import Image from "next/image";
import {
  Button as HeadlessButton,
  ButtonProps as HeadlessButtonProps,
} from "@headlessui/react";
import { forwardRef } from "react";
import { User } from "@/types/users";
import styles from "./user.module.scss";

const UserDropdownTrigger = forwardRef<
  HTMLButtonElement,
  HeadlessButtonProps & { user: User }
>(({ user, ...props }, ref) => {
  return (
    <HeadlessButton ref={ref} aria-label="User Dropdown" {...props}>
      <Image
        className={styles.dropdownTriggerImage}
        src="/photo.jpg"
        alt={`${user.username} avatar`}
        width={48}
        height={48}
      />
    </HeadlessButton>
  );
});

export default function UserDropdown({ user }: { user: User }) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger
        className={styles.dropdownTrigger}
        user={user}
        as={UserDropdownTrigger}
      />
      <DropdownMenuItems anchor="bottom end" style={{ width: "206px" }}>
        <div style={{ padding: "12px" }}>
          <h5 style={{ fontSize: "16px" }}>{user.username}</h5>
          <p style={{ fontSize: "14px" }}>{user.email}</p>
        </div>

        <DropdownMenuSeparator />

        <DropdownMenuItem as={Link} href={`/${user.username}`}>
          <UserIcon className="icon" />
          Your Profile
        </DropdownMenuItem>

        <DropdownMenuItem as={Link} href={`/${user.username}/albums`}>
          <FolderIcon className="icon" />
          Your Albums
        </DropdownMenuItem>

        <DropdownMenuItem as={Link} href="/account/favorites">
          <HeartIcon className="icon" />
          Your Favorites
        </DropdownMenuItem>

        <DropdownMenuItem as={Link} href="/account/settings">
          <Cog6ToothIcon className="icon" />
          Settings
        </DropdownMenuItem>

        <DropdownMenuSeparator />
        <DropdownMenuItem>
          <SunIcon className="icon" />
          Light Mode
        </DropdownMenuItem>
        <DropdownMenuItem>
          <ArrowRightStartOnRectangleIcon className="icon" />
          Logout
        </DropdownMenuItem>
      </DropdownMenuItems>
    </DropdownMenu>
  );
}
