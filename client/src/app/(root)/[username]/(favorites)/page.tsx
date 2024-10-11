import NothingHere from './nothing-here';

export default async function ProfilePage({ params }: { params: Record<string, string> }) {
	return <NothingHere username={params.username} />;
}
