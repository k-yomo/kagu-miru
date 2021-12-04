import React, { memo } from 'react';
import Image from 'next/image';
import Link from 'next/link';
import { routes } from '@src/routes/routes';
import PostCategoryBadge from '@src/components/PostCategoryBadge';
import { truncate } from '@src/lib/string';
import { buildSanityImageSrc } from '@src/lib/sanityClient';
import { formatDistance, parseISO } from 'date-fns';
import { SanityImageSource } from '@sanity/image-url/lib/types/types';

export interface PostMeta {
  slug: string;
  title: string;
  description: string;
  mainImage: SanityImageSource;
  publishedAt?: string;
  categories?: Array<{ id: string; names: string[] }>;
}

interface Props {
  postMeta: PostMeta;
}

export default memo(function PostCard({ postMeta }: Props) {
  const mainImageUrl = buildSanityImageSrc(postMeta.mainImage);
  const blurImgUrl = buildSanityImageSrc(postMeta.mainImage).blur(10).url();
  return (
    <div className="mx-auto max-w-[400px] h-full shadow-md rounded-md">
      <Link href={routes.mediaPost(postMeta.slug)}>
        <a>
          <div className="relative w-full h-[200px] rounded-t-md">
            <Image
              src={mainImageUrl.url()}
              alt={postMeta.title}
              blurDataURL={blurImgUrl}
              placeholder="blur"
              layout="fill"
              objectFit="cover"
              objectPosition="center"
              loading="lazy"
              lazyBoundary="600px"
              unoptimized
            />
          </div>
          <div className="p-2">
            {postMeta.categories && (
              <div className="mb-2">
                {postMeta.categories.map((category) => (
                  <PostCategoryBadge
                    key={category.id}
                    id={category.id}
                    name={category.names[category.names.length - 1]}
                  />
                ))}
              </div>
            )}
            <h3 className="font-bold mb-1">{postMeta.title}</h3>
            <span className="text-sm text-text-secondary dark:text-text-secondary-dark">
              {truncate(postMeta.description, 50)}
            </span>
            <div className="my-2 text-sm text-text-secondary dark:text-text-secondary-dark">
              {postMeta.publishedAt &&
                formatDistance(parseISO(postMeta.publishedAt), new Date(), {
                  addSuffix: true,
                })}
            </div>
          </div>
        </a>
      </Link>
    </div>
  );
});
