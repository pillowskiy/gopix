import styles from './header.module.scss';

export default function Header() {
	return (
		<header className={styles.header}>
			<div className={styles.headerContent}>
				<h1 className={styles.headerContentTitle}>Gopix</h1>
			</div>
		</header>
	);
}
