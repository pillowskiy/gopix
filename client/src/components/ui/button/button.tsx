'use client';

import { forwardRef } from 'react';
import {
	Button as HeadlessButton,
	type ButtonProps as HeadlessButtonProps
} from '@headlessui/react';
import cc from 'classcat';
import styles from './button.module.scss';

const buttonVariantsStyles = {
	accent: styles.btnVariantAccent,
	ghost: styles.btnVariantGhost,
	secondary: styles.btnVariantSecondary
} as const;

const buttonSizesStyles = {
	micro: styles.btnSizeMicro,
	small: styles.btnSizeSmall,
	medium: styles.btnSizeMedium,
	large: styles.btnSizeLarge,
	icon: styles.btnSizeIcon
} as const;

export interface ButtonProps extends HeadlessButtonProps {
	variant?: keyof typeof buttonVariantsStyles;
	size?: keyof typeof buttonSizesStyles;
}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
	({ children, variant = 'accent', size = 'medium', className, ...props }, ref) => {
		return (
			<HeadlessButton
				ref={ref}
				className={cc([
					styles.btn,
					buttonVariantsStyles[variant],
					buttonSizesStyles[size],
					className
				])}
				{...props}
			>
				{children}
			</HeadlessButton>
		);
	}
);
