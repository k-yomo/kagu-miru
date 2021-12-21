import React, { Fragment, useEffect, useRef, useState } from 'react';
import { Dialog, Transition } from '@headlessui/react';
import { AdjustmentsIcon, XIcon } from '@heroicons/react/outline';
import CategoryList from '@src/components/CategoryList';
import {
  ItemColor,
  ItemSellingPlatform,
  SearchFilter,
} from '@src/generated/graphql';
import RatingSelect from '@src/components/RatingSelect';
import {
  defaultSearchFilter,
  SearchActionType,
  useSearch,
} from '@src/contexts/search';
import PlatformSelect from '@src/components/PlatformSelect';
import ColorSelect from '@src/components/ColorSelect';

export default function MobileSearchFilterModal() {
  const { searchState, dispatch } = useSearch();
  const [searchFilter, setSearchFilter] = useState<SearchFilter>(
    searchState.searchInput.filter
  );
  const [open, setOpen] = useState(false);
  const cancelButtonRef = useRef(null);

  const setCategoryIds = (categoryIds: string[]) => {
    setSearchFilter((prevState: SearchFilter) => ({
      ...prevState,
      categoryIds,
    }));
  };

  const setPlatforms = (platforms: ItemSellingPlatform[]) => {
    setSearchFilter((prevState: SearchFilter) => ({ ...prevState, platforms }));
  };

  const setColors = (colors: ItemColor[]) => {
    setSearchFilter((prevState: SearchFilter) => ({ ...prevState, colors }));
  };
  const setMinPrice = (price: number) => {
    const minPrice = price && !isNaN(price) ? price : undefined;
    setSearchFilter((prevState: SearchFilter) => ({ ...prevState, minPrice }));
  };

  const setMaxPrice = (price: number) => {
    const maxPrice = price && !isNaN(price) ? price : undefined;
    setSearchFilter((prevState: SearchFilter) => ({ ...prevState, maxPrice }));
  };

  const setMinRating = (minRating?: number) => {
    setSearchFilter((prevState: SearchFilter) => ({ ...prevState, minRating }));
  };

  const onClickApply = () => {
    dispatch({ type: SearchActionType.SET_FILTER, payload: searchFilter });
    setOpen(false);
  };

  const onClickClear = () => {
    setSearchFilter(defaultSearchFilter);
  };

  useEffect(() => {
    setSearchFilter(searchState.searchInput.filter);
  }, [searchState.searchInput.filter]);

  return (
    <>
      <div className="sm:hidden z-10 fixed right-5 bottom-5">
        <div
          onClick={() => setOpen((prevState: boolean) => !prevState)}
          className="flex flex-col items-center px-2 py-2.5 bg-black dark:bg-white rounded-lg cursor-pointer text-white dark:text-black"
        >
          <AdjustmentsIcon className="w-6 h-6" />
          <span className="mt-1 text-xs tracking-tighter font-bold">
            絞り込み
          </span>
        </div>
      </div>
      <Transition.Root show={open} as={Fragment}>
        <Dialog
          initialFocus={cancelButtonRef}
          as="div"
          className="overflow-y-auto fixed inset-0 z-10"
          onClose={() => setOpen(false)}
        >
          <div className="flex justify-center items-center min-h-screen text-center">
            <Transition.Child
              as={Fragment}
              enter="ease-out duration-300"
              enterFrom="opacity-0"
              enterTo="opacity-100"
              leave="ease-in duration-200"
              leaveFrom="opacity-100"
              leaveTo="opacity-0"
            >
              <Dialog.Overlay className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" />
            </Transition.Child>

            {/* This element is to trick the browser into centering the modal contents. */}
            <span className="hidden" aria-hidden="true">
              &#8203;
            </span>
            <Transition.Child
              as={Fragment}
              enter="ease-out duration-300"
              enterFrom="opacity-0 translate-y-8"
              enterTo="opacity-100 translate-y-0"
              leave="ease-in duration-200"
              leaveFrom="opacity-100 translate-y-0"
              leaveTo="opacity-0 translate-y-4"
            >
              <div className="fixed flex flex-col bottom-0 w-screen h-[90vh] overflow-hidden bg-white dark:bg-black rounded-xl transition-all transform">
                <Dialog.Title
                  as="h3"
                  className="my-5 text-xl font-bold text-center"
                >
                  絞り込み検索
                </Dialog.Title>
                <div className="absolute top-0 right-0 mt-3 pr-3">
                  <button
                    type="button"
                    className="rounded-md text-gray-400"
                    onClick={() => setOpen(false)}
                  >
                    <span className="sr-only">Close</span>
                    <XIcon className="h-6 w-6" aria-hidden="true" />
                  </button>
                </div>
                <hr className="border-gray-100 dark:border-gray-800" />
                <div className="flex-1 px-6 mt-2 overflow-y-scroll text-left">
                  <h3 className="my-2 text-md font-bold">カテゴリー</h3>
                  <CategoryList
                    categoryIds={searchFilter.categoryIds}
                    showCategoryCount={5}
                    onClickCategory={(categoryId: string) =>
                      setCategoryIds([categoryId])
                    }
                    onClearCategory={() => setCategoryIds([])}
                  />
                  <hr className="my-3 border-gray-100 dark:border-gray-800" />
                  <h3 className="my-2 text-md font-bold">ECサイト</h3>
                  <PlatformSelect
                    platforms={searchFilter.platforms}
                    onChangePlatforms={setPlatforms}
                    htmlIdPrefix="mobile"
                  />
                  <hr className="my-3 border-gray-100 dark:border-gray-800" />
                  <h3 className="my-2 text-md font-bold">カラー</h3>
                  <ColorSelect
                    colors={searchFilter.colors}
                    onChangeColors={setColors}
                  />
                  <hr className="my-3 border-gray-100 dark:border-gray-800" />
                  <h3 className="my-2 text-md font-bold">価格</h3>
                  <div className="flex items-center">
                    <div>
                      <input
                        type="text"
                        inputMode="numeric"
                        value={searchFilter.minPrice || ''}
                        onChange={(e) => setMinPrice(parseInt(e.target.value))}
                        className="w-[5rem] bg-white mr-1 p-1 dark:bg-gray-800 border border-gray-700 leading-5 text-gray-900 dark:text-gray-300 placeholder-gray-500 dark:placeholder-gray-400 focus:outline-none focus:ring-1 focus:ring-black dark:focus:ring-gray-400"
                      />
                      円
                    </div>
                    {'　'}〜{'　'}
                    <div>
                      <input
                        type="text"
                        inputMode="numeric"
                        value={searchFilter.maxPrice || ''}
                        onChange={(e) => setMaxPrice(parseInt(e.target.value))}
                        className="w-[5rem] bg-white mr-1 p-1 dark:bg-gray-800 border border-gray-700 leading-5 text-gray-900 dark:text-gray-300 placeholder-gray-500 dark:placeholder-gray-400 focus:outline-none focus:ring-1 focus:ring-black dark:focus:ring-gray-400"
                      />
                      円
                    </div>
                  </div>
                  <hr className="my-3 border-gray-100 dark:border-gray-800" />
                  <h3 className="my-2 text-md font-bold">レビュー評価</h3>
                  <RatingSelect
                    curMinRating={searchFilter.minRating || undefined}
                    onChangeRating={setMinRating}
                  />
                  <div className="mb-2" />
                </div>
                <hr className="border-gray-100 dark:border-gray-800" />
                <div className="flex justify-center my-4 space-x-2 px-2">
                  <button
                    type="button"
                    className="w-full px-4 py-2 border border-black dark:border-white shadow-sm font-medium bg-white hover:bg-gray-50 dark:bg-black dark:hover:bg-gray-800 focus:outline-none"
                    onClick={onClickClear}
                  >
                    クリア
                  </button>
                  <button
                    type="button"
                    className="w-full px-4 py-2  bg-gradient-to-r from-pink-500 dark:from-pink-600 to-rose-500 dark:to-rose-600 text-white focus:outline-none"
                    onClick={onClickApply}
                  >
                    適用
                  </button>
                </div>
              </div>
            </Transition.Child>
          </div>
        </Dialog>
      </Transition.Root>
    </>
  );
}
