import cc from 'classcat';

const sectionVariantStyles = {
	default: '',
	rounded: 'section__rounded'
} as const;

interface SectionProps extends React.ComponentProps<'div'> {
	container?: boolean;
	variant?: keyof typeof sectionVariantStyles;
}

export default function Section({
	children,
	variant = 'default',
	className,
	container = true,
	...props
}: SectionProps) {
	return (
		<div className={cc(['section__wrapper', sectionVariantStyles[variant]])} {...props}>
			{variant === 'default' && <div className='section__backdrop' />}
			<main className={cc([container && 'section__container', className])}>{children}</main>
		</div>
	);
}
