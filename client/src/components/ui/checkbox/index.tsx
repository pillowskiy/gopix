import * as Headless from '@headlessui/react';
import styles from './checkbox.module.scss';
import cc from 'classcat';
import { forwardRef } from 'react';

const checkboxSizesStyles = {
	small: styles.checkboxSizeSmall,
	medium: styles.checkboxSizeMedium,
	large: styles.checkboxSizeLarge
} as const;

interface CheckboxProps extends Headless.CheckboxProps {
	size?: keyof typeof checkboxSizesStyles;
}

export const Checkbox = forwardRef<HTMLSpanElement, CheckboxProps>(
	({ className, size = 'medium', ...props }, ref) => (
		<Headless.Checkbox
			ref={ref}
			className={cc([styles.checkbox, checkboxSizesStyles[size], className])}
			{...props}
		>
			{/* Checkmark icon */}
			<svg className={styles.checkboxCheckmark} viewBox='0 0 14 14' fill='none'>
				<title>Checkmark</title>
				<path
					d='M3 8L6 11L11 3.5'
					stroke='currentColor'
					strokeWidth={2}
					strokeLinecap='round'
					strokeLinejoin='round'
				/>
			</svg>
		</Headless.Checkbox>
	)
);
