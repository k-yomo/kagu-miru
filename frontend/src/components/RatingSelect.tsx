import React, { memo } from 'react';
import Rating from '@src/components/Rating';

interface Props {
  curMinRating?: number;
  onChangeRating: (rating?: number) => void;
}

export default memo(function RatingSelect({
  curMinRating,
  onChangeRating,
}: Props) {
  return (
    <div>
      <div className="space-y-3 text-xs">
        <div
          className={`cursor-pointer flex items-center hover:font-bold ${
            curMinRating === 5 ? 'font-bold' : ''
          }`}
          onClick={() => onChangeRating(5)}
        >
          <Rating
            size={20}
            rating={5}
            maxRating={5}
            grayOut={!!(curMinRating && curMinRating !== 5)}
          />
          以上
        </div>
        <div
          className={`cursor-pointer flex items-center hover:font-bold ${
            curMinRating === 4 ? 'font-bold' : ''
          }`}
          onClick={() => onChangeRating(4)}
        >
          <Rating
            size={20}
            rating={4}
            maxRating={5}
            grayOut={!!(curMinRating && curMinRating !== 4)}
          />
          以上
        </div>
        <div
          className={`cursor-pointer flex items-center hover:font-bold ${
            curMinRating === 3 ? 'font-bold' : ''
          }`}
          onClick={() => onChangeRating(3)}
        >
          <Rating
            size={20}
            rating={3}
            maxRating={5}
            grayOut={!!(curMinRating && curMinRating !== 3)}
          />
          以上
        </div>
      </div>
      <div className="flex items-center justify-end mt-2 text-sm">
        <button
          type="button"
          className="inline-flex items-center px-2.5 py-1.5 border border-black dark:border-white shadow-sm text-xs font-medium bg-white hover:bg-gray-50 dark:bg-black dark:hover:bg-gray-800 focus:outline-none"
          onClick={() => onChangeRating()}
        >
          クリア
        </button>
      </div>
    </div>
  );
});
