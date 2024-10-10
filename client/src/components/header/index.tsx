import { ArrowUpTrayIcon } from '@heroicons/react/16/solid';
import { AttentionButton, Button } from '../ui/button';
import UserDropdown from '../user/user-dropdown';
import { BellIcon } from '@heroicons/react/24/outline';
import { CursorArrowRaysIcon } from '@heroicons/react/20/solid';
import styles from './header.module.scss';
import HeaderNav from './header-nav';

export default function Header() {
	return (
		<header className={styles.header}>
			<div className={styles.headerContent}>
				<div className={styles.headerSection}>
					<span className={styles.headerLogo}>GoPix</span>
					<AttentionButton style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
						Upload <ArrowUpTrayIcon style={{ height: '16px', width: '16px' }} />
					</AttentionButton>
				</div>

				<HeaderNav />

				<div className={styles.headerSection}>
					<div className={styles.headerActions}>
						<Button size='icon' variant='ghost'>
							<CursorArrowRaysIcon />
						</Button>
						<Button size='icon' variant='ghost'>
							<BellIcon />
						</Button>
					</div>

					<UserDropdown />
				</div>
			</div>
		</header>
	);
}
