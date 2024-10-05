'use client';

import { ButtonProps, Button as HeadlessButton } from '@headlessui/react';
import cc from 'classcat';
import { forwardRef } from 'react';
import styles from './button.module.scss';

export const AttentionButton = forwardRef<HTMLButtonElement, ButtonProps>(
	({ children, className, ...props }, ref) => {
		return (
			<div className={styles.attentionBtnWrapper}>
				<span className={styles.attentionBtnSpinner} />
				<HeadlessButton
					ref={ref}
					className={cc([styles.attentionBtn, styles.btn, className])}
					{...props}
				>
					{children}
				</HeadlessButton>
			</div>
		);
	}
);
