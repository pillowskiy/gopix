import { PresentCard } from '@/components/present-card';
import { Button } from '@/components/ui/button';
import styles from './typo.module.scss';
import { Switch } from '@/components/ui/switch';
import { Input } from '@/components/ui/input/input';
import { DetailedInput } from '@/components/ui/input/detailed-input';
import { BackgroundGradient } from '@/components/backdroung-texture';
import { Checkbox } from '@/components/ui/checkbox';
import { Link } from '@/components/ui/link';
import { Skeleton } from '@/components/ui/skeleton';

export default function TypoPage() {
	return (
		<>
			<BackgroundGradient />
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

					<PresentCard title='Input'>
						<Input style={{ width: '100%' }} placeholder='Input' />
					</PresentCard>

					<PresentCard title='Detailed Input'>
						<DetailedInput
							style={{ width: '100%' }}
							className={styles.element}
							title='Title'
							description='Description'
						/>
					</PresentCard>

					<PresentCard title='Checkbox'>
						<Checkbox />
					</PresentCard>

					<PresentCard title='Link'>
						<Link href='#x'>Link</Link>
					</PresentCard>

					<PresentCard title='Skeleton'>
						<div
							style={{
								maxWidth: '300px',
								width: '100%',
								display: 'flex',
								alignItems: 'center',
								gap: '1rem'
							}}
						>
							<div>
								<Skeleton
									style={{
										width: '48px',
										height: '48px',
										borderRadius: '50%',
										overflow: 'hidden'
									}}
								/>
							</div>
							<div
								style={{ width: '100%', display: 'flex', flexDirection: 'column', gap: '0.5rem' }}
							>
								<Skeleton style={{ width: '60%', height: '14px', borderRadius: '9999px' }} />
								<Skeleton style={{ width: '80%', height: '12px', borderRadius: '9999px' }} />
							</div>
						</div>
					</PresentCard>
				</div>
			</div>
		</>
	);
}
