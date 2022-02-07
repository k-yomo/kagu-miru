import React, { memo } from 'react';

export default memo(function ({
  value,
}: {
  value: { title?: string; content: string };
}) {
  return (
    <div style={{ padding: 10 }}>
      <div
        style={{
          marginBottom: 5,
          paddingBottom: 5,
          borderBottom: 'solid 1px black',
        }}
      >
        {value.title}
      </div>
      <div>{value.content}</div>
    </div>
  );
});
