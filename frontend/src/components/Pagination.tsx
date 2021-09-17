import {
  ArrowNarrowLeftIcon,
  ArrowNarrowRightIcon,
} from '@heroicons/react/solid';

function range(start: number, end: number) {
  return Array(end - start + 1)
    .fill(0)
    .map((_, idx) => start + idx);
}

interface Props {
  page: number;
  totalPage: number;
  onClickPage: (page: number) => void;
}

export default function Pagination({ page, totalPage, onClickPage }: Props) {
  const minPage = Math.max(page - 5, 1);
  const maxPage = Math.min(page + 9 - (page - minPage), totalPage);

  const pages = [
    ...range(minPage, page - 1),
    page,
    ...range(page + 1, maxPage),
  ];

  return (
    <nav className="w-full border-t border-gray-200 dark:border-gray-400 px-4 flex items-center justify-between sm:px-0">
      <div className="-mt-px w-0 flex-1 flex">
        <button
          disabled={page === 1}
          onClick={() => onClickPage(page - 1)}
          className="border-t-2 disabled:border-none border-transparent pt-4 pr-1 inline-flex items-center text-sm ftext-gray-500 dark:text-gray-300 disabled:text-gray-300 dark:disabled:text-gray-500 hover:text-gray-700 dark:hover:text-gray-100 disabled:hover:text-gray-300 dark:disabled:hover:text-gray-500 disabled:hover:border-transparent"
        >
          <ArrowNarrowLeftIcon
            className={`mr-3 h-5 w-5 ${
              page === 1
                ? 'text-gray-300 dark:text-gray-700'
                : 'text-gray-500 dark:text-gray-300'
            }`}
            aria-hidden="true"
          />
          前へ
        </button>
      </div>
      <div className="hidden md:-mt-px md:flex">
        {pages.map((p) => (
          <span
            key={p}
            className={`${
              p === page
                ? 'border-black dark:border-white text-black dark:text-text-primary-dark font-bold'
                : 'border-transparent text-gray-400 hover:text-gray-700 dark:hover:text-gray-100 hover:border-gray-300'
            } cursor-pointer border-t-2 pt-4 px-4 inline-flex items-center text-sm`}
            onClick={() => onClickPage(p)}
          >
            {p}
          </span>
        ))}
      </div>
      <div className="-mt-px w-0 flex-1 flex justify-end">
        <button
          disabled={page === totalPage}
          onClick={() => onClickPage(page + 1)}
          className="border-t-2 disabled:border-none border-transparent pt-4 pr-1 inline-flex items-center text-sm text-gray-500 dark:text-gray-300 disabled:text-gray-300 dark:disabled:text-gray-500 hover:text-gray-700 dark:hover:text-gray-100 disabled:hover:text-gray-300 dark:disabled:hover:text-gray-500 disabled:hover:border-transparent"
        >
          次へ
          <ArrowNarrowRightIcon
            className={`ml-3 h-5 w-5 ${
              page === totalPage
                ? 'text-gray-500 dark:text-gray-700'
                : 'text-gray-500 dark:text-gray-300'
            }`}
            aria-hidden="true"
          />
        </button>
      </div>
    </nav>
  );
}
