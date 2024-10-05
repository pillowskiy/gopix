'use client';

import {
	Dialog as HeadlessDialog,
	DialogBackdrop as HeadlessDialogBackdrop,
	DialogPanel as HeadlessDialogPanel,
	DialogTitle as HeadlessDialogTitle,
	Button as HeadlessButton,
	type DialogProps,
	type DialogPanelProps
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

export const Dialog = ({
	open: defaultOpen = false,
	children
}: React.PropsWithChildren<{ open?: boolean }>) => {
	const [isOpen, setIsOpen] = useState(defaultOpen);

	const close = useCallback(() => setIsOpen(false), []);
	const open = useCallback(() => setIsOpen(true), []);

	return (
		<DialogContext.Provider value={{ open, close, isOpen }}>{children}</DialogContext.Provider>
	);
};

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

export const DialogTrigger = forwardRef<HTMLButtonElement, ButtonProps>((props, ref) => {
	const dialog = useDialog();

	return <Button ref={ref} onClick={dialog.open} {...props} />;
});

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

export const DialogTitle = HeadlessDialogTitle;

export const DialogClose = forwardRef<HTMLButtonElement, ButtonProps>((props, ref) => {
	const dialog = useDialog();
	return <Button ref={ref} onClick={dialog.close} {...props} />;
});
