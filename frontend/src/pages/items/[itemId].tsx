import React, { useEffect } from 'react';
import gql from 'graphql-tag';
import Head from 'next/head';
import { GetServerSideProps } from 'next';
import apolloClient from '@src/lib/apolloClient';
import {
  Action,
  EventId,
  ItemDetailPageGetItemDocument,
  ItemDetailPageGetItemQuery,
  SearchClickItemActionParams,
  SimilarItemsDisplayItemsActionParams,
  useItemDetailPageGetSimilarItemsLazyQuery,
  useTrackEventMutation,
} from '@src/generated/graphql';
import ItemList from '@src/components/ItemList';
import Loading from '@src/components/Loading';
import Pagination from '@src/components/Pagination';

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
  const { data, errors } = await apolloClient.query<ItemDetailPageGetItemQuery>(
    {
      query: ItemDetailPageGetItemDocument,
      variables: { id: ctx.query.itemId as string },
    }
  );
  if (errors) {
    return { notFound: true };
  }
  ctx.res.setHeader('Cache-Control', 'public, max-age=3600');
  return { props: { item: data.getItem } };
};

interface Props {
  item: ItemDetailPageGetItemQuery['getItem'];
}

// TODO: complete the page
export default function ItemDetailPage({ item }: Props) {
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

  const onClickItem = (itemId: string) => {
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
  };

  useEffect(() => {
    getSimilarItems({
      variables: { input: { itemId: item.id, page: 0, pageSize: 100 } },
    }).catch((e) => {
      console.log(e);
    });
  }, [getSimilarItems]);

  return (
    <>
      <Head>
        {/* TODO: remove no index when the page is ready */}
        <meta name="robots" content="noindex,nofollow,noarchive" />
      </Head>
      <div className="max-w-[1200px] mx-auto mt-3 mb-6">
        <h1>{item.name}</h1>
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
