import React, { memo } from 'react';
import { ImStarEmpty, ImStarHalf, ImStarFull } from 'react-icons/im';

interface Props {
  size?: number;
  rating: number;
  maxRating: number;
}

export default memo(function Rating({ size, rating, maxRating }: Props) {
  let fullStarCount = Math.floor(rating);
  let halfStarCount = 0;
  if (rating % 1 !== 0) {
    if (rating % 1 >= 0.8) {
      fullStarCount += 1;
    } else if (rating % 1 >= 0.3) {
      halfStarCount = 1;
    }
  }
  const emptyStarCount = maxRating - fullStarCount - halfStarCount;

  return (
    <span className="flex sm:space-x-[0.05rem] text-amber-400">
      {Array(fullStarCount).fill(<ImStarFull size={size} />)}
      {Array(halfStarCount).fill(<ImStarHalf size={size} />)}
      {Array(emptyStarCount).fill(<ImStarEmpty size={size} />)}
    </span>
  );
});
