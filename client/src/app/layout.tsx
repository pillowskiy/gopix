import { NoiseBackground } from '@/components/backdroung-texture';
import AuthClientWrapper from '@/providers/auth';
import '@/styles/globals.scss';
import type { Metadata } from 'next';
import { Inter } from 'next/font/google';

export const metadata: Metadata = {
	title: 'Gopix'
};

const inter = Inter({ subsets: ['latin'] });

export default function RootLayout({
	children
}: Readonly<{
	children: React.ReactNode;
}>) {
	return (
		<html lang='en'>
			<body className={inter.className}>
				<NoiseBackground />
				<AuthClientWrapper>{children}</AuthClientWrapper>
			</body>
		</html>
	);
}
