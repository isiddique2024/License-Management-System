import React from "react";
import { handleBanLicense, handleDeleteLicense } from "../licenses/utils";
import { showToast } from "../misc/Toast";

interface TableRowProps {
  item: any;
  columns: {
    key: string;
    label: string;
    render?: (item: any) => React.ReactNode;
  }[];
  onSelect: (key: string, checked: boolean) => void;
  isSelected: boolean;
}

const TableRow: React.FC<TableRowProps> = ({
  item,
  columns,
  onSelect,
  isSelected,
}) => {
  const handleCellClick = async (event: React.MouseEvent) => {
    event.stopPropagation();
    const target = event.target as HTMLElement;
    const text = target.textContent ?? "";

    const range = document.createRange();
    range.selectNodeContents(target);
    const selection = window.getSelection();
    selection?.removeAllRanges();
    selection?.addRange(range);

    try {
      await navigator.clipboard.writeText(text);
      showToast({ message: "License copied to clipboard", type: "success" });
    } catch (err) {
      console.error("Failed to copy text: ", err);
    }
  };

  return (
    <tr className="bg-[#333b45] hover:bg-[#38404A] border-transparent">
      <th>
        <input
          id={`select-${item.key}`}
          type="checkbox"
          className="checkbox border-blue-400 [--chkbg:theme(colors.blue.400)] [--chkfg:white] checked:border-blue-400"
          checked={isSelected}
          onChange={(e) => onSelect(item.key, e.target.checked)}
        />
      </th>
      {columns.map((column) => {
        const value = item[column.key];
        const shouldBlur =
          ["ip", "hwid"].includes(column.key) && value !== "N/A";
        const isBanned = column.key === "status" && value === "Banned";
        return (
          <td
            key={column.key}
            onClick={column.key === "key" ? handleCellClick : undefined}
            className={`${shouldBlur ? "blurred-text" : ""} ${
              isBanned ? "text-red-500" : ""
            }`}
          >
            {column.render ? column.render(item) : value}
          </td>
        );
      })}
      <td>
        <div className="dropdown">
          <button tabIndex={0} className="btn btn-xs text-center">
            ...
          </button>
          <ul className="dropdown-content menu menu-xs p-2 shadow bg-[#2A323C] rounded-2xl z-20">
            <li>
              <button
                className="text-red-500"
                onClick={async (event) => {
                  event.preventDefault();
                  await handleDeleteLicense(item.key);
                }}
              >
                Delete
              </button>
              <button
                className="text-red-500"
                onClick={async (event) => {
                  event.preventDefault();
                  await handleBanLicense(item.key);
                }}
              >
                Ban
              </button>
            </li>
          </ul>
        </div>
      </td>
    </tr>
  );
};

export default TableRow;
