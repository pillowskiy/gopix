import bundleAnalyzer from '@next/bundle-analyzer';

console.log('HOSTNAME', process.env.NEXT_PUBLIC_CDN_HOSTNAME);
const HOSTNAME = process.env.NEXT_PUBLIC_CDN_HOSTNAME;
const PATHNAME = process.env.NEXT_PUBLIC_CDN_PATHNAME;

/** @type {import('next').NextConfig} */
const nextConfig = {
	reactStrictMode: true,
	output: 'standalone',
	swcMinify: true,
	experimental: {
		optimizePackageImports: ['package-name']
	},
	images: {
		remotePatterns: [
			{
				protocol: 'https',
				hostname: HOSTNAME,
				pathname: PATHNAME ? `/${PATHNAME}/*` : '*'
			}
		]
	}
};

const withBundleAnalyzer = bundleAnalyzer({
	enabled: process.env.ANALYZE === 'true'
});

export default withBundleAnalyzer(nextConfig);
