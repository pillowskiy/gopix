'use client';

import {
	Button as HeadlessButton,
	Dialog as HeadlessDialog,
	DialogBackdrop as HeadlessDialogBackdrop,
	DialogPanel as HeadlessDialogPanel,
	DialogTitle as HeadlessDialogTitle,
	type DialogPanelProps,
	type DialogProps
} from '@headlessui/react';
import { XMarkIcon } from '@heroicons/react/16/solid';
import cc from 'classcat';
import { createContext, forwardRef, useCallback, useContext, useState } from 'react';
import { Button, type ButtonProps } from '../button';
import styles from './dialog.module.scss';

interface DialogContextState {
	open: () => void;
	close: () => void;
	isOpen: boolean;
}

const DialogContext = createContext<DialogContextState | null>(null);

const useDialog = (): DialogContextState => {
	const ctx = useContext(DialogContext);
	if (!ctx) {
		throw new Error('useDialog must be used within a DialogProvider');
	}

	return ctx;
};

export function Dialog({
	open: defaultOpen = false,
	children
}: React.PropsWithChildren<{ open?: boolean }>) {
	const [isOpen, setIsOpen] = useState(defaultOpen);

	const close = useCallback(() => setIsOpen(false), []);
	const open = useCallback(() => setIsOpen(true), []);

	return (
		<DialogContext.Provider value={{ open, close, isOpen }}>{children}</DialogContext.Provider>
	);
}

export const DialogOverlay = forwardRef<
	HTMLDivElement,
	React.PropsWithChildren<Omit<DialogProps, 'onClose'>>
>(({ className, children, ...props }, ref) => {
	const dialog = useDialog();

	return (
		<HeadlessDialog
			ref={ref}
			as='div'
			open={dialog.isOpen}
			onClose={dialog.close}
			className={cc([styles.dialog, className])}
			{...props}
		>
			<HeadlessDialogBackdrop transition className={styles.dialogBackdrop} />
			<div className={styles.dialogWrapper}>{children}</div>
		</HeadlessDialog>
	);
});
DialogOverlay.displayName = 'DialogOverlay';

export function DialogTrigger(props: Omit<React.PropsWithChildren<ButtonProps>, 'as'>) {
	const dialog = useDialog();

	return <Button onClick={dialog.open} {...props} />;
}

export const DialogPanel = forwardRef<HTMLDivElement, React.PropsWithChildren<DialogPanelProps>>(
	({ className, transition = true, children, ...props }, ref) => {
		const dialog = useDialog();
		return (
			<HeadlessDialogPanel
				ref={ref}
				className={cc([styles.dialogPanel, className])}
				transition={transition}
				{...props}
			>
				<HeadlessButton
					type='button'
					aria-label='Close dialog button'
					className={styles.dialogPanelX}
					onClick={dialog.close}
				>
					<XMarkIcon />
				</HeadlessButton>
				{children}
			</HeadlessDialogPanel>
		);
	}
);
DialogPanel.displayName = 'DialogPanel';

export const DialogTitle = HeadlessDialogTitle;
DialogTitle.displayName = 'DialogTitle';

export function DialogClose(props: React.PropsWithChildren<ButtonProps>) {
	const dialog = useDialog();
	return <Button onClick={dialog.close} {...props} />;
}
