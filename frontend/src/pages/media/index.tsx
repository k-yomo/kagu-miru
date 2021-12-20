import React from 'react';
import { GetServerSideProps } from 'next';
import { useRouter } from 'next/router';
import groq from 'groq';
import { sanityClient } from '@src/lib/sanityClient';
import MediaTopImg from '@public/images/media_top.jpg';
import SEOMeta from '@src/components/SEOMeta';
import PostCard, { PostMeta } from '@src/components/PostCard';
import Image from 'next/image';

const fetchRecentlyPublishedPostsQuery = groq`*[_type == "post"]{
  "slug": slug.current,
  title,
  description,
  mainImage,
  publishedAt,
  categories,
} | order(publishedAt desc) [0..9]`;

export const getServerSideProps: GetServerSideProps = async (ctx) => {
  const props = await sanityClient.fetch(fetchRecentlyPublishedPostsQuery);
  ctx.res.setHeader('Cache-Control', 'public, max-age=3600');
  return { props: { posts: props } };
};

interface Props {
  posts: PostMeta[];
}

export default function MediaTopPage({ posts }: Props) {
  const router = useRouter();
  return (
    <>
      <SEOMeta
        title="カグミル - 家具・インテリア専門メディア"
        excludeSiteTitle
        description="おすすめのインテリアから家具の選び方・比較まで、日々の生活を彩る家具・インテリア情報を発信します。"
        img={{ srcPath: MediaTopImg.src }}
        path={router.asPath}
      />
      <div className="max-w-[1000px] sm:mx-auto mb-8">
        <div className="relative w-full h-[200px] sm:h-[300px] rounded-t-md">
          <Image
            src={MediaTopImg.src}
            alt="トップ画像"
            layout="fill"
            objectFit="cover"
            objectPosition="center"
            priority
            unoptimized
          />
        </div>
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
