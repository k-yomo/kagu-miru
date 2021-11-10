import React, { memo } from 'react';
import { ItemSellingPlatform } from '@src/generated/graphql';

interface Props {
  platform: ItemSellingPlatform;
}

function platFormText(platform: ItemSellingPlatform) {
  switch (platform) {
    case ItemSellingPlatform.Rakuten:
      return '楽天';
    case ItemSellingPlatform.YahooShopping:
      return 'Yahoo';
  }
}

function platFormColor(platform: ItemSellingPlatform) {
  switch (platform) {
    case ItemSellingPlatform.Rakuten:
      return 'text-rakuten';
    case ItemSellingPlatform.YahooShopping:
      return 'text-yahoo-shopping';
  }
}

export default memo(function PlatformBadge({ platform }: Props) {
  const text = platFormText(platform);
  const color = platFormColor(platform);
  return (
    <span className="inline-flex items-center rounded text-xs font-bold">
      <svg
        className={`mr-0.5 h-2 w-2 ${color}`}
        fill="currentColor"
        viewBox="0 0 8 8"
      >
        <circle cx={4} cy={4} r={3} />
      </svg>
      {text}
    </span>
  );
});
