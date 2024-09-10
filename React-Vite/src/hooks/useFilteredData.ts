import { signal } from "@preact/signals-react";

const searchTerm = signal<string>("");

interface FilterableItem {
  [key: string]: any;
}

const useFilteredData = <T extends FilterableItem>(
  data: T[],
  getValue: (item: T) => string
) => {
  const filteredData = data.filter((item) =>
    getValue(item).toLowerCase().includes(searchTerm.value.toLowerCase())
  );

  return { filteredData, searchTerm };
};

export default useFilteredData;
