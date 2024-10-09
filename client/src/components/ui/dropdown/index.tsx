'use client';

import {
	Menu as HeadlessMenu,
	MenuButton as HeadlessMenuButton,
	MenuItems as HeadlessMenuItems,
	MenuItem as HeadlessMenuItem,
	type MenuButtonProps,
	type MenuItemsProps,
	type ButtonProps as HeadlessButtonProps
} from '@headlessui/react';
import cc from 'classcat';
import styles from './dropdown.module.scss';
import { forwardRef } from 'react';
import { Button, ButtonProps } from '@/components/ui/button';
import { ChevronDownIcon } from '@heroicons/react/16/solid';

export const DropdownMenu = HeadlessMenu;

export const DropdownMenuTrigger = forwardRef<HTMLButtonElement, MenuButtonProps>(
	({ children, className, as, ...props }, ref) => (
		<HeadlessMenuButton
			ref={ref}
			className={cc([styles.dropdownTrigger, className])}
			as={as ?? DropdownDefaultTrigger}
			{...props}
		>
			{children}
		</HeadlessMenuButton>
	)
);

const DropdownDefaultTrigger = forwardRef<
	HTMLButtonElement,
	Omit<React.PropsWithChildren<ButtonProps>, 'as'>
>(({ children, className, ...props }, ref) => (
	<Button ref={ref} className={cc([styles.dropdownTrigger, className])} {...props}>
		{children}
		<ChevronDownIcon className={styles.dropdownTriggerIcon} />
	</Button>
));

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

export const DropdownMenuItem = forwardRef<
	HTMLElement,
	React.PropsWithChildren<HeadlessButtonProps>
>(({ children, className, as: Component = 'button', ...props }, ref) => (
	<HeadlessMenuItem ref={ref} {...props}>
		<Component className={cc([styles.dropdownItem, className])}>{children}</Component>
	</HeadlessMenuItem>
));

export const DropdownMenuSeparator = forwardRef<
	HTMLDivElement,
	Omit<React.ComponentProps<'div'>, 'children'>
>(({ className, ...props }) => (
	<div className={cc([styles.dropdownSeparator, className])} {...props} />
));
