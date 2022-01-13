import React, { Fragment, useCallback, useEffect, useState } from 'react';
import gql from 'graphql-tag';
import Head from 'next/head';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { GetServerSideProps } from 'next';
import { ApolloError } from '@apollo/client';
import { Swiper, SwiperSlide } from 'swiper/react';
import apolloClient, { isErrorIncludes } from '@src/lib/apolloClient';
import {
  Action,
  ErrorCode,
  EventId,
  ItemDetailPageGetItemDocument,
  ItemDetailPageGetItemQuery,
  ItemSellingPlatform,
  SearchClickItemActionParams,
  SimilarItemsDisplayItemsActionParams,
  useItemDetailPageGetSimilarItemsLazyQuery,
  useTrackEventMutation,
} from '@src/generated/graphql';
import ItemList from '@src/components/ItemList';
import Loading from '@src/components/Loading';
import { changeItemImageSize } from '@src/lib/platformImage';
import PlatformBadge from '@src/components/PlatformBadge';
import Rating from '@src/components/Rating';
import SEOMeta from '@src/components/SEOMeta';
import {
  findCategoryIdsById,
  findCategoryNameById,
} from '@src/lib/itemCategories';
import { ChevronRightIcon } from '@heroicons/react/solid';
import { routes } from '@src/routes/routes';

import SwiperCore, { FreeMode, Pagination, Thumbs, Zoom } from 'swiper';
import 'swiper/css';
import 'swiper/css/free-mode';
import 'swiper/css/navigation';
import 'swiper/css/thumbs';
import 'swiper/css/zoom';

SwiperCore.use([FreeMode, Pagination, Thumbs, Zoom]);

gql`
  query itemDetailPageGetItem($id: ID!) {
    getItem(id: $id) {
      id
      name
      description
      status
      url
      affiliateUrl
      price
      imageUrls
      averageRating
      reviewCount
      categoryId
      platform
    }
  }

  query itemDetailPageGetSimilarItems($input: GetSimilarItemsInput!) {
    getSimilarItems(input: $input) {
      searchId
      itemConnection {
        pageInfo {
          page
          totalPage
          totalCount
        }
        nodes {
          ...itemListItemFragment
        }
      }
    }
  }
`;

export const getServerSideProps: GetServerSideProps = async (ctx) => {
  try {
    const { data } = await apolloClient.query<ItemDetailPageGetItemQuery>({
      query: ItemDetailPageGetItemDocument,
      variables: { id: ctx.query.itemId as string },
    });

    ctx.res.setHeader(
      'Cache-Control',
      'public, max-age=600, stale-while-revalidate=86400'
    );
    return { props: { item: data.getItem } };
  } catch (e) {
    if (e instanceof ApolloError && isErrorIncludes(e, ErrorCode.NotFound)) {
      return { notFound: true };
    }
    throw e;
  }
};

interface Props {
  item: ItemDetailPageGetItemQuery['getItem'];
}

export default function ItemDetailPage({ item }: Props) {
  const router = useRouter();
  const [thumbsSwiper, setThumbsSwiper] = useState<any>(null);
  const [getSimilarItems, { data, loading }] =
    useItemDetailPageGetSimilarItemsLazyQuery({
      fetchPolicy: 'no-cache',
      nextFetchPolicy: 'no-cache',
      onCompleted: (data) => {
        const params: SimilarItemsDisplayItemsActionParams = {
          searchId: data.getSimilarItems.searchId,
          getSimilarItemsInput: { itemId: item.id, page: 0, pageSize: 100 },
          itemIds: data.getSimilarItems.itemConnection.nodes.map(
            (item) => item.id
          ),
        };
        trackEvent({
          variables: {
            event: {
              id: EventId.Search,
              action: Action.Display,
              createdAt: new Date(),
              params,
            },
          },
        }).catch(() => {
          // do nothing
        });
      },
    });
  const [trackEvent] = useTrackEventMutation();
  const similarItems = data?.getSimilarItems?.itemConnection.nodes;

  const mainImgUrl = changeItemImageSize(item.imageUrls[0], item.platform, 512);
  // Because of implementation error, yahoo shopping item has two dup images (medium and small)
  // So restricting images to remove the dup.
  // TODO: remove this condition after reindex
  const imageUrls = (
    [
      ItemSellingPlatform.YahooShopping,
      ItemSellingPlatform.PaypayMall,
    ] as ItemSellingPlatform[]
  ).includes(item.platform)
    ? [item.imageUrls[0]]
    : item.imageUrls;

  const onClickItem = useCallback(
    (itemId: string) => {
      const params: SearchClickItemActionParams = {
        searchId: data!.getSimilarItems.searchId,
        itemId,
      };
      trackEvent({
        variables: {
          event: {
            id: EventId.SimilarItems,
            action: Action.ClickItem,
            createdAt: new Date(),
            params,
          },
        },
      }).catch(() => {
        // do nothing
      });
    },
    [data, trackEvent]
  );

  useEffect(() => {
    getSimilarItems({
      variables: { input: { itemId: item.id, page: 0, pageSize: 100 } },
    }).catch((e) => {
      console.log(e);
    });
  }, [getSimilarItems, item.id]);

  return (
    <>
      <Head>
        <SEOMeta
          title={item.name}
          description={`${item.name}の詳細ページです。最安値のECサイトや関連商品の表示あり！`}
          img={{ src: mainImgUrl }}
          path={router.asPath}
        />
      </Head>
      <div className="max-w-[1200px] mx-auto mb-6">
        <div className="relative w-full h-[350px] sm:h-[600px]">
          <Swiper thumbs={{ swiper: thumbsSwiper }} zoom pagination>
            {imageUrls.map((imageUrl, i) => (
              <SwiperSlide key={imageUrl}>
                <div className="swiper-zoom-container">
                  <img
                    src={changeItemImageSize(imageUrl, item.platform, 600)}
                    alt={`${item.name} 商品画像 ${i + 1}`}
                    className="w-full h-[300px] sm:h-[500px] object-contain object-center"
                  />
                </div>
              </SwiperSlide>
            ))}
          </Swiper>
          <Swiper
            onSwiper={setThumbsSwiper}
            spaceBetween={5}
            slidesPerView={6}
            freeMode
          >
            {imageUrls.map((imageUrl, i) => (
              <SwiperSlide
                key={imageUrl}
                className="max-w-[50px] sm:max-w-[100px]"
              >
                <img
                  src={changeItemImageSize(imageUrl, item.platform, 512)}
                  alt={`${item.name} プレビュー画像 ${i + 1}`}
                  className="w-[50px] sm:w-[100px] h-[50px] sm:h-[100px] object-cover object-center"
                />
              </SwiperSlide>
            ))}
          </Swiper>
        </div>
        <div className="mx-3">
          <div className="mt-1">
            <span className="text-xs font-bold text-text-secondary dark:text-text-secondary-dark">
              カテゴリー
            </span>
            <div className="flex items-center text-sm">
              {findCategoryIdsById(item.categoryId).map((categoryId, i) => (
                <Fragment key={categoryId}>
                  {i !== 0 && (
                    <ChevronRightIcon className="w-5 h-5 text-text-secondary dark:text-text-secondary-dark" />
                  )}
                  <Link href={`${routes.top()}?categoryIds=${categoryId}`}>
                    <a className="block text-primary-500 cursor-pointer">
                      {findCategoryNameById(categoryId)}
                    </a>
                  </Link>
                </Fragment>
              ))}
            </div>
          </div>
          <h1 className="my-2 text-lg font-bold">{item.name}</h1>
          <div className="my-2">
            <PlatformBadge platform={item.platform} size="md" />
            <div className="flex items-center">
              <Rating rating={item.averageRating} maxRating={5} />
              <div className="ml-1 text-gray-600 dark:text-gray-300">
                {item.reviewCount}件
              </div>
            </div>
          </div>
          <div className="my-4 text-xl font-bold">
            <span className="text-text-secondary dark:text-text-secondary-dark">
              価格:
            </span>{' '}
            <span className="text-3xl text-rose-600 dark:text-rose-500">
              {item.price.toLocaleString()}円
            </span>
          </div>
          <hr className="border-gray-100 dark:border-gray-800" />
          <p
            dangerouslySetInnerHTML={{ __html: item.description }}
            className="my-4 text-text-secondary dark:text-text-secondary-dark"
          />
        </div>
        <a href={item.affiliateUrl}>
          <button
            type="button"
            className="block mx-auto my-4 px-2.5 py-3 w-[95%] rounded  bg-gradient-to-r from-pink-500 dark:from-pink-600 to-rose-500 dark:to-rose-600 text-center text-white focus:outline-none"
          >
            商品購入ページへ
          </button>
        </a>
        <hr className="border-gray-100 dark:border-gray-800" />
        <div className="mx-3">
          <h2 className="my-2 text-2xl">関連商品</h2>
          {loading ? <Loading /> : <></>}
          <div className="flex flex-col items-center">
            <div className="relative grid grid-cols-2 sm:grid-cols-4 md:grid-cols-5 gap-3 md:gap-4 text-sm sm:text-md">
              {similarItems && (
                <ItemList
                  items={similarItems}
                  onClickItem={onClickItem}
                  isAdmin={false}
                />
              )}
            </div>
            {/*{data?.getSimilarItems?.itemConnection?.pageInfo && (*/}
            {/*  <div className="my-4 w-full">*/}
            {/*    <Pagination />*/}
            {/*  </div>*/}
            {/*)}*/}
          </div>
        </div>
      </div>
    </>
  );
}
