'use client';

import styles from './background-texture.module.scss';

export default function BackgroundTexture() {
	return (
		<>
			<div className={styles.top} style={{ backgroundImage: 'url(/bg-top.jpg)' }} />
			<div className={styles.bottom} style={{ backgroundImage: 'url(/bg-bottom.jpg)' }} />
			<div className={styles.noise} style={{ backgroundImage: 'url(/bg-noise.png)' }} />
		</>
	);
}
