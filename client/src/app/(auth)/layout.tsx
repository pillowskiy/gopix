import { ArrowUturnLeftIcon } from '@heroicons/react/20/solid';
import Link from 'next/link';
import styles from './auth.module.scss';

export default function AuthLayout({ children }: React.PropsWithChildren) {
	return (
		<main className={styles.wrapper}>
			<div className={styles.container}>
				<Link className={styles.homeAnchor} aria-label='Go to home' href='/'>
					<ArrowUturnLeftIcon className={styles.homeAnchorIcon} />
				</Link>

				<span className={styles.logo}>G</span>
				<p className={styles.description}>Welcome back</p>

				{children}
			</div>
			<div className={styles.previewContainer}>
				<div className={styles.previewAuthor}>
					<div className={styles.previewAuthorAvatar}></div>
					<div>
						<div className={styles.previewAuthorName}>John Doe</div>
						<div className={styles.previewAuthorTitle}>Software Developer</div>
					</div>
				</div>
				<div className={styles.previewImage} style={{ backgroundImage: 'url(/photo.jpg)' }} />
			</div>
		</main>
	);
}
