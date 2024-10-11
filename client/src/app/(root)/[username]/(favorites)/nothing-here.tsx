'use client';

import { Button } from '@/components/ui/button';
import styles from './favorites-page.module.scss';
import { useUserStore } from '@/providers/auth/store';

export default function NothingHere({ username }: { username: string }) {
	const user = useUserStore();

	if (user.data?.username === username) {
		return (
			<section className={styles.emptySection}>
				<p className={styles.emptyText}>There is nothing here yet.</p>
				<Button>Follow some users and start saving</Button>
			</section>
		);
	}

	return (
		<section className={styles.emptySection}>
			<p className={styles.emptyText}>This user did not save anything yet.</p>
			<p className={styles.emptyText}>Come back later.</p>
		</section>
	);
}
