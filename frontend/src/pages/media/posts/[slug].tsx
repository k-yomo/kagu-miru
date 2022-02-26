import React from 'react';
import { GetStaticPaths, GetStaticProps } from 'next';
import { useRouter } from 'next/router';
import groq from 'groq';
import { buildSanityImageSrc, sanityClient } from '@src/lib/sanityClient';
import SEOMeta from '@src/components/SEOMeta';
import PostDetail, {
  postFragmentForPostDetail,
  PostFragment,
} from '@src/components/PostDetail';
import { routes } from '@src/routes/routes';

const fetchPostQuery = groq`*[_type == "post" && slug.current == $slug][0]{
  ${postFragmentForPostDetail}
}`;

export const getStaticPaths: GetStaticPaths = async () => {
  const query = groq`{
      "posts": *[_type == 'post']{slug},
    }`;
  const { posts } = await sanityClient.fetch(query);
  const paths = posts.map(({ slug }: { slug: { current: string } }) =>
    routes.mediaPost(slug.current)
  );
  return { paths, fallback: 'blocking' };
};

export const getStaticProps: GetStaticProps = async ({ params }) => {
  if (!params) {
    return { notFound: true };
  }
  const props = await sanityClient.fetch(fetchPostQuery, {
    slug: params.slug,
  });
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
