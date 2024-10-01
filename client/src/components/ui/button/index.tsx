'use client';

import { Button as HeadlessButton } from '@headlessui/react';
import cc from 'classcat';
import { forwardRef } from 'react';
import styles from './button.module.scss';

const buttonVariantsStyles = {
	accent: styles.btnAccent,
	ghost: styles.btnGhost,
	secondary: styles.btnSecondary
} as const;

interface ButtonProps extends React.ComponentProps<'button'> {
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
