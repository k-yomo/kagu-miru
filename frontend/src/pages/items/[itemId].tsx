import React from 'react';
import gql from 'graphql-tag';
import Head from 'next/head';
import { GetServerSideProps } from 'next';
import apolloClient from '@src/lib/apolloClient';
import {
  ItemDetailPageGetItemDocument,
  ItemDetailPageGetItemQuery,
} from '@src/generated/graphql';

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
  console.log(data.getItem);
  ctx.res.setHeader('Cache-Control', 'public, max-age=3600');
  return { props: { item: data.getItem } };
};

interface Props {
  item: ItemDetailPageGetItemQuery['getItem'];
}

// TODO: complete the page
export default function ItemDetailPage({ item }: Props) {
  return (
    <div>
      <Head>
        {/* TODO: remove no index when the page is ready */}
        <meta name="robots" content="noindex,nofollow,noarchive" />
      </Head>
      <h1>{item.name}</h1>
    </div>
  );
}
