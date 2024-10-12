'use client';

import { forwardRef } from 'react';
import cc from 'classcat';
import styles from './skeleton.module.scss';

interface SkeletonProps extends React.ComponentPropsWithoutRef<'div'> {
	loaded?: boolean;
}

export const Skeleton = forwardRef<HTMLDivElement, SkeletonProps>(
	({ className, loaded = false, ...props }, ref) => (
		<div ref={ref} className={cc([styles.skeleton, className])} data-loaded={loaded} {...props} />
	)
);
Skeleton.displayName = 'Skeleton';
