import React, { memo } from 'react';
import { ItemColor } from '@src/generated/graphql';
import { SearchActionType, useSearch } from '@src/contexts/search';
import ColorSelect from '@src/components/ColorSelect';

export default memo(function ColorFilter() {
  const { searchState, dispatch } = useSearch();

  const onChangeColors = (colors: ItemColor[]) => {
    dispatch({
      type: SearchActionType.SET_COLOR_FILTER,
      payload: colors,
    });
  };

  return (
    <ColorSelect
      colors={searchState.searchInput.filter.colors}
      onChangeColors={onChangeColors}
    />
  );
});
