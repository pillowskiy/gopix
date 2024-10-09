'use client';

import { useUserStore } from '@/providers/auth/store';
import { ArrowRightEndOnRectangleIcon } from '@heroicons/react/20/solid';
import {
	DropdownMenu,
	DropdownMenuItem,
	DropdownMenuItems,
	DropdownMenuSeparator,
	DropdownMenuTrigger
} from '../ui/dropdown';
import {
	ArrowRightStartOnRectangleIcon,
	Cog6ToothIcon,
	FolderIcon,
	HeartIcon,
	SunIcon,
	UserIcon
} from '@heroicons/react/16/solid';
import Link from 'next/link';
import Image from 'next/image';
import { Button, ButtonProps } from '@headlessui/react';
import { forwardRef } from 'react';
import { User } from '@/types/users';
import styles from './user.module.scss';

const UserDropdownTrigger = forwardRef<HTMLButtonElement, ButtonProps & { user: User }>(
	({ user, ...props }, ref) => {
		return (
			<Button ref={ref} aria-label='User Dropdown' {...props}>
				<span>{user.username}</span>
				<Image
					className={styles.dropdownTriggerImage}
					src='/photo.jpg'
					alt={`${user.username} avatar`}
					width={48}
					height={48}
				/>
			</Button>
		);
	}
);

export default function UserDropdown() {
	const { data: user } = useUserStore();

	if (!user) {
		return (
			<a aria-label='Login' href='/login'>
				<ArrowRightEndOnRectangleIcon />
			</a>
		);
	}

	return (
		<DropdownMenu>
			<DropdownMenuTrigger
				className={styles.dropdownTrigger}
				user={user}
				as={UserDropdownTrigger}
			/>
			<DropdownMenuItems anchor='bottom end' style={{ width: '206px' }}>
				<div style={{ padding: '16px' }}>
					<h5 style={{ fontSize: '16px' }}>{user.username}</h5>
					<p style={{ fontSize: '14px' }}>{user.email}</p>
				</div>

				<DropdownMenuSeparator />

				<Link href={`/${user.username}`} passHref legacyBehavior>
					<DropdownMenuItem as='a'>
						<UserIcon className='icon' />
						Your Profile
					</DropdownMenuItem>
				</Link>

				<Link href={`/${user.username}/albums`} passHref legacyBehavior>
					<DropdownMenuItem as='a'>
						<FolderIcon className='icon' />
						Your Albums
					</DropdownMenuItem>
				</Link>

				<Link href='/account/favorites' passHref legacyBehavior>
					<DropdownMenuItem as='a'>
						<HeartIcon className='icon' />
						Your Favorites
					</DropdownMenuItem>
				</Link>

				<Link href='/account/favorites' passHref legacyBehavior>
					<DropdownMenuItem as='a'>
						<Cog6ToothIcon className='icon' />
						Settings
					</DropdownMenuItem>
				</Link>

				<DropdownMenuSeparator />
				<DropdownMenuItem>
					<SunIcon className='icon' />
					Light Mode
				</DropdownMenuItem>
				<DropdownMenuItem>
					<ArrowRightStartOnRectangleIcon className='icon' />
					Logout
				</DropdownMenuItem>
			</DropdownMenuItems>
		</DropdownMenu>
	);
}
