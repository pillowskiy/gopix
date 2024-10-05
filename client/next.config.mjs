import bundleAnalyzer from '@next/bundle-analyzer';

/** @type {import('next').NextConfig} */
const nextConfig = {
	reactStrictMode: true,
	output: 'standalone',
	swcMinify: true,
	experimental: {
		optimizePackageImports: ['package-name']
	}
};

const withBundleAnalyzer = bundleAnalyzer({
	enabled: process.env.ANALYZE === 'true'
});

export default withBundleAnalyzer(nextConfig);
