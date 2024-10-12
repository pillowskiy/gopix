'use client';

import { Input as HeadlessInput, type InputProps as HeadlessInputProps } from '@headlessui/react';
import { forwardRef } from 'react';
import cc from 'classcat';
import styles from './input.module.scss';

export const Input = forwardRef<HTMLInputElement, HeadlessInputProps>(
	({ className, ...props }, ref) => (
		<HeadlessInput ref={ref} className={cc([styles.input, className])} {...props} />
	)
);
Input.displayName = 'Input';
