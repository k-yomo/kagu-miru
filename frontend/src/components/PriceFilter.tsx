import React, { useState, memo } from 'react';

interface Props {
  defaultMinPrice?: number;
  defaultMaxPrice?: number;
  onSubmit: (minPrice?: number, maxPrice?: number) => void;
  onClear: () => void;
}

export default memo(function PriceFilter({
  defaultMinPrice,
  defaultMaxPrice,
  onSubmit,
  onClear,
}: Props) {
  const [minPrice, setMinPrice] = useState(defaultMinPrice);
  const [maxPrice, setMaxPrice] = useState(defaultMaxPrice);

  const onClickClear = () => {
    if (minPrice || maxPrice) {
      setMinPrice(undefined);
      setMaxPrice(undefined);
      onClear();
    }
  };

  const onClickApply = () => {
    onSubmit(
      minPrice && !isNaN(minPrice) ? minPrice : undefined,
      maxPrice && !isNaN(maxPrice) ? maxPrice : undefined
    );
  };

  return (
    <div>
      <div className="flex items-center text-xs">
        <div>
          <input
            type="text"
            inputMode="numeric"
            value={minPrice || ''}
            onChange={(e) => setMinPrice(parseInt(e.target.value))}
            className="w-[5rem] bg-white mr-1 p-1 dark:bg-gray-800 border border-gray-700 leading-5 text-gray-900 dark:text-gray-300 placeholder-gray-500 dark:placeholder-gray-400 focus:outline-none focus:ring-1 focus:ring-black dark:focus:ring-gray-400 text-xs"
          />
          円
        </div>
        {'　'}〜{'　'}
        <div>
          <input
            type="text"
            inputMode="numeric"
            value={maxPrice || ''}
            onChange={(e) => setMaxPrice(parseInt(e.target.value))}
            className="w-[5rem] bg-white mr-1 p-1 dark:bg-gray-800 border border-gray-700 leading-5 text-gray-900 dark:text-gray-300 placeholder-gray-500 dark:placeholder-gray-400 focus:outline-none focus:ring-1 focus:ring-black dark:focus:ring-gray-400 text-xs"
          />
          円
        </div>
      </div>
      <div className="flex items-center justify-end mt-4 space-x-2 text-sm">
        <button
          type="button"
          className="inline-flex items-center px-2.5 py-1.5 border border-black dark:border-white shadow-sm text-xs font-medium bg-white hover:bg-gray-50 dark:bg-black dark:hover:bg-gray-800 focus:outline-none"
          onClick={onClickClear}
        >
          クリア
        </button>
        <button
          type="button"
          className="inline-flex items-center px-2.5 py-1.5 border border-indigo-700 text-xs font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none"
          onClick={onClickApply}
        >
          適用
        </button>
      </div>
    </div>
  );
});
