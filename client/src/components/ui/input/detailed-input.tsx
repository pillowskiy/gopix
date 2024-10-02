'use client';

import { forwardRef } from 'react';
import * as Headless from '@headlessui/react';
import { Input } from './input';
import cc from 'classcat';
import styles from './input.module.scss';

interface DetailedInputProps extends Headless.InputProps {
	title?: string;
	description?: string;
}

export const DetailedInput = forwardRef<HTMLInputElement, DetailedInputProps>(
	({ title, description, className, style, ...props }, ref) => (
		<div className={cc([styles.detailedContainer, className])} style={style}>
			<Headless.Field>
				{title && <Headless.Label className={styles.detailedLabel}>{title}</Headless.Label>}
				{description && (
					<Headless.Description className={styles.detailedDescription}>
						{description}
					</Headless.Description>
				)}
				<Input className={styles.detailedInput} ref={ref} {...props} />
			</Headless.Field>
		</div>
	)
);
