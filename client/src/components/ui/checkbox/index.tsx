import {
	Checkbox as HeadlessCheckbox,
	type CheckboxProps as HeadlessCheckboxProps
} from '@headlessui/react';
import cc from 'classcat';
import { forwardRef } from 'react';
import styles from './checkbox.module.scss';

const checkboxSizesStyles = {
	small: styles.checkboxSizeSmall,
	medium: styles.checkboxSizeMedium,
	large: styles.checkboxSizeLarge
} as const;

interface CheckboxProps extends HeadlessCheckboxProps {
	size?: keyof typeof checkboxSizesStyles;
}

export const Checkbox = forwardRef<HTMLSpanElement, CheckboxProps>(
	({ className, size = 'medium', ...props }, ref) => (
		<HeadlessCheckbox
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
		</HeadlessCheckbox>
	)
);
