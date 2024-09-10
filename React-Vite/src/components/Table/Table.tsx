// Table.tsx
import React from "react";
import { Signal } from "@preact/signals-react";
import TableHeader from "./TableHeader";
import TableRow from "./TableRow";
import Pagination from "./Pagination";
import { CSSTransition, TransitionGroup } from "react-transition-group";
import {
  currentPage,
  ITEMS_PER_PAGE,
  selectedLicenses,
  handleSelectLicense,
  handleSelectAll,
} from "../licenses/utils";
import "./Transition.css";

interface TableProps {
  filteredData: any[];
  searchTerm: Signal<string>;
  columns: {
    key: string;
    label: string;
    render?: (item: any) => React.ReactNode;
  }[];
  toolbar: React.ReactNode;
}

const Table: React.FC<TableProps> = ({
  filteredData,
  searchTerm,
  columns,
  toolbar,
}) => {
  const totalPages = Math.ceil(filteredData.length / ITEMS_PER_PAGE);
  const startIndex = (currentPage.value - 1) * ITEMS_PER_PAGE;
  const currentData = filteredData.slice(
    startIndex,
    startIndex + ITEMS_PER_PAGE
  );

  return (
    <div className="flex flex-col items-center justify-center ">
      {toolbar}
      <div className="flex flex-col w-11/12 mb-4 p-4">
        <input
          type="text"
          value={searchTerm.value}
          onChange={(e) => (searchTerm.value = e.target.value)}
          className="input p-1 rounded-full text-center"
          placeholder="Search"
        />
      </div>
      <div className="flex w-screen flex-col ">
        <div className="divider divider-neutral"></div>
      </div>
      <div className="relative flex-shrink flex-col overflow-x-visible overflow-y-hidden w-screen justify-center items-center px-16">
        <table className="table w-full font-bold ">
          <TableHeader
            headers={columns}
            selectAll={selectedLicenses.value.length === currentData.length}
            onSelectAll={(selected) => handleSelectAll(selected, currentData)}
          />
          <TransitionGroup component="tbody">
            {currentData.map((item) => (
              <CSSTransition key={item.key} timeout={300} classNames="fade">
                <TableRow
                  key={item.key}
                  item={item}
                  columns={columns}
                  onSelect={handleSelectLicense}
                  isSelected={selectedLicenses.value.includes(item.key)}
                />
              </CSSTransition>
            ))}
          </TransitionGroup>
        </table>
      </div>
      <Pagination totalPages={totalPages} currentPage={currentPage} />
    </div>
  );
};

export default Table;
