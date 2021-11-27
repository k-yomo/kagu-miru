import React, { useEffect, useState } from 'react';
import Head from 'next/head';
import Error from 'next/error';
import { sanityPreviewClient } from '@src/lib/sanityClient';
import { useRouter } from 'next/router';
import Post, { Props } from '@src/pages/media/posts/[slug]';
import groq from 'groq';

// Copied from `@src/pages/media/posts/[slug]`
export const fetchPostQuery = groq`*[_type == "post" && slug.current == $slug][0]{
  title,
  description,
  mainImage,
  publishedAt,
  "tags": tags[]->value,
  "authorName": author->name,
  "authorImage": author->image,
  body[]{
    ...,
    _type == "internalLink" => {
      "slug": @->slug.current,
      "title": @->title,
      "description": @->description,
      "mainImage": @->mainImage,
    } 
  }
}`;

export default function PostPreview() {
  const [data, setData] = useState<Props | undefined>(undefined);
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
      <Post {...data} />
    </>
  );
}
