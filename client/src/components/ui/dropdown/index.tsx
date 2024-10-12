'use client';

import { Button, ButtonProps } from '@/components/ui/button';
import {
	Menu as HeadlessMenu,
	MenuButton as HeadlessMenuButton,
	MenuItem as HeadlessMenuItem,
	MenuItems as HeadlessMenuItems,
	type MenuButtonProps,
	type MenuItemProps,
	type MenuItemsProps
} from '@headlessui/react';
import { ChevronDownIcon } from '@heroicons/react/16/solid';
import cc from 'classcat';
import { forwardRef } from 'react';
import styles from './dropdown.module.scss';

export const DropdownMenu = HeadlessMenu;

export function DropdownMenuTrigger<E extends React.ElementType = typeof DropdownDefaultTrigger>({
	children,
	className,
	as,
	...props
}: MenuButtonProps<E> & { className?: string }) {
	return (
		// @ts-expect-error
		<HeadlessMenuButton
			as={as ?? DropdownDefaultTrigger}
			className={cc([styles.dropdownTrigger, className])}
			{...props}
		>
			{children}
		</HeadlessMenuButton>
	);
}

function DropdownDefaultTrigger({
	children,
	className,
	...props
}: Omit<React.PropsWithChildren<ButtonProps>, 'as'>) {
	return (
		<Button className={cc([styles.dropdownTrigger, className])} {...props}>
			{children}
			<ChevronDownIcon className={styles.dropdownTriggerIcon} />
		</Button>
	);
}

export const DropdownMenuItems = forwardRef<HTMLDivElement, MenuItemsProps>(
	({ className, transition = true, ...props }, ref) => (
		<HeadlessMenuItems
			ref={ref}
			className={cc([styles.dropdownItems, className])}
			transition={transition}
			{...props}
		/>
	)
);
DropdownMenuItems.displayName = 'DropdownMenuItems';

export function DropdownMenuItem<E extends React.ElementType = 'button'>({
	children,
	className,
	as,
	...props
}: MenuItemProps<E> & { className?: string }) {
	return (
		// @ts-expect-error
		<HeadlessMenuItem
			as={as ?? 'button'}
			className={cc([styles.dropdownItem, className])}
			{...props}
		>
			{children}
		</HeadlessMenuItem>
	);
}

export const DropdownMenuSeparator = forwardRef<
	HTMLDivElement,
	Omit<React.ComponentProps<'div'>, 'children'>
>(({ className, ...props }, ref) => (
	<div ref={ref} className={cc([styles.dropdownSeparator, className])} {...props} />
));
DropdownMenuSeparator.displayName = 'DropdownMenuSeparator';
