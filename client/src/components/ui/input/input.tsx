'use client';

import * as Headless from '@headlessui/react';
import { forwardRef } from 'react';
import cc from 'classcat';
import styles from './input.module.scss';

export const Input = forwardRef<HTMLInputElement, Headless.InputProps>(
	({ className, ...props }, ref) => (
		<Headless.Input ref={ref} className={cc([styles.input, className])} {...props} />
	)
);
