import React from 'react';

// Use PageLoading if you have some data / steps you need to complete before first rendering.
// Loading per component works for most of the case.
export default function Loading() {
  return (
    <div className="p-4">
      <div className="mb-12 w-12 h-12 rounded-full border-4 border-t-4 border-gray-200 ease-linear animate-spin page-loader" />
      <h2 className="text-xl font-semibold text-center text-white">
        Loading...
      </h2>
    </div>
  );
}
