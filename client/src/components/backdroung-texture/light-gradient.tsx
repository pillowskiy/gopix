import cc from 'classcat';
import styles from './background-texture.module.scss';

export function LightGradient() {
	return (
		<>
			<div
				className={cc([styles.lightRight, styles.dark])}
				style={{ backgroundImage: 'url(https://nextui.org/gradients/docs-right.png)' }}
			/>

			<div
				className={cc([styles.lightLeft, styles.dark])}
				style={{ backgroundImage: 'url(https://nextui.org/gradients/docs-right.png)' }}
			/>
		</>
	);
}
