import { ArrowUpTrayIcon } from "@heroicons/react/16/solid";
import { AttentionButton } from "../ui/button";
import styles from "./header.module.scss";
import HeaderNav from "./header-nav";
import Link from "next/link";
import HeaderActions from "./header-actions";

export default function Header() {
  return (
    <header className={styles.header}>
      <div className={styles.headerContent}>
        <div className={styles.headerSection}>
          <Link href="/" className={styles.headerLogo}>
            GoPix
          </Link>
          <AttentionButton
            style={{ display: "flex", alignItems: "center", gap: "8px" }}
          >
            Upload <ArrowUpTrayIcon style={{ height: "16px", width: "16px" }} />
          </AttentionButton>
        </div>

        <HeaderNav />

        <HeaderActions />
      </div>
    </header>
  );
}
