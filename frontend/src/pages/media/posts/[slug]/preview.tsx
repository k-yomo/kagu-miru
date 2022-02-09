import React, { useEffect, useState } from 'react';
import Head from 'next/head';
import Error from 'next/error';
import { useRouter } from 'next/router';
import groq from 'groq';
import { sanityPreviewClient } from '@src/lib/sanityClient';
import PostDetail, {
  postFragmentForPostDetail,
  PostFragment,
} from '@src/components/PostDetail';

const fetchPostQuery = groq`*[_type == "post" && slug.current == $slug][0]{
  ${postFragmentForPostDetail}
}`;

const PostPreview = () => {
  const [data, setData] = useState<PostFragment | undefined>(undefined);
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  useEffect(() => {
    if (router.query.slug) {
      sanityPreviewClient
        .fetch(fetchPostQuery, { slug: router.query.slug })
        .then((data) => {
          setData(data);
        })
        .finally(() => {
          setLoading(false);
        });
    }
  }, [router.query]);

  if (loading) {
    return <></>;
  }

  if (!data) {
    return <Error statusCode={404} />;
  }

  return (
    <>
      <Head>
        <meta name="robots" content="noindex,nofollow,noarchive" />
      </Head>
      <PostDetail {...data} />
    </>
  );
};

PostPreview.theme = 'light';
