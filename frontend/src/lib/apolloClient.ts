import { ApolloClient, ApolloLink, createHttpLink, InMemoryCache } from '@apollo/client'
import { withScalars } from 'apollo-link-scalars'
import { buildClientSchema, IntrospectionQuery } from 'graphql'
import { DateTimeResolver } from 'graphql-scalars'
import fetch from 'isomorphic-unfetch'
import introspectionResult from '../../graphql.schema.json'
import { GRAPHQL_API_URL } from '@src/config/env'

const schema = buildClientSchema(
  introspectionResult as unknown as IntrospectionQuery,
)
const typesMap = {
  Date: DateTimeResolver,
}
const link = ApolloLink.from([
  withScalars({ schema, typesMap }) as unknown as ApolloLink,
  createHttpLink({
    uri: GRAPHQL_API_URL,
    credentials: 'include',
    fetch,
    headers: {
      'X-Requested-By': 'kagu-miru-frontend', // for CSRF validation
    },

  }),
])

const apolloClient = new ApolloClient({
  ssrMode: true,
  link,
  cache: new InMemoryCache(),
  defaultOptions: {
    watchQuery: {
      fetchPolicy: 'no-cache',
    },
  },
})

export default apolloClient
