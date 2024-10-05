/** @type {import('next').NextConfig} */
const nextConfig = {
	reactStrictMode: true,
	standalone: true,
	bundlePagesRouterDependencies: true,
	optimizePackageImports: ['@headlessui/react'],
	swcMinify: true,
	webpack: (config, { isServer }) => {
		if (!isServer) {
			Object.assign(config.resolve.alias, {
				'react/jsx-runtime.js': 'preact/compat/jsx-runtime',
				react: 'preact/compat',
				'react-dom/test-utils': 'preact/test-utils',
				'react-dom': 'preact/compat'
			});
		}
		return config;
	}
};

export default nextConfig;
