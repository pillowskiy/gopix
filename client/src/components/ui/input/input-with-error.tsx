'use client';

import { forwardRef } from 'react';
import { type InputProps as HeadlessInputProps } from '@headlessui/react';
import { Input } from './input';
import cc from 'classcat';
import styles from './input.module.scss';

interface InputWithErrorProps extends HeadlessInputProps {
	error?: string;
}

export const InputWithError = forwardRef<HTMLInputElement, InputWithErrorProps>(
	({ error, className, style, ...props }, ref) => (
		<div className={cc([styles.detailedContainer, className])} style={style}>
			<Input
				className={styles.detailedInput}
				data-error={error ? '' : undefined}
				ref={ref}
				{...props}
			/>
			{error && <p className={styles.detailedError}>{error}</p>}
		</div>
	)
);
