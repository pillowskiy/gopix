'use client';

import { useState } from 'react';
import styles from './present-card.module.scss';
import cc from 'classcat';
import { Switch } from '../ui/switch';

interface PresentCardProps extends React.ComponentProps<'div'> {
	title: string;
}

export function PresentCard({ title, children, className, ...props }: PresentCardProps) {
	const [withHighlight, setWithHighlight] = useState(true);

	return (
		<div className={cc([styles.container, className])} {...props}>
			<div className={styles.card}>
				<div className={styles.cardActions}>
					Highlight: <Switch checked={withHighlight} onChange={setWithHighlight} size='small' />
				</div>
				{withHighlight && <div className={styles.cardHighlight} />}
				{children}
			</div>
			<h5 className={styles.title}>{title}</h5>
		</div>
	);
}
