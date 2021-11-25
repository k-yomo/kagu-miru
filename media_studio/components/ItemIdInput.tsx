import React, { forwardRef } from 'react'
import { TextInput, Stack, Label } from '@sanity/ui'
import PatchEvent, {set, unset} from '@sanity/form-builder/PatchEvent'
import { KAGU_MIRU_URL } from "../config/env"

// TODO: Fix type
export default forwardRef<HTMLInputElement>((props: any, ref) => {
  const handleChange = React.useCallback(
    // useCallback will help with performance
    (event) => {
      const inputValue = event.currentTarget.value // get current value
      // if the value exists, set the data, if not, unset the data
      props.onChange(PatchEvent.from(inputValue ? set(inputValue) : unset()))
    },
    [props.onChange]
  )
    return (
      <Stack space={2}>
        <Label>{props.type.title}</Label>
        <TextInput
          ref={ref}
          value={props.value}
          onChange={handleChange}
        />
        <a href={`${KAGU_MIRU_URL}/admin/search`} target="_blank">商品を検索する</a>
      </Stack>
    )
  }
)
