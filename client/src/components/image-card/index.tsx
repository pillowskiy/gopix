import type { ImageWithAuthor } from '@/types/images';
import { CDNImage } from '../image/cdn-image';
import styles from './image-card.module.scss';
import cc from 'classcat';

interface ImageCardProps extends React.ComponentProps<'div'> {
	image: ImageWithAuthor;
	withAuthor?: boolean;
}

export function ImageCard({ className, image, ...props }: ImageCardProps) {
	return (
		<div className={cc([styles.card, className])} {...props}>
			<div className={styles.cardImage}>
				<CDNImage
					path={image.path}
					alt={image.title || `Image ${image.path}`}
					width={260}
					height={260}
				/>
			</div>
		</div>
	);
}
