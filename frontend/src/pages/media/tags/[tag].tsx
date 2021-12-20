import React from 'react';
import { GetServerSideProps } from 'next';
import groq from 'groq';
import { sanityClient } from '@src/lib/sanityClient';
import MediaTopImg from '@public/images/media_top.jpg';
import SEOMeta from '@src/components/SEOMeta';
import PostCard, { PostMeta } from '@src/components/PostCard';
import { useRouter } from 'next/router';

export const getServerSideProps: GetServerSideProps = async (ctx) => {
  const tag = ctx.query.tag as string;
  const query = groq`*[_type == "post" && "${tag}" in tags[].value]{
  "slug": slug.current,
  title,
  description,
  mainImage,
  publishedAt,
  categories
} | order(publishedAt desc) [0..9]`;

  const props = await sanityClient.fetch(query);
  ctx.res.setHeader('Cache-Control', 'public, max-age=3600');
  return { props: { tag, posts: props } };
};

interface Props {
  tag: string;
  posts: PostMeta[];
}

export default function MediaTopPage({ tag, posts }: Props) {
  const router = useRouter();
  return (
    <>
      <SEOMeta
        title={`#${tag} 記事一覧`}
        description={`${tag}に関する記事の一覧ページです`}
        img={{ srcPath: MediaTopImg.src }}
        path={router.asPath}
      />
      <div className="max-w-[1000px] sm:mx-auto mb-8">
        <div className="mx-4">
          <h2 className="text-2xl font-bold mt-8 mb-4">#{tag}</h2>
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
