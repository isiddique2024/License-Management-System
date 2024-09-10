import React from "react";

interface TableHeaderProps {
  headers: { key: string; label: string }[];
  selectAll: boolean;
  onSelectAll: (checked: boolean) => void;
}

const TableHeader: React.FC<TableHeaderProps> = ({
  headers,
  selectAll,
  onSelectAll,
}) => {
  return (
    <thead className="text-white bg-[#2a323c] rounded-t-lg ">
      <tr>
        <th>
          <input
            id="select-all-checkbox"
            type="checkbox"
            className="checkbox border-blue-400 [--chkbg:theme(colors.blue.400)] [--chkfg:white] checked:border-blue-400"
            checked={selectAll}
            onChange={(e) => onSelectAll(e.target.checked)}
          />
        </th>
        {headers.map((header) => (
          <th key={header.key} className="px-2 py-1 sm:px-4 sm:py-2">
            {header.label}
          </th>
        ))}
      </tr>
    </thead>
  );
};

export default TableHeader;
