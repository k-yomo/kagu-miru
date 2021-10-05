import React, { Fragment, memo } from 'react';
import { Menu, Transition } from '@headlessui/react';

interface Props {
  show: boolean;
  suggestedQueries: string[];
  onClickQuery: (query: string) => void;
}

function classNames(...classes: string[]) {
  return classes.filter(Boolean).join(' ');
}

export default memo(function QuerySuggestionsDropdown({
  show,
  suggestedQueries,
  onClickQuery,
}: Props) {
  return (
    <Menu as="div" className="relative block w-full text-left">
      <Transition
        as={Fragment}
        show={show}
        enter="transition ease-out duration-100"
        enterFrom="transform opacity-0"
        enterTo="transform opacity-100 scale-100"
        leave="transition ease-in duration-75"
        leaveFrom="transform opacity-100 scale-100"
        leaveTo="transform opacity-0"
      >
        <Menu.Items className="origin-top-right absolute left-0 mt-2 w-full rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 focus:outline-none">
          <div className="py-1 w-full">
            {suggestedQueries.map((query) => (
              <Menu.Item key={query}>
                {({ active }) => (
                  <span
                    onClick={() => onClickQuery(query)}
                    className={classNames(
                      active ? 'bg-gray-100 text-gray-900' : 'text-gray-700',
                      'block px-4 py-2 cursor-pointer text-sm'
                    )}
                  >
                    {query}
                  </span>
                )}
              </Menu.Item>
            ))}
          </div>
        </Menu.Items>
      </Transition>
    </Menu>
  );
});
