'use client';

import { forwardRef } from 'react';
import {
	Field as HeadlessField,
	Label as HeadlessLabel,
	Description as HeadlessDescription,
	type InputProps as HeadlessInputProps
} from '@headlessui/react';
import { Input } from './input';
import cc from 'classcat';
import styles from './input.module.scss';

interface DetailedInputProps extends HeadlessInputProps {
	title?: string;
	description?: string;
}

export const DetailedInput = forwardRef<HTMLInputElement, DetailedInputProps>(
	({ title, description, className, style, ...props }, ref) => (
		<div className={cc([styles.detailedContainer, className])} style={style}>
			<HeadlessField>
				{title && <HeadlessLabel className={styles.detailedLabel}>{title}</HeadlessLabel>}
				{description && (
					<HeadlessDescription className={styles.detailedDescription}>
						{description}
					</HeadlessDescription>
				)}
				<Input className={styles.detailedInput} ref={ref} {...props} />
			</HeadlessField>
		</div>
	)
);
