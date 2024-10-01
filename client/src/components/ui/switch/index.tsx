'use client';

import * as Headless from '@headlessui/react';
import cc from 'classcat';
import { forwardRef } from 'react';
import styles from './switch.module.scss';

const switchSizesStyles = {
	small: styles.switchSizeSmall,
	medium: styles.switchSizeMedium,
	large: styles.switchSizeLarge
} as const;

interface SwitchProps extends Headless.SwitchProps {
	size?: keyof typeof switchSizesStyles;
}

export const Switch = forwardRef<HTMLButtonElement, SwitchProps>(
	({ className, size = 'medium', ...props }, ref) => (
		<Headless.Switch
			ref={ref}
			className={cc([styles.switch, switchSizesStyles[size], className])}
			{...props}
		>
			<span className={styles.switchCircle}></span>
		</Headless.Switch>
	)
);
