import { getUserAlbums } from '@/shared/albums';
import styles from './albums-page.module.scss';
import { AlbumCard, AlbumCardSkeleton } from '@/components/album-card';

export interface UserAlbumsProps {
	userId: string;
}

export async function UserAlbums({ userId }: UserAlbumsProps) {
	const albums = await getUserAlbums(userId);

	if (!albums) return null;

	return (
		<section className={styles.section}>
			{albums.map((album) => (
				<AlbumCard key={album.id} album={album} />
			))}
		</section>
	);
}

export function UserAlbumsFallback() {
	return (
		<section className={styles.section}>
			{Array.from({ length: 5 }).map((_, i) => (
				// biome-ignore lint/suspicious/noArrayIndexKey: <explanation>
				<AlbumCardSkeleton key={i} />
			))}
		</section>
	);
}
