'use client';

import { usePathname } from 'next/navigation';
import { useEffect, useRef, useState } from 'react';
import cc from 'classcat';
import styles from './header.module.scss';
import Link from 'next/link';
import { useUserStore } from '@/providers/auth/store';

interface NavItem {
	name: string;
	href: `/${string}`;
}

const baseNavItems = [
	{
		name: 'Suggest',
		href: '/'
	},
	{
		name: 'Search',
		href: '/search'
	}
] satisfies NavItem[];

const buildNavItems = (username?: string): NavItem[] => {
	const items: NavItem[] = [...baseNavItems];

	const authItem: NavItem = {
		name: 'You',
		href: username ? `/${username}` : '/login'
	};
	items.splice(1, 0, authItem);

	return items;
};

export default function HeaderNav() {
	const [navItems, setNavItems] = useState(buildNavItems());
	const user = useUserStore();
	const path = usePathname();

	useEffect(() => {
		setNavItems(buildNavItems(user.data?.username));
	}, [user.data?.username]);

	console.log(path);

	return (
		<nav className={cc([styles.headerSection, styles.nav])}>
			{navItems.map((item) => (
				<Link
					data-active={path === item.href ? '' : void 0}
					className={styles.navItem}
					href={item.href}
					key={item.name}
				>
					{item.name}
				</Link>
			))}
		</nav>
	);
}
