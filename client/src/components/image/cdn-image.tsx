'use client';

import { useState } from 'react';
import Image, { ImageProps } from 'next/image';
import cc from 'classcat';
import styles from './image.module.scss';

const HOSTNAME = process.env.NEXT_PUBLIC_CDN_HOSTNAME;
const PATHNAME = process.env.NEXT_PUBLIC_CDN_PATHNAME;

export const buildCdnUrl = (pathname: string) => {
	return `https://${HOSTNAME}/${PATHNAME}/${pathname}`;
};

export interface CDNImageProps extends Omit<ImageProps, 'src'> {
	path: string;
}

export function CDNImage({ className, path, quality = 100, ...props }: CDNImageProps) {
	const [isLoading, setIsLoading] = useState(true);
	const src = buildCdnUrl(path);

	return (
		<Image
			data-loading={isLoading ? '' : undefined}
			className={cc([styles.cdnImage, className])}
			src={src}
			placeholder='blur'
			blurDataURL={`/_next/image?url=${encodeURIComponent(src)}&q=1&w=64`}
			onLoadingComplete={() => setIsLoading(false)}
			quality={quality}
			{...props}
		/>
	);
}
