import { Suspense } from 'react';
import { getByUsername } from '@/shared/users';
import { UserAlbums, UserAlbumsFallback } from './user-albums-section';

export default async function UserAlbumsPage({ params }: { params: Record<string, string> }) {
	const user = await getByUsername(params.username);

	return (
		<Suspense fallback={<UserAlbumsFallback />}>
			<UserAlbums userId={user.id} />
		</Suspense>
	);
}
