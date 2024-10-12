import NextLink, { type LinkProps as NextLinkProps } from 'next/link';
import { forwardRef } from 'react';
import cc from 'classcat';
import styles from './link.module.scss';

const linkSizesStyles = {
	small: styles.linkSmall,
	base: styles.linkBase,
	large: styles.linkLarge,
	none: ''
} as const;

interface LinkProps extends NextLinkProps, Omit<React.ComponentPropsWithoutRef<'a'>, 'href'> {
	size?: keyof typeof linkSizesStyles;
	disabled?: boolean;
}

export const Link = forwardRef<HTMLAnchorElement, LinkProps>(
	({ size = 'none', disabled, className, children, ...props }, ref) => (
		<NextLink
			ref={ref}
			aria-disabled={disabled}
			data-disabled={disabled}
			className={cc([styles.link, linkSizesStyles[size], className])}
			{...props}
		>
			{children}
		</NextLink>
	)
);
Link.displayName = 'Link';
