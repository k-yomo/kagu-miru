import { ChangeEvent, KeyboardEvent, useCallback, useEffect, useState } from "react"
import type { NextPage } from 'next'
import Link from 'next/link'
import Image from 'next/image'
import gql from "graphql-tag"
import { SearchIcon } from '@heroicons/react/solid'
import { useHomePageSearchItemsLazyQuery } from "@src/generated/graphql"
import SEOMeta from '@src/components/SEOMeta'
import PageLoading from "@src/components/PageLoading"
import { useRouter } from "next/router"

gql`
    query homePageSearchItems($input: SearchItemsInput!) {
        searchItems(input: $input) {
            id
            name
            description
            status
            sellingPageURL
            price
            imageUrls
            averageRating
            reviewCount
            platform
        }
    }
`

const Home: NextPage = () => {
  const router = useRouter()
  const [searchQuery, setSearchQuery] = useState('')
  const [page, setPage] = useState<number>(1);
  const [searchItems, { data, loading, error }] = useHomePageSearchItemsLazyQuery({
    fetchPolicy: 'no-cache',
    nextFetchPolicy: 'no-cache',
  })

  const onChangeSearchInput = useCallback((e: ChangeEvent<HTMLInputElement>) => {
    setSearchQuery(e.target.value.trim())
  }, [setSearchQuery])

  const onSearchKeyPress = useCallback((e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key == 'Enter') {
      e.preventDefault()
      refreshPageWithParams()
    }
  }, [searchQuery])

  const refreshPageWithParams = () => {
    router.push(`${router.pathname}?q=${searchQuery}&page=${page}`, undefined, { shallow: true });
  }

  useEffect(() => {
    const page = parseInt(router.query.page as string) || 1;
    setPage(page)
    if (router.query.q) {
      const query = router.query.q as string
      setSearchQuery(query)
      searchItems(
        {
          variables: { input: { query, page } },
        },
      )
    }
  }, [router.query.q])

  return (
    <div>
      <SEOMeta
        title="カグミル - 家具検索サービス"
        excludeSiteTitle
        description="カグミルはオンラインで買える家具を横断で検索出来るサービスです。"
        // img={{ srcPath: TopImg.src }}
      />
      <div className="m-4">
        <h1 className="text-2xl text-black dark:text-white">商品検索</h1>
        <div className="my-4 max-w-xl w-full lg:max-w-lg">
          <div className="relative text-gray-400 focus-within:text-gray-600">
            <div className="pointer-events-none absolute inset-y-0 left-0 pl-3 flex items-center">
              <SearchIcon className="h-5 w-5" aria-hidden="true"/>
            </div>
            <input
              id="search"
              className="block w-full bg-white py-3 pl-10 pr-3 dark:bg-gray-800 border border-gray-700 rounded-sm leading-5 text-gray-900 dark:text-gray-300 placeholder-gray-500 dark:placeholder-gray-400 focus:outline-none focus:outline-none focus:ring-1 focus:ring-black dark:focus:ring-gray-400"
              placeholder="Search"
              type="search"
              name="search"
              value={searchQuery}
              onChange={onChangeSearchInput}
              onKeyPress={onSearchKeyPress}
            />
          </div>
        </div>
        {loading ? <PageLoading/> : <></>}
        <div className="flex flex-col items-center sm:m-6">
          <div className="relative grid grid-cols-3 md:grid-cols-6 gap-8 px-3 w-full">
            {
              data && data.searchItems.map(item => (
                <Link key={item.id} href={item.sellingPageURL}>
                  <a>
                    <div className="rounded-sm shadow">
                      <Image
                        src={item.imageUrls[0] || 'https://via.placeholder.com/300'}
                        alt={item.name}
                        width={300}
                        height={300}
                        layout="responsive"
                        objectFit="cover"
                        className="w-20 h-20"
                      />
                      <div className="p-2">
                        <h4 className="truncate">{item.name}</h4>
                        <span>￥{item.price}</span>
                        <span>{item.platform}</span>
                      </div>
                    </div>
                  </a>
                </Link>
              ))
            }
          </div>
        </div>
      </div>
    </div>
  )
}

export default Home
