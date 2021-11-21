import React, { useEffect } from 'react';
import { GetServerSideProps } from 'next';
import { useRouter } from 'next/router';
import Image from 'next/image';
import { formatDistance, parseISO } from 'date-fns';
import groq from 'groq';
import BlockContent from '@sanity/block-content-to-react';
import {
  SanityImageObject,
  SanityImageSource,
} from '@sanity/image-url/lib/types/types';
import { buildSanityImageSrc, sanityClient } from '@src/lib/sanityClient';
import SEOMeta from '@src/components/SEOMeta';
import CategoryTag from '@src/components/CategoryTag';
import TableOfContents from '@src/components/TableOfContents';
import AuthorIcon from '@src/components/AuthorIcon';
import LinkWithThumbnail from '@src/components/LinkWithThumbnail';
import { routes } from '@src/routes/routes';
import ItemDetailCard from '@src/components/ItemDetailCard';

// Copy to `@src/pages/media/posts/preview/[slug]`
// TODO: Fix to use identical query
//  For now importing fetchPostQuery results in `undefined`
export const fetchPostQuery = groq`*[_type == "post" && slug.current == $slug][0]{
  title,
  description,
  mainImage,
  publishedAt,
  "categories": categories[]->name,
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

interface InternalLink {
  slug: string;
  title: string;
  description: string;
  mainImage: SanityImageSource;
}

const serializers = {
  marks: {
    link: ({ mark, children }: { mark: { href: string }; children: any }) => {
      return (
        <a href={mark.href} rel="noopener" className="underline">
          {children}
        </a>
      );
    },
  },
  types: {
    internalLink: ({ node }: { node: InternalLink }) => {
      return (
        <LinkWithThumbnail
          url={routes.post(node.slug)}
          title={node.title}
          subTitle={node.description}
          imgSrc={buildSanityImageSrc(node.mainImage).url()}
        />
      );
    },
    image: ({ node }: { node: SanityImageObject }) => {
      const imgUrl = buildSanityImageSrc(node).url();
      const blurImgUrl = buildSanityImageSrc(node).blur(10).url();
      return (
        <div className="relative w-full h-[250px] sm:h-[500px]">
          <Image
            src={imgUrl}
            blurDataURL={blurImgUrl}
            placeholder="blur"
            alt=""
            layout="fill"
            objectFit="cover"
            objectPosition="center"
            loading="lazy"
            lazyBoundary="600px"
            unoptimized
          />
        </div>
      );
    },
    item: ({ node }: { node: { id: string } }) => {
      return <ItemDetailCard itemId={node.id} />;
    },
  },
};

export const getServerSideProps: GetServerSideProps = async (ctx) => {
  const { slug } = ctx.query;
  const props = await sanityClient.fetch(fetchPostQuery, { slug });
  if (Object.keys(props).length === 0) {
    return { notFound: true };
  }
  ctx.res.setHeader('Cache-Control', 'public, max-age=86400');
  return { props };
};

export interface Props {
  title: string;
  description: string;
  mainImage: SanityImageSource;
  authorName?: string;
  authorImage?: SanityImageSource;
  publishedAt?: string;
  categories?: string[];
  body: any[];
}

const Post = ({
  title,
  description,
  mainImage,
  authorName,
  authorImage,
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

  const mainImgUrl = buildSanityImageSrc(mainImage).width(1000).url()!;
  const mainImgBlurUrl = buildSanityImageSrc(mainImage)
    .blur(10)
    .width(1000)
    .url()!;
  const authorImgUrl = authorImage
    ? buildSanityImageSrc(authorImage).width(100).url()!
    : '';

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
        <div className="relative w-full h-[300px] sm:h-[600px]">
          <Image
            src={mainImgUrl}
            blurDataURL={mainImgBlurUrl}
            alt={title}
            placeholder="blur"
            layout="fill"
            objectFit="cover"
            objectPosition="center"
            priority
            unoptimized
          />
        </div>
        <div className="mx-3">
          <div className="my-4 space-x-2">
            {categories?.map((category) => (
              <CategoryTag key={category} name={category} />
            ))}
          </div>

          <h1 className="my-4 text-3xl font-bold">{title}</h1>
          <div className="flex items-center my-4 text-sm text-gray-400">
            {authorName && (
              <AuthorIcon name={authorName} imgSrc={authorImgUrl} />
            )}
            <span className="ml-2">
              {publishedAt &&
                formatDistance(parseISO(publishedAt), new Date(), {
                  addSuffix: true,
                })}
            </span>
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
            serializers={serializers}
            imageOptions={{ w: 320, h: 240, fit: 'max' }}
            {...sanityClient.config()}
          />
        </div>
      </article>
    </>
  );
};

export default Post;
