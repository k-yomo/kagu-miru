import React from 'react';
import Link from 'next/link';
import { routes } from '@src/routes/routes';

export default function Footer() {
  return (
    <footer className="bg-black pt-5 dark:border-t dark:border-gray-800">
      <div className="max-w-6xl m-auto text-gray-800 flex flex-wrap justify-left">
        <div className="p-5 w-1/2 sm:w-4/12 md:w-3/12">
          <div className="text-xs uppercase text-gray-400  mb-6">
            クレジット
          </div>

          <div className="my-3 block text-gray-300 hover:text-gray-100 text-sm duration-700">
            <a href="https://developers.rakuten.com/" target="_blank">
              Supported by Rakuten Developers
            </a>
          </div>
          <div className="my-3 block text-gray-300 hover:text-gray-100 text-sm duration-700">
            <a href="https://developer.yahoo.co.jp/sitemap/" target="_blank">
              Web Services by Yahoo! JAPAN
            </a>
          </div>
        </div>
        <div className="p-5 w-1/2 sm:w-4/12 md:w-3/12">
          <div className="text-xs uppercase text-gray-400  mb-6">
            ヘルプ＆ガイド
          </div>

          <Link href={routes.contact()}>
            <a className="my-3 block text-gray-300 hover:text-gray-100 text-sm duration-700">
              お問い合せ
            </a>
          </Link>
        </div>

        <div className="p-5 w-1/2 sm:w-4/12 md:w-3/12">
          <div className="text-xs uppercase text-gray-400  mb-6">
            プライバシーと利用規約
          </div>

          <Link href={routes.privacyPolicy()}>
            <a className="my-3 block text-gray-300 hover:text-gray-100 text-sm duration-700">
              プライバシーポリシー
            </a>
          </Link>
        </div>

        <div className="p-5 w-1/2 sm:w-4/12 md:w-3/12">
          <div className="text-xs uppercase text-gray-400  mb-6">SNS</div>

          <a
            href="https://www.instagram.com/kagu_miru_official/"
            className="my-3 block text-gray-300 hover:text-gray-100 text-sm duration-700"
          >
            Instagram
          </a>
          <a
            href="https://twitter.com/kagu_miru"
            className="my-3 block text-gray-300 hover:text-gray-100 text-sm duration-700"
          >
            Twitter
          </a>
        </div>
      </div>

      <div className="pt-2">
        <div
          className="flex pb-5 px-3 m-auto pt-5
            border-t border-gray-800 text-gray-400 text-sm
            flex-col md:flex-row max-w-6xl"
        >
          <div className="mt-2">© 2021 kagu-miru.com</div>
        </div>
      </div>
    </footer>
  );
}
