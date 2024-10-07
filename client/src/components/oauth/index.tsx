import { forwardRef } from 'react';
import Link, { LinkProps } from 'next/link';
import cc from 'classcat';
import styles from './oauth.module.scss';
import { Button } from '../ui/button';

type OAuthService = 'google' | 'apple';

const servicesIcons = {
	google: {
		name: 'Google',
		icon: (props: React.ComponentProps<'svg'>) => (
			<svg version='1.1' xmlns='http://www.w3.org/2000/svg' viewBox='0 0 48 48' {...props}>
				<title>Google Icon</title>
				<g>
					<path
						fill='#EA4335'
						d='M24 9.5c3.54 0 6.71 1.22 9.21 3.6l6.85-6.85C35.9 2.38 30.47 0 24 0 14.62 0 6.51 5.38 2.56 13.22l7.98 6.19C12.43 13.72 17.74 9.5 24 9.5z'
					></path>
					<path
						fill='#4285F4'
						d='M46.98 24.55c0-1.57-.15-3.09-.38-4.55H24v9.02h12.94c-.58 2.96-2.26 5.48-4.78 7.18l7.73 6c4.51-4.18 7.09-10.36 7.09-17.65z'
					></path>
					<path
						fill='#FBBC05'
						d='M10.53 28.59c-.48-1.45-.76-2.99-.76-4.59s.27-3.14.76-4.59l-7.98-6.19C.92 16.46 0 20.12 0 24c0 3.88.92 7.54 2.56 10.78l7.97-6.19z'
					></path>
					<path
						fill='#34A853'
						d='M24 48c6.48 0 11.93-2.13 15.89-5.81l-7.73-6c-2.15 1.45-4.92 2.3-8.16 2.3-6.26 0-11.57-4.22-13.47-9.91l-7.98 6.19C6.51 42.62 14.62 48 24 48z'
					></path>
					<path fill='none' d='M0 0h48v48H0z'></path>
				</g>
			</svg>
		)
	},
	apple: {
		name: 'Apple',
		icon: (props: React.ComponentProps<'svg'>) => (
			<svg
				width='256px'
				height='315px'
				viewBox='0 0 256 315'
				version='1.1'
				xmlns='http://www.w3.org/2000/svg'
				xmlnsXlink='http://www.w3.org/1999/xlink'
				preserveAspectRatio='xMidYMid'
				{...props}
			>
				<title>Apple Icon</title>
				<g>
					<path
						d='M213.803394,167.030943 C214.2452,214.609646 255.542482,230.442639 256,230.644727 C255.650812,231.761357 249.401383,253.208293 234.24263,275.361446 C221.138555,294.513969 207.538253,313.596333 186.113759,313.991545 C165.062051,314.379442 158.292752,301.507828 134.22469,301.507828 C110.163898,301.507828 102.642899,313.596301 82.7151126,314.379442 C62.0350407,315.16201 46.2873831,293.668525 33.0744079,274.586162 C6.07529317,235.552544 -14.5576169,164.286328 13.147166,116.18047 C26.9103111,92.2909053 51.5060917,77.1630356 78.2026125,76.7751096 C98.5099145,76.3877456 117.677594,90.4371851 130.091705,90.4371851 C142.497945,90.4371851 165.790755,73.5415029 190.277627,76.0228474 C200.528668,76.4495055 229.303509,80.1636878 247.780625,107.209389 C246.291825,108.132333 213.44635,127.253405 213.803394,167.030988 M174.239142,50.1987033 C185.218331,36.9088319 192.607958,18.4081019 190.591988,0 C174.766312,0.636050225 155.629514,10.5457909 144.278109,23.8283506 C134.10507,35.5906758 125.195775,54.4170275 127.599657,72.4607932 C145.239231,73.8255433 163.259413,63.4970262 174.239142,50.1987249'
						fill='#ffffff'
					></path>
				</g>
			</svg>
		)
	}
} as const satisfies Record<
	OAuthService,
	{ name: string; icon: React.FC<React.ComponentProps<'svg'>> }
>;

interface OAuthButtonProps extends LinkProps, Omit<React.ComponentProps<'a'>, 'href'> {
	service: OAuthService;
	prefix?: string;
}

export const OAuthButton = forwardRef<HTMLAnchorElement, OAuthButtonProps>(
	({ prefix = 'Continue with', href = '/#', service, className, ...props }, ref) => {
		const { icon: ServiceIcon, name } = servicesIcons[service];
		return (
			<Link ref={ref} href={href} passHref legacyBehavior {...props}>
				<Button as='a' variant='secondary' className={cc([styles.oauthBtn, className])}>
					<ServiceIcon className={styles.oauthBtnIcon} />
					{prefix} {name}
				</Button>
			</Link>
		);
	}
);
