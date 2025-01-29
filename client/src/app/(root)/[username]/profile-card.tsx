import type { DetailedUser } from "@/types/users";
import styles from "./profile.module.scss";
import Image from "next/image";
import FollowButton from "@/components/user/follow-button";

interface ProfileCardProps {
  user: DetailedUser;
}

export default function ProfileCard({ user }: ProfileCardProps) {
  return (
    <section className={styles.userCardContainer}>
      <div className={styles.userCardInfo}>
        <div className={styles.userCardAvatar}>
          <Image
            src="/photo.jpg"
            alt={`${user.username} avatar`}
            width={128}
            height={128}
          />
        </div>

        <h3 className={styles.userCardName}>{user.username}</h3>
        <p className={styles.userCardDetails}>{user.email}</p>
        <p className={styles.userCardSubscription}>
          {user.subscription.followers} followers ·{" "}
          {user.subscription.following} following
        </p>
      </div>

      <FollowButton size="large" isFollowing={user.subscription.isFollowing} />
    </section>
  );
}
