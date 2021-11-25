import React, { memo } from "react"

export default memo(function ({ value }: { value: { html: string } }) {
  return (
    <div dangerouslySetInnerHTML={{ __html: value.html }}/>
  )
})
