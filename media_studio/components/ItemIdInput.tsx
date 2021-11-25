import React, { useEffect } from 'react'
import { TextInput, Stack, Label } from '@sanity/ui'
import { KAGU_MIRU_URL } from "../config/env"

// TODO: Fix type
export default React.forwardRef<HTMLInputElement>((props: any, ref) => {
    return (
      <Stack space={2}>
        <Label>{props.type.title}</Label>
        <TextInput ref={ref} value={props.value} />
        <a href={`${KAGU_MIRU_URL}/admin/search`} target="_blank">商品を検索する</a>
      </Stack>
    )
  }
)
