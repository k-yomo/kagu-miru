export function isNotEmptyArray(arr: any[] | undefined | null): boolean {
  if (arr === undefined || arr === null) {
    return false;
  }
  return arr.length > 0;
}
