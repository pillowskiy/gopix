import { getByUsername } from '@/shared/actions/users';
import { LightGradient } from '@/components/backdroung-texture';
import Section from '@/components/section';
import ProfileTabs from './profile-tabs';
import ProfileCard from './profile-card';

export default async function ProfileLayout({
	params,
	children
}: React.PropsWithChildren<{ params: Record<string, string> }>) {
	const user = await getByUsername(params.username);

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
