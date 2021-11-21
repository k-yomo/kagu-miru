import React, { memo } from 'react';
import { ItemSellingPlatform } from '@src/generated/graphql';
import { platFormColor, platFormText } from '@src/conv/platform';

interface Props {
  platform: ItemSellingPlatform;
  size?: 'xs' | 'sm' | 'md';
}

export default memo(function PlatformBadge({ platform, size }: Props) {
  const text = platFormText(platform, true);
  const color = platFormColor(platform);

  size = size || 'md';
  let textSize = `text-${size}`;
  let heightWidth = 'h-3 w-3';
  if (size === 'xs' || size == 'sm') {
    heightWidth = 'h-2 w-2';
  }
  return (
    <span className={`inline-flex items-center rounded ${textSize} font-bold`}>
      <svg
        className={`mr-0.5 ${heightWidth} ${color}`}
        fill="currentColor"
        viewBox="0 0 8 8"
      >
        <circle cx={4} cy={4} r={3} />
      </svg>
      {text}
    </span>
  );
});
