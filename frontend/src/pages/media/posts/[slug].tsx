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
import { useItemDetailPageGetItemQuery } from '@src/generated/graphql';
import Link from 'next/link';
import { truncate } from '@src/lib/string';
import Rating from '@src/components/Rating';
import PlatformBadge from '@src/components/PlatformBadge';

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
            unoptimized
          />
        </div>
      );
    },
    item: ({ node }: { node: { id: string } }) => {
      // eslint-disable-next-line react-hooks/rules-of-hooks
      const { data, loading, error } = useItemDetailPageGetItemQuery({
        variables: { id: node.id },
      });
      const item = data?.getItem;
      if (!item || loading) {
        return <div>loading...</div>;
      }
      if (error) {
        return (
          <div className="p-3 shadow-md rounded-md">
            商品の取得に失敗しました。
          </div>
        );
      }
      return (
        <a href={!!item.affiliateUrl ? item.affiliateUrl : item.url}>
          <div className="flex items-center shadow-md rounded-md">
            <div className="w-[40%] h-full overflow-hidden">
              <img
                src={item.imageUrls[0]}
                alt={item.name}
                className="w-full h-full max-h-[250px] object-cover object-center"
              />
            </div>
            <div className="ml-2 sm:ml-3 w-[70%]">
              <PlatformBadge platform={item.platform} size="sm" />
              <div className="my-1 line-clamp-3  font-bold text-sm sm:text-md">
                {item.name}
              </div>
              <div className="hidden sm:block text-text-secondary dark:text-text-secondary-dark">
                {truncate(item.description, 100)}
              </div>
              <div className="my-1 flex items-center">
                <Rating rating={item.averageRating} maxRating={5} />
                <div className="ml-1 text-sm text-gray-600 dark:text-gray-300">
                  {item.reviewCount}
                </div>
              </div>
              <div className="my-2 text-xl sm:text-2xl font-bold">
                {item.price}円
              </div>
            </div>
          </div>
        </a>
      );
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
