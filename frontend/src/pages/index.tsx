import React, { useState, useEffect, useCallback } from 'react';
import { GetServerSideProps } from 'next';
import { useRouter } from 'next/router';
import Link from 'next/link';
import Image from 'next/image';
import gql from 'graphql-tag';
import { ApolloError } from '@apollo/client';
import SEOMeta from '@src/components/SEOMeta';
import SearchPageScreenImg from '@public/images/search_screen.jpeg';
import HomeTopImg from '@public/images/home_top.jpg';
import apolloClient, { isErrorIncludes } from '@src/lib/apolloClient';
import {
  ErrorCode,
  HomeDocument,
  HomeQuery,
  HomeComponentPayloadItemGroups,
  HomeComponentPayloadCategories,
  HomeComponentPayloadMediaPosts,
  EventId,
  Action,
  HomeClickItemActionParams,
  useTrackEventMutation,
  SearchFrom,
} from '@src/generated/graphql';
import ItemList from '@src/components/ItemList';
import PostCard from '@src/components/PostCard';
import { routes } from '@src/routes/routes';

gql`
  query home {
    home {
      components {
        id
        payload {
          __typename
          ... on HomeComponentPayloadItemGroups {
            title
            payload {
              ... on HomeComponentPayloadItems {
                title
                items {
                  ...itemListItemFragment
                }
              }
            }
          }
          ... on HomeComponentPayloadItems {
            title
            items {
              ...itemListItemFragment
            }
          }
          ... on HomeComponentPayloadCategories {
            title
            categories {
              id
              name
              imageUrl
            }
          }
          ... on HomeComponentPayloadMediaPosts {
            title
            posts {
              slug
              title
              description
              mainImageUrl
              publishedAt
              categories {
                id
                names
              }
            }
          }
        }
      }
    }
  }
`;

export const getServerSideProps: GetServerSideProps = async (ctx) => {
  try {
    const { data } = await apolloClient.query<HomeQuery>({
      query: HomeDocument,
    });

    ctx.res.setHeader(
      'Cache-Control',
      'public, max-age=600, stale-while-revalidate=86400'
    );
    return { props: { components: data.home.components } };
  } catch (e) {
    if (e instanceof ApolloError && isErrorIncludes(e, ErrorCode.NotFound)) {
      return { notFound: true };
    }
    throw e;
  }
};

interface Props {
  components: HomeQuery['home']['components'];
}

export default function HomePage({ components }: Props) {
  const router = useRouter();
  const [showTopTitle, setShowTopTitle] = useState(false);

  useEffect(() => {
    setTimeout(() => {
      setShowTopTitle(true);
    }, 500);
  }, []);

  return (
    <>
      <SEOMeta
        title="カグミル - 家具検索サービス"
        excludeSiteTitle
        description="カグミルはオンラインで買える家具・インテリアを横断で一括検索・比較出来るサービスです。"
        img={{ srcPath: SearchPageScreenImg.src }}
        path={router.asPath}
      />
      <div className="relative w-full h-[200px] sm:h-[300px] rounded-t-md">
        <Image
          src={HomeTopImg.src}
          alt="トップ画像"
          layout="fill"
          objectFit="cover"
          objectPosition="center"
          priority
          unoptimized
        />
        <div className="flex justify-center items-center absolute bottom-0 w-full h-full p-2 bg-black bg-opacity-10">
          <h1
            className={`text-2xl sm:text-4xl font-bold text-white text-shadow ease-in-out duration-1000 ${
              showTopTitle ? 'opacity-100' : 'opacity-0 translate-y-2'
            } `}
          >
            カグミル
          </h1>
        </div>
      </div>
      <div className="mx-auto max-w-[1400px]">
        {components.map((c) => (
          <HomeComponent key={c.id} component={c} />
        ))}
      </div>
    </>
  );
}

function HomeComponent({
  component,
}: {
  component: HomeQuery['home']['components'][0];
}) {
  switch (component.payload.__typename) {
    case 'HomeComponentPayloadItemGroups': {
      return (
        <HomeComponentItemGroups
          id={component.id}
          payload={
            component.payload as unknown as HomeComponentPayloadItemGroups
          }
        />
      );
    }
    case 'HomeComponentPayloadItems':
      return <></>;
    case 'HomeComponentPayloadCategories': {
      return (
        <HomeComponentCategories
          payload={
            component.payload as unknown as HomeComponentPayloadCategories
          }
        />
      );
    }
    case 'HomeComponentPayloadMediaPosts':
      return (
        <HomeComponentPosts
          payload={
            component.payload as unknown as HomeComponentPayloadMediaPosts
          }
        />
      );
    default:
      const _: never = component.payload;
      return <></>;
  }
}

function HomeComponentItemGroups({
  id,
  payload,
}: {
  id: string;
  payload: HomeComponentPayloadItemGroups;
}) {
  const [trackEvent] = useTrackEventMutation();
  const onClickItem = useCallback(
    (itemId: string) => {
      const params: HomeClickItemActionParams = {
        componentId: id,
        itemId,
      };
      trackEvent({
        variables: {
          event: {
            id: EventId.Home,
            action: Action.ClickItem,
            createdAt: new Date(),
            params,
          },
        },
      }).catch(() => {
        // do nothing
      });
    },
    [id, trackEvent]
  );

  return (
    <div className="mx-3 my-4">
      <h2 className="pb-2 text-xl border-b-2 border-gray-200 dark:border-gray-800 font-bold">
        {payload.title}
      </h2>
      {payload.payload.map((itemsPayload) => (
        <div key={itemsPayload.title} className="py-2">
          <h2 className="mb-2 font-bold">{itemsPayload.title}</h2>
          <div className="grid grid-flow-col auto-cols-min space-x-4 overflow-x-auto">
            <ItemList
              items={itemsPayload.items}
              onClickItem={onClickItem}
              isAdmin={false}
            />
          </div>
        </div>
      ))}
    </div>
  );
}

function HomeComponentCategories({
  payload,
}: {
  payload: HomeComponentPayloadCategories;
}) {
  const excludedCategoryIds = [
    '101859',
    '207738',
    '500349',
    '566175',
    '564568',
    '568409',
  ];
  const displayableCategories = payload.categories.filter(
    (c) => !excludedCategoryIds.includes(c.id)
  );
  return (
    <div className="mx-3 my-4">
      <h2 className="pb-2 border-b-2 border-gray-200 dark:border-gray-800 text-xl font-bold">
        {payload.title}
      </h2>
      <div className="py-2">
        <div className="grid grid-cols-3 space-x-2">
          {displayableCategories.map((category) => (
            <Link
              key={category.id}
              href={`${routes.search()}?categoryIds=${category.id}&searchFrom=${
                SearchFrom.Home
              }`}
            >
              <a className="flex flex-col items-center justify-center my-1 border-[1px] border-gray-100 dark:border-gray-800 p-2 text-sm">
                {category.imageUrl && (
                  <img src={category.imageUrl} className="w-10 h-10 my-2" />
                )}
                {category.name}
              </a>
            </Link>
          ))}
        </div>
      </div>
    </div>
  );
}

function HomeComponentPosts({
  payload,
}: {
  payload: HomeComponentPayloadMediaPosts;
}) {
  return (
    <div className="mx-3 my-4">
      <h2 className="pb-2 border-b-2 border-gray-200 dark:border-gray-800 text-xl font-bold">
        {payload.title}
      </h2>
      <div className="grid sm:grid-cols-2 md:grid-cols-3 sm:gap-4 gap-y-4 sm:gap-y-8">
        {payload.posts.map((post) => (
          <PostCard
            key={post.slug}
            postMeta={{ ...post, mainImage: post.mainImageUrl }}
          />
        ))}
      </div>
      <Link href={routes.media()}>
        <a>
          <button
            type="button"
            className="block mx-auto my-4 px-2.5 py-3 w-[80%] rounded border-[1px] border-gray-800 dark:border-gray-200 text-center focus:outline-none"
          >
            記事をもっと見る
          </button>
        </a>
      </Link>
    </div>
  );
}
