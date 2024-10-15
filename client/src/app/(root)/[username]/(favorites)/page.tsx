import { Suspense } from 'react';
import NothingHere from './nothing-here';
import { ErrorBoundary } from 'react-error-boundary';
import FavoritesSection from './favorites-section';
import getByUsernameCached from '../getByUsernameCached';

export default async function ProfilePage({ params }: { params: Record<string, string> }) {
	const user = await getByUsernameCached(params.username);

	return (
		<Suspense fallback={<div />}>
			<ErrorBoundary fallback={<NothingHere username={user.username} />}>
				<FavoritesSection userId={user.id} />
			</ErrorBoundary>
		</Suspense>
	);
}
