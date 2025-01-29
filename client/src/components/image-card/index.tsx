import type { ImageWithMeta } from "@/types/images";
import { CDNImage } from "../image/cdn-image";
import styles from "./image-card.module.scss";
import cc from "classcat";

interface ImageCardProps extends React.ComponentProps<"div"> {
  image: ImageWithMeta;
  withAuthor?: boolean;
}

export function ImageCard({ className, image, ...props }: ImageCardProps) {
  return (
    <div className={cc([styles.card, className])} {...props}>
      <div
        style={{
          // @ts-expect-error
          "--aspect-ratio": image.properties.width / image.properties.height,
        }}
        className={styles.cardImage}
      >
        <CDNImage
          path={image.path}
          alt={image.title || `Image ${image.path}`}
          width={image.properties.width}
          height={image.properties.height}
        />
      </div>

      <div className={styles.cardActions}>
        <div className={styles.cardActionsAuthor}>
          <div className={styles.cardActionsAuthorAvatar}>
            <img src={"/photo.jpg"} alt={`${image.author.username}'s avatar`} />
          </div>
          <span>{image.author.username}</span>
        </div>
      </div>
    </div>
  );
}
