import { ItemListItemFragmentFragment } from '@src/generated/graphql';
import React, { memo, PropsWithChildren } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import gql from 'graphql-tag';
import PlatformBadge from '@src/components/PlatformBadge';
import Rating from '@src/components/Rating';
import { routes } from '@src/routes/routes';

gql`
  fragment itemListItemFragment on Item {
    id
    name
    categoryId
    url
    affiliateUrl
    price
    imageUrls
    averageRating
    reviewCount
    platform
  }
`;

interface Props {
  items: Array<ItemListItemFragmentFragment>;
  onClickItem: (itemId: string) => void;
  isAdmin: boolean;
}

export default memo(function ItemList({ items, onClickItem, isAdmin }: Props) {
  const ItemWrapper = function ({
    id,
    url,
    children,
  }: PropsWithChildren<{ id: string, url: string }>) {
    if (isAdmin) {
      return <button onClick={() => onClickItem(id)}>{children}</button>;
    } else {
      return (
        <Link href={url}>
          <a onClick={() => onClickItem(id)}>{children}</a>
        </Link>
      );
    }
  };
  return (
    <>
      {items.map((item) => (
        <ItemWrapper key={item.id} id={item.id} url={item.affiliateUrl}>
          <div className="min-w-[140px] rounded-md sm:shadow">
            <Image
              src={item.imageUrls[0] || 'https://via.placeholder.com/300'}
              alt={item.name}
              width={300}
              height={300}
              layout="responsive"
              objectFit="cover"
              className="w-20 h-20 rounded-t-md"
              unoptimized
            />
            <div className="py-0.5 sm:p-2">
              <PlatformBadge platform={item.platform} size="xs" />
              <div className="flex items-center">
                <Rating rating={item.averageRating} maxRating={5} />
                <div className="ml-1 text-xs">{item.reviewCount}件</div>
              </div>
              <h4 className="mt-1 break-all line-clamp-2 text-sm sm:text-md">
                {item.name}
              </h4>
              <div className="text-lg font-bold text-black dark:text-white">
                {item.price.toLocaleString()}円
              </div>
            </div>
          </div>
        </ItemWrapper>
      ))}
    </>
  );
});
