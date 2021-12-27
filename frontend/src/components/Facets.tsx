import React, { Fragment, memo, useRef, useState } from 'react';
import { Dialog, Transition } from '@headlessui/react';
import { ChevronDownIcon } from '@heroicons/react/solid';
import { SearchActionType, useSearch } from '@src/contexts/search';
import { FacetType, ItemColor, SearchQuery } from '@src/generated/graphql';

export default memo(function Facets() {
  const { facets, searchState, dispatch } = useSearch();
  const [openFacetTitle, setOpenFacetTitle] = useState<string | undefined>();

  const getSelectedIds = (facetType: FacetType, title: string) => {
    switch (facetType) {
      case FacetType.CategoryIds:
        return searchState.searchInput.filter.categoryIds;
      case FacetType.BrandNames:
        return searchState.searchInput.filter.brandNames;
      case FacetType.Colors:
        return searchState.searchInput.filter.colors;
      case FacetType.Metadata:
        return (
          searchState.searchInput.filter.metadata.find((m) => m.name == title)
            ?.values || []
        );
      default:
        return [];
    }
  };

  const onClickFacet = (
    facetType: FacetType,
    selectedId: string,
    title: string
  ) => {
    switch (facetType) {
      case FacetType.CategoryIds:
        let categoryIds = searchState.searchInput.filter.categoryIds;
        if (categoryIds.includes(selectedId)) {
          categoryIds = categoryIds.filter((id) => id !== selectedId);
        } else {
          categoryIds = Array.from(new Set([...categoryIds, selectedId]));
        }
        dispatch({
          type: SearchActionType.SET_CATEGORY_FILTER,
          payload: categoryIds,
        });
        return;
      case FacetType.BrandNames:
        let brandNames = searchState.searchInput.filter.brandNames;
        if (brandNames.includes(selectedId)) {
          brandNames = brandNames.filter((name) => name !== selectedId);
        } else {
          brandNames = Array.from(new Set([...brandNames, selectedId]));
        }
        dispatch({
          type: SearchActionType.SET_BRAND_FILTER,
          payload: brandNames,
        });
        return;
      case FacetType.Colors:
        let colors = searchState.searchInput.filter.colors;
        if (colors.includes(selectedId as ItemColor)) {
          colors = colors.filter(
            (color) => color !== (selectedId as ItemColor)
          );
        } else {
          colors = Array.from(new Set([...colors, selectedId as ItemColor]));
        }
        dispatch({
          type: SearchActionType.SET_COLOR_FILTER,
          payload: colors,
        });
        return;
      case FacetType.Metadata:
        let metadata = searchState.searchInput.filter.metadata;
        const selectedMetadata = metadata.find((m) => m.name === title);
        // already selected
        if (selectedMetadata) {
          if (selectedMetadata.values.includes(selectedId)) {
            selectedMetadata.values = selectedMetadata.values.filter(
              (v) => v !== selectedId
            );
          } else {
            selectedMetadata.values = Array.from(
              new Set([...selectedMetadata.values, selectedId])
            );
          }
        } else {
          metadata = [...metadata, { name: title, values: [selectedId] }];
        }
        dispatch({
          type: SearchActionType.SET_METADATA_FILTER,
          payload: metadata,
        });
        return;
    }
  };

  return (
    <div className="flex w-[95vw] sm:w-[90%] space-x-2 overflow-auto whitespace-nowrap">
      {facets.map((facet) => {
        const selectedIds = getSelectedIds(facet.facetType, facet.title);
        return (
          <div key={facet.title}>
            <div>
              <button
                onClick={() => setOpenFacetTitle(facet.title)}
                className={`inline-flex items-center justify-center w-full rounded-full border ${
                  selectedIds.length > 0
                    ? 'border-rose-500'
                    : 'border-black dark:border-gray-200'
                }  px-2 py-1.5 bg-white dark:bg-black text-xs`}
              >
                {facet.title}
                <ChevronDownIcon className="ml-2 h-5 w-5" aria-hidden="true" />
              </button>
            </div>
            <FacetDropdown
              isOpen={facet.title === openFacetTitle}
              facet={facet}
              selectedIds={selectedIds}
              onClickFacet={onClickFacet}
              onClose={() => setOpenFacetTitle(undefined)}
            />
          </div>
        );
      })}
    </div>
  );
});

interface FacetDropdownProps {
  isOpen: boolean;
  facet: SearchQuery['search']['facets'][number];
  selectedIds: string[];
  onClickFacet: (facetType: FacetType, id: string, title: string) => void;
  onClose: () => void;
}

function FacetDropdown({
  isOpen,
  facet,
  selectedIds,
  onClickFacet,
  onClose,
}: FacetDropdownProps) {
  const cancelButtonRef = useRef(null);
  const selectedIdMap: { [key: string]: boolean } = selectedIds.reduce(
    (m: { [key: string]: boolean }, v) => ((m[v] = true), m),
    {}
  );

  return (
    <>
      <Transition.Root show={isOpen} as={Fragment}>
        <Dialog
          initialFocus={cancelButtonRef}
          as="div"
          className="overflow-y-auto fixed inset-0 z-10"
          onClose={onClose}
        >
          <div className="flex justify-center items-center">
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
              <div className="fixed flex flex-col bottom-0 w-screen h-[70vh] overflow-y-auto bg-white dark:bg-black rounded-t-xl transition-all transform">
                <Dialog.Title as="h3" className="mt-4 ml-4 text-xl font-bold">
                  {facet.title}
                </Dialog.Title>
                <div className="py-1 divide-y-2">
                  {facet.values.map((facetValue) => (
                    <div
                      key={facetValue.id}
                      className="flex items-center flex-between px-4 py-2 border-gray-100 dark:border-gray-800 text-sm"
                      onClick={() =>
                        onClickFacet(
                          facet.facetType,
                          facetValue.id,
                          facet.title
                        )
                      }
                    >
                      <input
                        type="checkbox"
                        name={facetValue.name}
                        checked={selectedIdMap[facetValue.id] || false}
                        readOnly
                        className="h-5 w-5 rounded cursor-pointer bg-gray-200 dark:bg-gray-800 border-none text-rose-500 focus:ring-0 form-checkbox"
                      />
                      <div className="ml-2 cursor-pointer text-sm">
                        <div>{facetValue.name}</div>
                        <div className="text-xs text-text-secondary dark:text-text-secondary-dark">
                          {facetValue.count.toLocaleString()}件
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </Transition.Child>
          </div>
        </Dialog>
      </Transition.Root>
    </>
  );
}
