import React, { ChangeEvent, memo } from 'react';
import { ItemSellingPlatform } from '@src/generated/graphql';
import { platFormText } from '@src/conv/platform';

interface Props {
  platforms: ItemSellingPlatform[];
  onChangePlatforms: (platforms: ItemSellingPlatform[]) => void;
}

export default memo(function PlatformSelect({
  platforms,
  onChangePlatforms,
}: Props) {
  const onChange = (e: ChangeEvent<HTMLInputElement>) => {
    const platform = e.target.name as ItemSellingPlatform;
    if (e.target.checked) {
      if (!platforms.includes(platform)) {
        onChangePlatforms([...platforms, platform]);
      }
    } else {
      onChangePlatforms(platforms.filter((p) => p !== platform));
    }
  };

  return (
    <div>
      <div className="space-y-2">
        <div className="flex items-center">
          <input
            id="platformCheckBoxRakuten"
            type="checkbox"
            name={ItemSellingPlatform.Rakuten}
            checked={platforms.includes(ItemSellingPlatform.Rakuten)}
            className="h-4 w-4 rounded cursor-pointer"
            onChange={onChange}
          />
          <label
            htmlFor="platformCheckBoxRakuten"
            className="ml-2 cursor-pointer text-sm"
          >
            {platFormText(ItemSellingPlatform.Rakuten)}
          </label>
        </div>
        <div className="flex items-center">
          <input
            id="platformCheckBoxYahooShopping"
            type="checkbox"
            name={ItemSellingPlatform.YahooShopping}
            checked={platforms.includes(ItemSellingPlatform.YahooShopping)}
            className="h-4 w-4 rounded cursor-pointer"
            onChange={onChange}
          />
          <label
            htmlFor="platformCheckBoxYahooShopping"
            className="ml-2 cursor-pointer text-sm"
          >
            {platFormText(ItemSellingPlatform.YahooShopping)}
          </label>
        </div>
      </div>
    </div>
  );
});
