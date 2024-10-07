'use client';

import { forwardRef } from 'react';
import {
	Button as HeadlessButton,
	type ButtonProps as HeadlessButtonProps
} from '@headlessui/react';
import cc from 'classcat';
import styles from './button.module.scss';

const buttonVariantsStyles = {
	accent: styles.btnAccent,
	ghost: styles.btnGhost,
	secondary: styles.btnSecondary
} as const;

export interface ButtonProps extends HeadlessButtonProps {
	variant?: keyof typeof buttonVariantsStyles;
}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
	({ children, variant = 'accent', className, ...props }, ref) => {
		return (
			<HeadlessButton
				ref={ref}
				className={cc([styles.btn, buttonVariantsStyles[variant], className])}
				{...props}
			>
				{children}
			</HeadlessButton>
		);
	}
);
