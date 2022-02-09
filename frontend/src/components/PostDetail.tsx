import {
  SanityImageObject,
  SanityImageSource,
} from '@sanity/image-url/lib/types/types';
import LinkWithThumbnail from '@src/components/LinkWithThumbnail';
import { routes } from '@src/routes/routes';
import { buildSanityImageSrc, sanityClient } from '@src/lib/sanityClient';
import { SearchFrom } from '@src/generated/graphql';
import SearchPageScreenImg from '@public/images/search_screen.jpeg';
import Image from 'next/image';
import ItemDetailCard from '@src/components/ItemDetailCard';
import { useRouter } from 'next/router';
import React, { useEffect } from 'react';
import PostCategoryBadge from '@src/components/PostCategoryBadge';
import PostTagBadge from '@src/components/PostTagBadge';
import AuthorIcon from '@src/components/AuthorIcon';
import { formatDistance, parseISO } from 'date-fns';
import jaLocale from 'date-fns/locale/ja';
import BlockContent from '@sanity/block-content-to-react';
import TableOfContents from '@src/components/TableOfContents';

export const postFragmentForPostDetail = `
  title,
  description,
  mainImage,
  publishedAt,
  categories,
  tags,
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
`;

export interface PostFragment {
  title: string;
  description: string;
  mainImage: SanityImageSource;
  authorName?: string;
  authorImage?: SanityImageSource;
  publishedAt?: string;
  categories?: Array<{ id: string; names: [] }>;
  tags?: Array<{ value: string }>;
  body: any[];
}

interface Box {
  title?: string;
  content: string;
}

interface InternalLink {
  slug: string;
  title: string;
  description: string;
  mainImage: SanityImageSource;
}

interface SearchPageLink {
  title: string;
  url: string;
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
    box: ({ node }: { node: Box }) => {
      return (
        <div className="my-4 px-2 py-4 bg-slate-100 rounded-md">
          {node.title && (
            <div className="pb-2 mb-2 border-b-[1px] border-gray-200 dark:border-gray-800 font-bold">
              {node.title}
            </div>
          )}
          <pre className="whitespace-pre-wrap">{node.content}</pre>
        </div>
      );
    },
    internalLink: ({ node }: { node: InternalLink }) => {
      return (
        <div className="my-4">
          <LinkWithThumbnail
            url={routes.mediaPost(node.slug)}
            title={node.title}
            subTitle={node.description}
            imgSrc={buildSanityImageSrc(node.mainImage).width(300).url()}
          />
        </div>
      );
    },
    searchPageLink: ({ node }: { node: SearchPageLink }) => {
      const url = new URL(node.url);
      url.searchParams.set('searchFrom', SearchFrom.Media);
      const urlAs = new URL(node.url);
      urlAs.searchParams.delete('searchFrom');
      return (
        <div className="my-4">
          <LinkWithThumbnail
            url={url.toString()}
            urlAs={urlAs.toString()}
            title={node.title}
            subTitle=""
            imgSrc={SearchPageScreenImg.src}
          />
        </div>
      );
    },
    image: ({ node }: { node: SanityImageObject }) => {
      const imgUrl = buildSanityImageSrc(node).quality(70).width(600).url();
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
            lazyBoundary="10000px"
            unoptimized
          />
        </div>
      );
    },
    item: ({ node }: { node: { id: string } }) => {
      return (
        <div className="my-4">
          <ItemDetailCard itemId={node.id} />
        </div>
      );
    },
    customHtml: ({ node }: { node: { html: string } }) => {
      return <div dangerouslySetInnerHTML={{ __html: node.html }} />;
    },
  },
};

interface Props extends PostFragment {}

export default function PostDetail({
  title,
  mainImage,
  authorName,
  authorImage,
  publishedAt,
  categories,
  tags,
  body,
}: Props) {
  const router = useRouter();
  // remove line breaks
  body = body.map((block) => {
    if (!['h2', 'h3', 'h4', 'normal'].includes(block.style) || !block.children)
      return block;
    block.children = block.children.map((child: any) => {
      child.text = child.text.replace(/\n/g, '').replace(/<br>/g, '');
      return child;
    });
    return block;
  });
  // remove empty lines
  body = body.filter(
    (block) =>
      !['h2', 'h3', 'h4', 'normal'].includes(block.style) ||
      (block.children && !!block.children[0].text)
  );
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
    <article
      id="post"
      className="max-w-[1000px] mx-auto mb-4 sm:my-8 leading-8"
    >
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
            <PostCategoryBadge
              key={category.id}
              id={category.id}
              name={category.names[category.names.length - 1]}
              enableLink
            />
          ))}
        </div>
        <div className="my-4 space-x-2 space-y-2">
          {tags?.map((tag) => (
            <PostTagBadge key={tag.value} name={tag.value} />
          ))}
        </div>

        <h1 className="my-4 text-3xl font-bold">{title}</h1>
        <div className="flex items-center my-4 text-sm text-gray-400">
          {authorName && <AuthorIcon name={authorName} imgSrc={authorImgUrl} />}
          <span className="ml-2">
            {publishedAt &&
              formatDistance(parseISO(publishedAt), new Date(), {
                addSuffix: true,
                locale: jaLocale,
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
  );
}
