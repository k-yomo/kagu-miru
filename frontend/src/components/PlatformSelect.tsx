import React, { ChangeEvent, memo } from 'react';
import { ItemSellingPlatform } from '@src/generated/graphql';
import { allPlatforms, platFormText } from '@src/conv/platform';

interface Props {
  platforms: ItemSellingPlatform[];
  onChangePlatforms: (platforms: ItemSellingPlatform[]) => void;
  // TODO: Use unstable_useOpaqueIdentifier
  // https://www.dkrk-blog.net/react/useopaqueidentifier
  htmlIdPrefix: string;
}

export default memo(function PlatformSelect({
  platforms,
  onChangePlatforms,
  htmlIdPrefix,
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
              id={`${htmlIdPrefix}_platformCheckBox_${platform}`}
              type="checkbox"
              name={platform}
              checked={platforms.includes(platform)}
              className="h-5 w-5 rounded cursor-pointer bg-gray-200 dark:bg-gray-800 border-none text-rose-500 focus:ring-0 form-checkbox"
              onChange={onChange}
            />
            <label
              htmlFor={`${htmlIdPrefix}_platformCheckBox_${platform}`}
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
