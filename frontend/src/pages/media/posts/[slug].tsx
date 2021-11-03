import React, { useEffect } from 'react';
import { GetServerSideProps } from 'next';
import { parseISO, format } from 'date-fns';
import groq from 'groq';
import BlockContent from '@sanity/block-content-to-react';
import { SanityImageSource } from '@sanity/image-url/lib/types/types';
import { ClockIcon } from '@heroicons/react/solid';
import { sanityClient, buildSanityImageSrc } from '@src/lib/sanityClient';
import SEOMeta from '@src/components/SEOMeta';
import CategoryTag from '@src/components/CategoryTag';
import TableOfContents from '@src/components/TableOfContents';
import { useRouter } from 'next/router';

const fetchPostQuery = groq`*[_type == "post" && slug.current == $slug][0]{
  title,
  description,
  mainImage,
  publishedAt,
  "categories": categories[]->name,
  "authorName": author->name,
  "authorImage": author->image,
  body
}`;

export const getServerSideProps: GetServerSideProps = async (ctx) => {
  const { slug } = ctx.query;
  const props = await sanityClient.fetch(fetchPostQuery, { slug });
  if (Object.keys(props).length === 0) {
    return { notFound: true };
  }
  return { props };
};

interface Props {
  title: string;
  description: string;
  mainImage: SanityImageSource;
  publishedAt: string;
  categories: string[];
  authorName: string;
  authorImage: SanityImageSource;
  body: any[];
}

const Post = ({
  title,
  description,
  mainImage,
  publishedAt,
  categories,
  body,
}: Props) => {
  const router = useRouter();
  const headings: string[] = body
    .filter((block) => block.style === 'h2')
    .map((block) => block.children[0].text);
  const firstH2 = body.findIndex((block) => block.style === 'h2');
  const bodyBeforeTOC = [...body.slice(0, firstH2)];
  const bodyAfterTOC = [...body.slice(firstH2)];

  const mainImgUrl = buildSanityImageSrc(mainImage).url() || '';

  useEffect(() => {
    const hash = router.asPath.split('#')[1] ?? '';
    if (hash) {
      const section = document.getElementsByTagName('h2');
      section[parseInt(hash)]?.scrollIntoView(true);
    }
  }, [router.asPath]);

  return (
    <>
      <SEOMeta
        title={title}
        description={description}
        img={{ src: mainImgUrl }}
      />
      <article id="post" className="max-w-[1000px] mx-auto sm:my-8">
        <img src={mainImgUrl} alt={title} className="rounded-md" />
        <div className="mx-3">
          <div className="my-4 space-x-2">
            {categories.map((category) => (
              <CategoryTag key={category} name={category} />
            ))}
          </div>

          <h1 className="my-4 text-3xl font-bold">{title}</h1>
          <div className="flex items-center justify-end text-gray-400">
            <ClockIcon className="w-5 h-5 mr-1" />{' '}
            {format(parseISO(publishedAt), 'yyyy/M/d')}
          </div>
          <BlockContent
            blocks={bodyBeforeTOC}
            imageOptions={{ w: 320, h: 240, fit: 'max' }}
            {...sanityClient.config()}
          />
          <div className="my-10">
            <TableOfContents headings={headings} />
          </div>
          <BlockContent
            blocks={bodyAfterTOC}
            imageOptions={{ w: 320, h: 240, fit: 'max' }}
            {...sanityClient.config()}
          />
        </div>
      </article>
    </>
  );
};

export default Post;
