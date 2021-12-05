import React from 'react';
import { useItemDetailPageGetItemQuery } from '@src/generated/graphql';
import PlatformBadge from '@src/components/PlatformBadge';
import Rating from '@src/components/Rating';

interface Props {
  itemId: string;
}

export default function ItemDetailCard({ itemId }: Props) {
  const { data, loading, error } = useItemDetailPageGetItemQuery({
    variables: { id: itemId },
  });
  const item = data?.getItem;
  if (!item || loading) {
    return <div>loading...</div>;
  }
  if (error) {
    return (
      <div className="p-3 shadow-md rounded-md">商品の取得に失敗しました。</div>
    );
  }
  return (
    <a
      href={!!item.affiliateUrl ? item.affiliateUrl : item.url}
      className="no-underline"
    >
      <div className="flex items-center border-[1px] border-gray-200 dark:border-gray-800 rounded-md">
        <div className="w-[40%] sm:w-[30%] h-full overflow-hidden">
          <img
            src={item.imageUrls[0]}
            alt={item.name}
            className="w-full h-full object-cover object-center rounded-tl-md rounded-bl-md"
          />
        </div>
        <div className="ml-2 sm:ml-3 w-[60%] sm:w-[70%]">
          <PlatformBadge platform={item.platform} size="sm" />
          <div className="my-1 line-clamp-3 font-bold text-sm sm:text-lg">
            {item.name}
          </div>
          <div className="hidden sm:block">
            <span className="line-clamp-3 text-sm text-text-secondary dark:text-text-secondary-dark">
              {item.description}
            </span>
          </div>
          <div className="my-1 flex items-center">
            <Rating rating={item.averageRating} maxRating={5} />
            <div className="ml-1 text-sm text-gray-600 dark:text-gray-300">
              {item.reviewCount}
            </div>
          </div>
          <div className="my-2 text-xl sm:text-2xl font-bold">
            {item.price.toLocaleString()}円
          </div>
        </div>
      </div>
    </a>
  );
}
