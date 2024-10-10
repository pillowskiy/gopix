import cc from 'classcat';

const sectionVariantStyles = {
	default: '',
	rounded: 'section__rounded'
} as const;

interface SectionProps extends React.ComponentProps<'div'> {
	container?: boolean;
	variant?: keyof typeof sectionVariantStyles;
}

export default function Section({ children, className, container = true, ...props }: SectionProps) {
	return (
		<div className='section__wrapper' {...props}>
			<main className={cc([container && 'section__container', className])}>{children}</main>
		</div>
	);
}
