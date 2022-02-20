import React, { Fragment, useCallback, useRef, useState } from 'react';
import { Dialog, Transition } from '@headlessui/react';
import { SearchIcon, XIcon } from '@heroicons/react/outline';
import { SearchFrom } from '@src/generated/graphql';
import { buildSearchUrlQuery, useSearch } from '@src/contexts/search';
import SearchBar from '@src/components/SearchBar';
import { useRouter } from 'next/router';
import { routes } from '@src/routes/routes';

export default function SearchBarModal() {
  const router = useRouter();
  const { searchState, dispatch } = useSearch();
  const [open, setOpen] = useState(false);
  const cancelButtonRef = useRef(null);

  const onSubmitQuery = useCallback(
    (query: string, searchFrom: SearchFrom) => {
      router.push({
        pathname: routes.search(),
        query: buildSearchUrlQuery(query, searchFrom),
      });
      setOpen(false);
    },
    [dispatch]
  );

  return (
    <>
      <div>
        <div onClick={() => setOpen((prevState: boolean) => !prevState)}>
          <SearchIcon className="w-5 h-5" />
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
              enter="ease-out duration-200"
              enterFrom="opacity-0"
              enterTo="opacity-100"
              leave="ease-in duration-100"
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
              enter="ease-out duration-200"
              enterFrom="opacity-0 translate-y-8"
              enterTo="opacity-100 translate-y-0"
              leave="ease-in duration-200"
              leaveFrom="opacity-100 translate-y-0"
              leaveTo="opacity-0 translate-y-4"
            >
              <div className="flex flex-col w-screen h-screen overflow-hidden bg-white dark:bg-black transition-all transform">
                <div className="flex items-center space-x-2 px-4 py-4">
                  <SearchBar
                    query={searchState.searchInput.query}
                    onSubmit={onSubmitQuery}
                  />
                  <XIcon className="w-7 h-7" onClick={() => setOpen(false)} />
                </div>
              </div>
            </Transition.Child>
          </div>
        </Dialog>
      </Transition.Root>
    </>
  );
}
