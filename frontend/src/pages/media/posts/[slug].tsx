import React from 'react';
import { GetServerSideProps } from 'next';
import { useRouter } from 'next/router';
import groq from 'groq';
import { buildSanityImageSrc, sanityClient } from '@src/lib/sanityClient';
import SEOMeta from '@src/components/SEOMeta';
import PostDetail, {
  postFragmentForPostDetail,
  PostFragment,
} from '@src/components/PostDetail';

const fetchPostQuery = groq`*[_type == "post" && slug.current == $slug][0]{
  ${postFragmentForPostDetail}
}`;

export const getServerSideProps: GetServerSideProps = async (ctx) => {
  const { slug } = ctx.query;
  const props = await sanityClient.fetch(fetchPostQuery, { slug });
  if (Object.keys(props).length === 0) {
    return { notFound: true };
  }
  ctx.res.setHeader(
    'Cache-Control',
    'public, max-age=600, s-maxage=3600, stale-while-revalidate=864000'
  );
  return { props };
};
const Post = (props: PostFragment) => {
  const { title, description, mainImage } = props;
  const router = useRouter();
  const mainImgUrl = buildSanityImageSrc(mainImage).width(1000).url()!;

  return (
    <>
      <SEOMeta
        title={title}
        description={description}
        img={{ src: mainImgUrl }}
        path={router.asPath}
      />
      <PostDetail {...props} />
    </>
  );
};

Post.theme = 'light';
export default Post;
