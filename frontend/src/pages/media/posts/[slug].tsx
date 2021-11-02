import React from 'react';
import { GetServerSideProps } from 'next';
import groq from 'groq';
import BlockContent from '@sanity/block-content-to-react';
import { SanityImageSource } from '@sanity/image-url/lib/types/types';
import { sanityClient, buildSanityImageSrc } from '@src/lib/sanityClient';
import SEOMeta from '@src/components/SEOMeta';

const fetchPostQuery = groq`*[_type == "post" && slug.current == $slug][0]{
  title,
  description,
  mainImage,
  "categories": categories[]->name,
  "authorName": author->name,
  "authorImage": author->image,
  body
}`;

export const getServerSideProps: GetServerSideProps = async (ctx) => {
  const { slug = '' } = ctx.query;
  return { props: await sanityClient.fetch(fetchPostQuery, { slug }) };
};

interface Props {
  title: string;
  description: string;
  mainImage: SanityImageSource;
  categories: string[];
  authorName: string;
  authorImage: SanityImageSource;
  body: any[] | any;
}

const Post = ({ title, description, mainImage, body }: Props) => {
  const mainImgUrl = buildSanityImageSrc(mainImage).url() || '';
  return (
    <>
      <SEOMeta
        title={title}
        description={description}
        img={{ src: mainImgUrl }}
      />
      <article className="max-w-[1000px] mx-auto my-8">
        <img src={mainImgUrl} alt={title} className="rounded-md" />
        <h1 className="my-4 text-4xl font-bold">{title}</h1>
        <BlockContent
          blocks={body}
          imageOptions={{ w: 320, h: 240, fit: 'max' }}
          {...sanityClient.config()}
        />
      </article>
    </>
  );
};

export default Post;
