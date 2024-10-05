'use client';

import {
	Switch as HeadlessSwitch,
	type SwitchProps as HeadlessSwitchProps
} from '@headlessui/react';
import cc from 'classcat';
import { forwardRef } from 'react';
import styles from './switch.module.scss';

const switchSizesStyles = {
	small: styles.switchSizeSmall,
	medium: styles.switchSizeMedium,
	large: styles.switchSizeLarge
} as const;

interface SwitchProps extends HeadlessSwitchProps {
	size?: keyof typeof switchSizesStyles;
}

export const Switch = forwardRef<HTMLButtonElement, SwitchProps>(
	({ className, size = 'medium', ...props }, ref) => (
		<HeadlessSwitch
			ref={ref}
			className={cc([styles.switch, switchSizesStyles[size], className])}
			{...props}
		>
			<span className={styles.switchCircle}></span>
		</HeadlessSwitch>
	)
);
