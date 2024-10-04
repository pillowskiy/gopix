'use client';

import cc from 'classcat';
import styles from './disclosure.module.scss';
import * as Headless from '@headlessui/react';
import { forwardRef } from 'react';
import { ChevronDownIcon } from '@heroicons/react/16/solid';

export const Disclosure = forwardRef<HTMLDivElement, React.ComponentPropsWithoutRef<'div'>>(
	({ className, children, ...props }, ref) => (
		<div className={cc([styles.disclosure, className])} ref={ref} {...props}>
			{children}
		</div>
	)
);

export const DisclosureItem = forwardRef<
	HTMLDivElement,
	React.PropsWithChildren<Headless.DisclosureProps & React.ComponentProps<'div'>>
>(({ className, children, ...props }, ref) => (
	<Headless.Disclosure
		as='div'
		className={cc([styles.disclosureItem, className])}
		ref={ref}
		{...props}
	>
		{children}
	</Headless.Disclosure>
));

export const DisclosureTrigger = forwardRef<
	HTMLButtonElement,
	React.PropsWithChildren<Headless.DisclosureButtonProps>
>(({ className, children, ...props }, ref) => (
	<Headless.DisclosureButton
		className={cc([styles.disclosureTrigger, className])}
		ref={ref}
		{...props}
	>
		<span className={styles.disclosureTriggerText}>{children}</span>
		<ChevronDownIcon className={styles.disclosureTriggerIcon} />
	</Headless.DisclosureButton>
));

export const DisclosurePanel = forwardRef<
	HTMLDivElement,
	React.PropsWithChildren<Headless.DisclosurePanelProps & React.ComponentProps<'div'>>
>(({ className, transition = true, children, ...props }, ref) => (
	<div className={styles.disclosurePanelWrapper}>
		<Headless.DisclosurePanel
			ref={ref}
			className={cc([styles.disclosurePanel, className])}
			transition={transition}
			{...props}
		>
			{children}
		</Headless.DisclosurePanel>
	</div>
));
