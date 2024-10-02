import styles from './background-texture.module.scss';

export function BackgroundGradient() {
	return (
		<>
			<div className={styles.top} style={{ backgroundImage: 'url(/bg-top.jpg)' }} />
			<div className={styles.bottom} style={{ backgroundImage: 'url(/bg-bottom.jpg)' }} />
		</>
	);
}
