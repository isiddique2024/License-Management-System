import React from "react";
import { Signal } from "@preact/signals-react";

interface PaginationProps {
  totalPages: number;
  currentPage: Signal<number>;
}

const Pagination: React.FC<PaginationProps> = ({ totalPages, currentPage }) => {
  const handlePageClick = (page: number) => {
    currentPage.value = page;
  };

  const generatePageNumbers = (): (number | string)[] => {
    const pageNumbers: (number | string)[] = [];
    const maxButtons = 5; // max buttons to display
    const currentPageValue = currentPage.value;

    const addRange = (start: number, end: number) => {
      for (let i = start; i <= end; i++) {
        pageNumbers.push(i);
      }
    };

    if (totalPages <= maxButtons) {
      addRange(1, totalPages);
    } else {
      let startPage = Math.max(1, currentPageValue - 2);
      let endPage = Math.min(totalPages, currentPageValue + 2);

      if (currentPageValue <= 3) {
        endPage = 5;
      } else if (currentPageValue + 2 >= totalPages) {
        startPage = totalPages - 4;
      }

      addRange(startPage, endPage);

      if (startPage > 1) {
        pageNumbers.unshift(startPage > 2 ? "..." : 1);
        pageNumbers.unshift(1);
      }

      if (endPage < totalPages) {
        pageNumbers.push(endPage < totalPages - 1 ? "..." : totalPages);
        pageNumbers.push(totalPages);
      }
    }

    return pageNumbers;
  };

  return (
    <div className="flex justify-center items-center mt-4">
      <div className="flex space-x-2">
        <button
          className="btn btn-sm"
          disabled={currentPage.value === 1}
          onClick={() => handlePageClick(currentPage.value - 1)}
        >
          Previous
        </button>
        {generatePageNumbers().map((page) =>
          typeof page === "number" ? (
            <button
              key={page}
              className={`btn btn-sm ${
                page === currentPage.value ? "btn-active" : ""
              }`}
              onClick={() => handlePageClick(page)}
            >
              {page}
            </button>
          ) : (
            <span key={`ellipsis-${page}`} className="btn btn-sm">
              {page}
            </span>
          )
        )}
        <button
          className="btn btn-sm"
          disabled={currentPage.value === totalPages}
          onClick={() => handlePageClick(currentPage.value + 1)}
        >
          Next
        </button>
      </div>
    </div>
  );
};

export default Pagination;
