import Section from '@/components/section';
import styles from './page.module.scss';
import { LightGradient } from '@/components/backdroung-texture';

export default function Home() {
	return (
		<>
			<LightGradient />
			<Section className={styles.page} container={false}>
				Hello World
			</Section>
		</>
	);
}
