'use client';

import * as Headless from '@headlessui/react';
import cc from 'classcat';
import styles from './dropdown.module.scss';
import { forwardRef } from 'react';
import { Button, ButtonProps } from '@/components/ui/button';
import { ChevronDownIcon } from '@heroicons/react/16/solid';

export const DropdownMenu = Headless.Menu;

export const DropdownMenuTrigger = forwardRef<HTMLButtonElement, Headless.MenuButtonProps>(
	({ children, className, as = DropdownDefaultTrigger, ...props }, ref) => (
		<Headless.MenuButton
			ref={ref}
			className={cc([styles.dropdownTrigger, className])}
			as={as}
			{...props}
		>
			{children}
		</Headless.MenuButton>
	)
);

const DropdownDefaultTrigger = forwardRef<HTMLButtonElement, Omit<ButtonProps, 'as'>>(
	({ children, className, ...props }, ref) => (
		<Button ref={ref} className={cc([styles.dropdownTrigger, className])} {...props}>
			{children}
			<ChevronDownIcon className={styles.dropdownTriggerIcon} />
		</Button>
	)
);

export const DropdownMenuItems = forwardRef<HTMLDivElement, Headless.MenuItemsProps>(
	({ className, transition = true, ...props }, ref) => (
		<Headless.MenuItems
			ref={ref}
			className={cc([styles.dropdownItems, className])}
			transition={transition}
			{...props}
		/>
	)
);

export const DropdownMenuItem = forwardRef<
	HTMLButtonElement,
	React.PropsWithChildren<Headless.ButtonProps>
>(({ children, className, as: Component = 'button', ...props }, ref) => (
	<Headless.MenuItem ref={ref} {...props}>
		<Component className={cc([styles.dropdownItem, className])}>{children}</Component>
	</Headless.MenuItem>
));

export const DropdownMenuSeparator = forwardRef<
	HTMLDivElement,
	Omit<React.ComponentProps<'div'>, 'children'>
>(({ className, ...props }) => (
	<div className={cc([styles.dropdownSeparator, className])} {...props} />
));
