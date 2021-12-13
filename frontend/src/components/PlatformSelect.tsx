import React, { ChangeEvent, memo } from 'react';
import { ItemSellingPlatform } from '@src/generated/graphql';
import { allPlatforms, platFormText } from '@src/conv/platform';

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
        {allPlatforms.map((platform) => (
          <div key={platform} className="flex items-center">
            <input
              id={`platformCheckBox_${platform}`}
              type="checkbox"
              name={platform}
              checked={platforms.includes(platform)}
              className="h-4 w-4 rounded cursor-pointer"
              onChange={onChange}
            />
            <label
              htmlFor={`platformCheckBox_${platform}`}
              className="ml-2 cursor-pointer text-sm"
            >
              {platFormText(platform)}
            </label>
          </div>
        ))}
      </div>
    </div>
  );
});
