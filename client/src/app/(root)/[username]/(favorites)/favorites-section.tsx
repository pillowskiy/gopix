import { ImageCard } from '@/components/image-card';
import styles from './favorites-page.module.scss';
import { getFavoriteImages } from '@/shared/actions/images';

interface FavoritesSectionProps {
	userId: string;
}

export default async function FavoritesSection({ userId }: FavoritesSectionProps) {
	const favorites = await getFavoriteImages(userId);

	if (!favorites.items.length) {
		throw new Error('Zero items provided');
	}

	return (
		<section className={styles.favoritesSection}>
			{favorites.items.map((image) => (
				<ImageCard key={image.id} image={image} />
			))}
		</section>
	);
}
