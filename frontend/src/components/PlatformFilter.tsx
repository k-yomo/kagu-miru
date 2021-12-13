import React, { memo } from 'react';
import { SearchActionType, useSearch } from '@src/contexts/search';
import { ItemSellingPlatform } from '@src/generated/graphql';
import PlatformSelect from '@src/components/PlatformSelect';

export default memo(function RatingFilter() {
  const { searchState, dispatch } = useSearch();

  const onChangePlatforms = (platforms: ItemSellingPlatform[]) => {
    dispatch({
      type: SearchActionType.SET_PLATFORM_FILTER,
      payload: platforms,
    });
  };

  return (
    <PlatformSelect
      platforms={searchState.searchInput.filter.platforms}
      onChangePlatforms={onChangePlatforms}
      htmlIdPrefix="desktop"
    />
  );
});
