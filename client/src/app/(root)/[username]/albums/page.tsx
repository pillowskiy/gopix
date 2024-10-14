import { Suspense } from 'react';
import { UserAlbums, UserAlbumsFallback } from './user-albums-section';
import getByUsernameCached from '../getByUsernameCached';

export default async function UserAlbumsPage({ params }: { params: Record<string, string> }) {
	const user = await getByUsernameCached(params.username);

	return (
		<Suspense fallback={<UserAlbumsFallback />}>
			<UserAlbums userId={user.id} />
		</Suspense>
	);
}
