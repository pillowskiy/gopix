'use client';

import * as Headless from '@headlessui/react';
import { forwardRef, createContext, useState, useContext, useCallback } from 'react';
import { Button, type ButtonProps } from '../button';
import { XMarkIcon } from '@heroicons/react/16/solid';
import cc from 'classcat';
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
	React.PropsWithChildren<Omit<Headless.DialogProps, 'onClose'>>
>(({ className, children, ...props }, ref) => {
	const dialog = useDialog();

	return (
		<Headless.Dialog
			ref={ref}
			as='div'
			open={dialog.isOpen}
			onClose={dialog.close}
			className={cc([styles.dialog, className])}
			{...props}
		>
			<Headless.DialogBackdrop transition className={styles.dialogBackdrop} />
			<div className={styles.dialogWrapper}>{children}</div>
		</Headless.Dialog>
	);
});

export const DialogTrigger = forwardRef<HTMLButtonElement, ButtonProps>((props, ref) => {
	const dialog = useDialog();

	return <Button ref={ref} onClick={dialog.open} {...props} />;
});

export const DialogPanel = forwardRef<
	HTMLDivElement,
	React.PropsWithChildren<Headless.DialogPanelProps>
>(({ className, transition = true, children, ...props }, ref) => {
	const dialog = useDialog();
	return (
		<Headless.DialogPanel
			ref={ref}
			className={cc([styles.dialogPanel, className])}
			transition={transition}
			{...props}
		>
			<Headless.Button
				type='button'
				aria-label='Close dialog button'
				className={styles.dialogPanelX}
				onClick={dialog.close}
			>
				<XMarkIcon />
			</Headless.Button>
			{children}
		</Headless.DialogPanel>
	);
});

export const DialogTitle = Headless.DialogTitle;

export const DialogClose = forwardRef<HTMLButtonElement, ButtonProps>((props, ref) => {
	const dialog = useDialog();
	return <Button ref={ref} onClick={dialog.close} {...props} />;
});
