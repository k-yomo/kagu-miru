import React, { memo, useState } from "react"
import gql from "graphql-tag"
import {
  ItemPreviewGetItemDocument,
  ItemPreviewGetItemQuery,
} from "../src/generated/graphql"
import apolloClient from "../lib/apolloClient"
import { GraphQLError } from "graphql"

gql`
  query itemPreviewGetItem($id: ID!) {
      getItem(id: $id) {
          id
          name
          status
          url
          affiliateUrl
          price
          imageUrls
          averageRating
          reviewCount
          categoryIds
          platform
      }
  }
`

export default memo(function ({ value }: { value: { id: string } }) {
  const [item, setItem] = useState<ItemPreviewGetItemQuery['getItem']>()
  const [error, setError] = useState<GraphQLError>()
  const { id } = value
  apolloClient.query<ItemPreviewGetItemQuery>({
    query: ItemPreviewGetItemDocument,
    variables: { id }
  }).then(({ data }) => {
    console.log(data.getItem)
    setItem(data.getItem)
  }).catch(e => {
    setError(e)
  })
  if (error) {
    return <div style={{ color: "red"}}>{error.name}</div>
  }
  if (!item) {
    return <div>loading...</div>
  }
  return (
    <div style={{ display: "flex", alignItems: "center", padding: 5 }}>
      <img src={item.imageUrls[0]} alt={item.name} style={{ width: 50, height: 50 }} />
      {item.name}
    </div>
  )
})