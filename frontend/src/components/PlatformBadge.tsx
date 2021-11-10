import React, { memo } from 'react';
import { ItemSellingPlatform } from '@src/generated/graphql';
import { platFormColor, platFormText } from '@src/conv/platform';

interface Props {
  platform: ItemSellingPlatform;
}

export default memo(function PlatformBadge({ platform }: Props) {
  const text = platFormText(platform, true);
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
