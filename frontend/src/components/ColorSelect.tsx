import { ItemColor } from '@src/generated/graphql';
import React, { memo } from 'react';
import { allColors, bgColorCSS } from '@src/conv/color';
import { CheckIcon } from '@heroicons/react/outline';

interface Props {
  colors: ItemColor[];
  onChangeColors: (platforms: ItemColor[]) => void;
}

export default memo(function ColorSelect({ colors, onChangeColors }: Props) {
  const onClick = (color: ItemColor) => {
    if (!colors.includes(color)) {
      onChangeColors([...colors, color]);
    } else {
      onChangeColors(colors.filter((p) => p !== color));
    }
  };

  return (
    <div>
      <div className="flex flex-col grid grid-cols-6">
        {allColors.map((color) => {
          return (
            <div
              key={color}
              onClick={() => onClick(color)}
              className={`flex items-center justify-center w-full h-6 border-[1px] border-gray-100 dark:border-gray-800 rounded-sm cursor-pointer ${bgColorCSS(
                color
              )}`}
            >
              {colors.includes(color) ? (
                <CheckIcon
                  className={`w-5 h-5 ${
                    [ItemColor.White, ItemColor.Transparent].includes(color)
                      ? 'text-black'
                      : 'text-white'
                  }`}
                />
              ) : (
                <></>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
});
