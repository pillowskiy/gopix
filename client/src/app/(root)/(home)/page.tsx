import Section from '@/components/section';
import styles from './page.module.scss';
import { BackgroundGradient } from '@/components/backdroung-texture';

export default async function Home() {
	return (
		<>
			<BackgroundGradient />
			<Section className={styles.page} container={false}>
				Hello World
				<section className={styles.container}></section>
			</Section>
		</>
	);
}
