import { PresentCard } from '@/components/present-card';
import { Button, AttentionButton } from '@/components/ui/button';
import { Switch } from '@/components/ui/switch';
import { Input } from '@/components/ui/input/input';
import { DetailedInput } from '@/components/ui/input/detailed-input';
import { LightGradient } from '@/components/backdroung-texture';
import { Checkbox } from '@/components/ui/checkbox';
import { Link } from '@/components/ui/link';
import { Skeleton } from '@/components/ui/skeleton';
import {
	Dialog,
	DialogOverlay,
	DialogTrigger,
	DialogPanel,
	DialogTitle,
	DialogClose
} from '@/components/ui/dialog';
import {
	ArchiveBoxXMarkIcon,
	ArrowUpTrayIcon,
	PencilIcon,
	Square2StackIcon,
	TrashIcon
} from '@heroicons/react/16/solid';
import {
	DropdownMenu,
	DropdownMenuItem,
	DropdownMenuItems,
	DropdownMenuSeparator,
	DropdownMenuTrigger
} from '@/components/ui/dropdown';
import {
	Disclosure,
	DisclosureItem,
	DisclosureTrigger,
	DisclosurePanel
} from '@/components/ui/disclosure';
import Section from '@/components/section';
import styles from './typo.module.scss';

export default function TypoPage() {
	return (
		<>
			<div className={styles.header}>
				<h1 className={styles.headerTitle}>
					Gopix UI brings headless UI components to life, inspired by HeadlessUI, seamlessly
					blending innovation with flexibility.
				</h1>
			</div>
			<LightGradient />
			<Section className={styles.container} variant='default'>
				<PresentCard title='Button Accent (default)'>
					<Button>Save changes</Button>
				</PresentCard>

				<PresentCard title='Button Ghost'>
					<Button variant='ghost'>Save changes</Button>{' '}
				</PresentCard>

				<PresentCard title='Button Secondary'>
					<Button variant='secondary'>Save changes</Button>
				</PresentCard>

				<PresentCard title='Attention Button'>
					<AttentionButton style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
						Upload <ArrowUpTrayIcon style={{ height: '16px', width: '16px' }} />
					</AttentionButton>
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
						<div style={{ width: '100%', display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
							<Skeleton style={{ width: '60%', height: '14px', borderRadius: '9999px' }} />
							<Skeleton style={{ width: '80%', height: '12px', borderRadius: '9999px' }} />
						</div>
					</div>
				</PresentCard>

				<PresentCard title='Dialog'>
					<Dialog>
						<DialogTrigger>Open Dialog</DialogTrigger>
						<DialogOverlay>
							<DialogPanel>
								<DialogTitle
									as='h3'
									style={{
										fontSize: '18px',
										lineHeight: '26px',
										fontWeight: '600'
									}}
								>
									Payment successful
								</DialogTitle>

								<p
									style={{
										marginTop: '10px',
										fontSize: '14px',
										lineHeight: '22px'
									}}
								>
									Your payment has been successfully submitted. We’ve sent you an email with all of
									the details of your order.
								</p>

								<div style={{ marginTop: '16px' }}>
									<DialogClose>Got it, thanks!</DialogClose>
								</div>
							</DialogPanel>
						</DialogOverlay>
					</Dialog>
				</PresentCard>

				<PresentCard title='Dropdown Menu'>
					<DropdownMenu>
						<DropdownMenuTrigger>Options</DropdownMenuTrigger>

						<DropdownMenuItems anchor='bottom' style={{ width: '206px' }}>
							<DropdownMenuItem>
								<PencilIcon className='icon' />
								Edit
							</DropdownMenuItem>

							<DropdownMenuItem>
								<Square2StackIcon className='icon' />
								Dublicate
							</DropdownMenuItem>

							<DropdownMenuSeparator />

							<DropdownMenuItem>
								<ArchiveBoxXMarkIcon className='icon' />
								Archive
							</DropdownMenuItem>

							<DropdownMenuItem>
								<TrashIcon className='icon' />
								Delete
							</DropdownMenuItem>
						</DropdownMenuItems>
					</DropdownMenu>
				</PresentCard>

				<PresentCard title='Disclosure'>
					<Disclosure>
						<DisclosureItem>
							<DisclosureTrigger>What is your refund policy?</DisclosureTrigger>
							<DisclosurePanel>
								If you're unhappy with your purchase, we'll refund you in full.
							</DisclosurePanel>
						</DisclosureItem>

						<DisclosureItem>
							<DisclosureTrigger>Do you offer technical support?</DisclosureTrigger>
							<DisclosurePanel>No.</DisclosurePanel>
						</DisclosureItem>
					</Disclosure>
				</PresentCard>
			</Section>
		</>
	);
}
