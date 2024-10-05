'use client';

import { useEffect, useRef, useState } from 'react';
import styles from './auth.module.scss';
import Link from 'next/link';

interface TabBackdropState {
	width: number;
	translateX: number;
}

const tabs = [
	{
		title: 'Log In',
		href: '/login'
	},
	{
		title: 'Sign Up',
		href: '/signup'
	}
] as const;

interface AuthTabsProps {
	defaultActive: keyof typeof tabs;
}

export default function AuthTabs({ defaultActive = 0 }: AuthTabsProps) {
	const targetRef = useRef<HTMLAnchorElement>(null);
	const [backdropSettings, setBackdropSettings] = useState<TabBackdropState | null>(null);

	useEffect(() => {
		if (targetRef.current) {
			setBackdropSettings({
				width: targetRef.current.offsetWidth,
				translateX: targetRef.current.offsetLeft
			});
		}
	}, []);

	const onTabClick = (e: React.MouseEvent<HTMLAnchorElement>) => {
		const target = e.currentTarget;
		if (!target) return;

		setBackdropSettings({
			width: target.offsetWidth,
			translateX: target.offsetLeft
		});
	};

	return (
		<div className={styles.loginTabs}>
			{backdropSettings && (
				<div
					className={styles.loginTabsBackdrop}
					style={{
						transform: `translateX(${backdropSettings.translateX}px)`,
						width: backdropSettings.width
					}}
				/>
			)}

			{tabs.map((tab, i) => (
				<Link
					key={tab.title}
					href={tab.href}
					onClick={onTabClick}
					ref={defaultActive === i ? targetRef : null}
					className={styles.loginTabsItem}
				>
					{tab.title}
				</Link>
			))}
		</div>
	);
}
