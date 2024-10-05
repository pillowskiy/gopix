'use client';

import cc from 'classcat';
import styles from './disclosure.module.scss';
import {
	Disclosure as HeadlessDisclosure,
	DisclosureButton as HeadlessDisclosureButton,
	DisclosurePanel as HeadlessDisclosurePanel,
	type DisclosureProps,
	type DisclosureButtonProps,
	type DisclosurePanelProps
} from '@headlessui/react';
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
	React.PropsWithChildren<DisclosureProps & React.ComponentProps<'div'>>
>(({ className, children, ...props }, ref) => (
	<HeadlessDisclosure
		as='div'
		className={cc([styles.disclosureItem, className])}
		ref={ref}
		{...props}
	>
		{children}
	</HeadlessDisclosure>
));

export const DisclosureTrigger = forwardRef<
	HTMLButtonElement,
	React.PropsWithChildren<DisclosureButtonProps>
>(({ className, children, ...props }, ref) => (
	<HeadlessDisclosureButton
		className={cc([styles.disclosureTrigger, className])}
		ref={ref}
		{...props}
	>
		<span className={styles.disclosureTriggerText}>{children}</span>
		<ChevronDownIcon className={styles.disclosureTriggerIcon} />
	</HeadlessDisclosureButton>
));

export const DisclosurePanel = forwardRef<
	HTMLDivElement,
	React.PropsWithChildren<DisclosurePanelProps & React.ComponentProps<'div'>>
>(({ className, transition = true, children, ...props }, ref) => (
	<div className={styles.disclosurePanelWrapper}>
		<HeadlessDisclosurePanel
			ref={ref}
			className={cc([styles.disclosurePanel, className])}
			transition={transition}
			{...props}
		>
			{children}
		</HeadlessDisclosurePanel>
	</div>
));
