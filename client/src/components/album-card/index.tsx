import type { DetailedAlbum } from '@/types/albums';
import { Skeleton } from '../ui/skeleton';
import Image from 'next/image';
import cc from 'classcat';
import styles from './album-card.module.scss';
import Link from 'next/link';

interface AlbumCardProps {
	album: DetailedAlbum;
}

// TEMP: It was just easier to hard code than to play with :nth-child selectors
const cardCoverStyles = {
	1: styles.cardCoverSingle,
	2: styles.cardCoverDouble,
	3: styles.cardCoverTriple
};

const calcCoverStyle = (length: number) => {
	return cardCoverStyles[length as keyof typeof cardCoverStyles];
};

export function AlbumCard({ album }: AlbumCardProps) {
	return (
		<Link href={`/a/${album.id}`} className={styles.cardWrapper}>
			<div className={cc([styles.cardCover, calcCoverStyle(album.cover.length)])}>
				{album.cover.map((image) => (
					<div className={styles.cardCoverImage} key={image.id}>
						<Image src='/photo.jpg' alt={image.title} width={256} height={256} />
					</div>
				))}
			</div>

			<div className={styles.cardDetails}>
				<h4 className={styles.cardDetailsTitle}>{album.name}</h4>
				<p className={styles.cardDetailsDescription}>{album.description}</p>
			</div>
		</Link>
	);
}

export function AlbumCardSkeleton() {
	return (
		<div className={styles.cardWrapper}>
			<div className={cc([styles.cardCover, calcCoverStyle(1)])}>
				<Skeleton className={styles.cardCoverImage} />
			</div>

			<div className={styles.cardDetails}>
				<Skeleton style={{ width: `${Math.max(Math.random() * 70, 40)}%`, height: '24px' }} />
				<Skeleton style={{ width: `${Math.max(Math.random() * 100, 60)}%`, height: '16px' }} />
			</div>
		</div>
	);
}
