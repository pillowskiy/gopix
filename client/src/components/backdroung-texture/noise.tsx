import styles from './background-texture.module.scss';

export function NoiseBackground() {
	return (
		<div
			className={styles.noise}
			style={{
				backgroundImage: 'url(/bg-noise.png)'
			}}
		/>
	);
}
