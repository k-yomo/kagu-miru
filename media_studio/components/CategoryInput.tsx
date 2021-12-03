import React, { forwardRef, ReactNode, useEffect, useState } from 'react';
import gql from 'graphql-tag';
import { Autocomplete, Box, Card, Label, Stack, Text } from '@sanity/ui';
import PatchEvent, { set, unset } from '@sanity/form-builder/PatchEvent';
import apolloClient from '../lib/apolloClient';
import {
  CategoryInputGetAllCategoriesDocument,
  CategoryInputGetAllCategoriesQuery,
} from '../generated/graphql';

gql`
  fragment subItemCategory on ItemCategory {
    id
    level
    name
    parentId
  }
  # max depth is 4 in our categories
  query categoryInputGetAllCategories {
    getAllItemCategories {
      ...subItemCategory
      children {
        ...subItemCategory
        children {
          ...subItemCategory
          children {
            ...subItemCategory
          }
        }
      }
    }
  }
`;

interface Props {
  value: {
    id?: string;
    names?: string[];
  };
  onChange: any;
}

// TODO: Fix type
export default forwardRef<HTMLInputElement, Props>((props: Props, ref) => {
  const [categories, setCategories] = useState<
    Array<{ id: string; names: string[] }>
  >([]);

  useEffect(() => {
    apolloClient
      .query<CategoryInputGetAllCategoriesQuery>({
        query: CategoryInputGetAllCategoriesDocument,
      })
      .then(({ data }) => {
        const idNamesMap = categoryIDNamesMap(
          data.getAllItemCategories,
          {},
          []
        );
        const topLevelCategoryMap: { [key: string]: boolean } =
          Object.fromEntries(
            data['getAllItemCategories'].map(({ id }) => [id, true])
          );
        // we don't allow to set top level category
        const categories = Object.entries(idNamesMap)
          .filter(([id]) => !topLevelCategoryMap[id])
          .map(([id, names]) => ({
            id,
            names,
          }));
        setCategories(categories);
      });
  }, []);

  const handleChange = React.useCallback(
    (selectedId: string) => {
      if (!selectedId) {
        props.onChange(PatchEvent.from(unset()));
      } else {
        const category = {
          _key: selectedId,
          id: selectedId,
          names: categories.find(({ id }) => id === selectedId)!.names,
        };
        props.onChange(PatchEvent.from(set(category)));
      }
    },
    [props.onChange, categories]
  );
  return (
    <Stack space={2}>
      <Label>カテゴリー</Label>
      <Autocomplete
        id="category_input"
        ref={ref}
        value={props.value?.id}
        filterOption={(query, option) =>
          option.payload.names
            .join(' > ')
            .toLowerCase()
            .indexOf(query.toLowerCase()) > -1
        }
        options={categories.map((category) => ({
          value: category.id,
          payload: { names: category.names },
        }))}
        renderOption={(option) => (
          <Card as="button">
            <Box padding={3}>
              <Text size={[2, 2, 3]}>{option.payload.names.join(' > ')}</Text>
            </Box>
          </Card>
        )}
        renderValue={(value, option) =>
          option?.payload.names.join(' > ') || value
        }
        onChange={handleChange}
      />
    </Stack>
  );
});

function categoryIDNamesMap(
  categories: any,
  map: { [key: string]: string[] },
  parentNames: string[]
): { [key: string]: string[] } {
  if (!categories || categories.length === 0) {
    return map;
  }
  categories.forEach((category: any) => {
    const names = [...parentNames, category.name];
    map[category.id] = names;
    categoryIDNamesMap(category.children, map, names);
  });

  return map;
}
