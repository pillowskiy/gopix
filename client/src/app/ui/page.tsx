import { PresentCard } from '@/components/present-card';
import { Button } from '@/components/ui/button';
import styles from './typo.module.scss';
import { Switch } from '@/components/ui/switch';
import BackgroundTexture from '@/components/backdroung-texture';

export default function TypoPage() {
	return (
		<>
			<BackgroundTexture />
			<div className={styles.header}>
				<h1 className={styles.headerTitle}>
					Gopix UI brings headless UI components to life, inspired by HeadlessUI, seamlessly
					blending innovation with flexibility.
				</h1>
			</div>
			<div className={styles.wrapper} style={{}}>
				<div className={styles.container}>
					<PresentCard title='Button Accent (default)'>
						<Button>Save changes</Button>
					</PresentCard>

					<PresentCard title='Button Ghost'>
						<Button variant='ghost'>Save changes</Button>
					</PresentCard>

					<PresentCard title='Button Secondary'>
						<Button variant='secondary'>Save changes</Button>
					</PresentCard>

					<PresentCard title='Switch'>
						<Switch size='medium' />
					</PresentCard>
				</div>
			</div>
		</>
	);
}
