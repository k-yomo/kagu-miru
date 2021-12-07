import React from 'react';
import { GetServerSideProps } from 'next';
import { buildSanityImageSrc, sanityClient } from '@src/lib/sanityClient';
import SEOMeta from '@src/components/SEOMeta';
import PostDetail, {
  postFragmentForPostDetail,
  PostFragment,
} from '@src/components/PostDetail';
import groq from 'groq';

const fetchPostQuery = groq`*[_type == "post" && slug.current == $slug][0]{
  ${postFragmentForPostDetail}
}`;

export const getServerSideProps: GetServerSideProps = async (ctx) => {
  const { slug } = ctx.query;
  const props = await sanityClient.fetch(fetchPostQuery, { slug });
  if (Object.keys(props).length === 0) {
    return { notFound: true };
  }
  ctx.res.setHeader('Cache-Control', 'public, max-age=3600');
  return { props };
};
const Post = (props: PostFragment) => {
  const { title, description, mainImage } = props;
  const mainImgUrl = buildSanityImageSrc(mainImage).width(1000).url()!;

  return (
    <>
      <SEOMeta
        title={title}
        description={description}
        img={{ src: mainImgUrl }}
      />
      <PostDetail {...props} />
    </>
  );
};

export default Post;
