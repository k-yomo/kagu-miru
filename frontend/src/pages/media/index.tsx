import React, { memo } from 'react';
import { GetServerSideProps } from 'next';
import Link from 'next/link';
import groq from 'groq';
import { buildSanityImageSrc, sanityClient } from '@src/lib/sanityClient';
import { SanityImageSource } from '@sanity/image-url/lib/types/types';
import { routes } from '@src/routes/routes';
import CategoryTag from '@src/components/PostTagBadge';
import MediaTopImg from '@public/images/media_top.jpg';
import { truncate } from '@src/lib/string';
import { formatDistance, parseISO } from 'date-fns';
import SEOMeta from '@src/components/SEOMeta';

const fetchRecentlyPublishedPostsQuery = groq`*[_type == "post"][0..9]{
  "slug": slug.current,
  title,
  description,
  mainImage,
  publishedAt,
  "categories": categories[]->name,
} | order(publishedAt desc)`;

export const getServerSideProps: GetServerSideProps = async (ctx) => {
  const props = await sanityClient.fetch(fetchRecentlyPublishedPostsQuery);
  ctx.res.setHeader('Cache-Control', 'public, max-age=3600');
  return { props: { posts: props } };
};

interface PostMeta {
  slug: string;
  title: string;
  description: string;
  mainImage: SanityImageSource;
  publishedAt?: string;
  categories?: string[];
}

interface Props {
  posts: PostMeta[];
}

export default function MediaTopPage({ posts }: Props) {
  return (
    <>
      <SEOMeta
        title="カグミル - 家具・インテリア専門メディア"
        excludeSiteTitle
        description="おすすめのインテリアから家具の選び方・比較まで、日々の生活を彩る家具・インテリア情報を発信します。"
        img={{ srcPath: MediaTopImg.src }}
      />
      <div className="max-w-[1000px] sm:mx-auto mb-8">
        <img
          src={MediaTopImg.src}
          alt="トップ画像"
          className="w-full h-[200px] sm:h-[300px] object-cover object-center"
        />
        <div className="mx-4">
          <h2 className="text-3xl font-bold mt-8 mb-4">新着記事一覧</h2>
          <div className="grid sm:grid-cols-2 md:grid-cols-3 sm:gap-4 gap-y-4 sm:gap-y-8">
            {posts.map((post) => (
              <PostCard key={post.slug} postMeta={post} />
            ))}
          </div>
        </div>
      </div>
    </>
  );
}

const PostCard = memo(function PostCard({ postMeta }: { postMeta: PostMeta }) {
  const mainImageUrl = buildSanityImageSrc(postMeta.mainImage);
  return (
    <div className="mx-auto max-w-[400px] h-full shadow-md rounded-md">
      <Link href={routes.post(postMeta.slug)}>
        <a>
          <img
            src={mainImageUrl.url()}
            alt={postMeta.title}
            className="w-full h-[200px] object-cover object-center rounded-t-md"
          />
          <div className="p-2">
            {postMeta.categories && (
              <div className="mb-2">
                {postMeta.categories.map((category) => (
                  <CategoryTag key={category} name={category} />
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
