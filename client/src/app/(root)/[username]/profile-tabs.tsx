'use client';

import Link from 'next/link';
import styles from './profile.module.scss';
import { usePathname } from 'next/navigation';
import { useRef } from 'react';

function buildProfileTabs(username: string): Readonly<{ name: string; href: string }[]> {
	return [
		{
			name: 'Favorites',
			href: `/${username}`
		},
		{
			name: 'Albums',
			href: `/${username}/albums`
		}
	] as const;
}

export interface ProfileTabsProps {
	username: string;
}

export default function ProfileTabs({ username }: ProfileTabsProps) {
	const tabs = useRef(buildProfileTabs(username)).current;
	const pathname = usePathname();

	return (
		<section className={styles.tabs}>
			<div className={styles.tabsContent}>
				{tabs.map((tab) => (
					<Link
						data-active={tab.href === pathname ? '' : undefined}
						key={tab.name}
						className={styles.tabsItem}
						href={tab.href}
					>
						{tab.name}
					</Link>
				))}
			</div>
		</section>
	);
}
