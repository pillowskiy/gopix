"use client";

import { Button, ButtonProps } from "@/components/ui/button";
import { useUserStore } from "@/providers/auth/store";

interface FollowButtonProps extends ButtonProps {
  isFollowing: boolean;
}

export default function FollowButton({
  isFollowing,
  ...props
}: FollowButtonProps) {
  const user = useUserStore();

  if (!user) {
    return (
      <Button {...props} disabled>
        Follow
      </Button>
    );
  }

  return (
    <Button variant={isFollowing ? "secondary" : "accent"} {...props}>
      {isFollowing ? "Unfollow" : "Follow"}
    </Button>
  );
}
