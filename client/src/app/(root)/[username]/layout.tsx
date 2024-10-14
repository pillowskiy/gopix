import { LightGradient } from '@/components/backdroung-texture';
import Section from '@/components/section';
import ProfileTabs from './profile-tabs';
import ProfileCard from './profile-card';
import getByUsernameCached from './getByUsernameCached';

export default async function ProfileLayout({
	params,
	children
}: { params: Record<string, string>; children: React.ReactElement }) {
	const user = await getByUsernameCached(params.username);

	return (
		<>
			<LightGradient />
			<ProfileCard user={user} />
			<Section>
				<ProfileTabs username={user.username} />
				{children}
			</Section>
		</>
	);
}
