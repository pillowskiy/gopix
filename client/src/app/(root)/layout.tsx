import Header from '@/components/header';

export default function RootLayout({ children }: React.PropsWithChildren) {
	return (
		<>
			<Header />
			{children}
		</>
	);
}
