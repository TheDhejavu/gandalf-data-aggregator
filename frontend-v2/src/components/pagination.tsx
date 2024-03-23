import React from 'react';

interface PaginationProps {
  total: number;
  limit: number;
  currentPage: number;
  onPageChange: (page: number) => void;
}

const Pagination: React.FC<PaginationProps> = ({ total, limit, currentPage, onPageChange }) => {
  const totalPages = Math.ceil(total / limit);
  
  const pageNumbersToShow = 6; 
  
  const getPageNumbers = () => {
    let startPage: number, endPage: number;

    if (totalPages <= pageNumbersToShow) {
      startPage = 1;
      endPage = totalPages;
    } else {
      if (currentPage <= Math.ceil(pageNumbersToShow / 2)) {
        startPage = 1;
        endPage = pageNumbersToShow;
      } else if (currentPage + Math.floor(pageNumbersToShow / 2) >= totalPages) {
        startPage = totalPages - (pageNumbersToShow - 1);
        endPage = totalPages;
      } else {
        startPage = currentPage - Math.floor(pageNumbersToShow / 2);
        endPage = currentPage + Math.floor(pageNumbersToShow / 2);
      }
    }

    const pages = [];
    for (let i = startPage; i <= endPage; i++) {
      pages.push(i);
    }
    return pages;
  };

  const handlePageChange = (page: number) => {
    if (page !== currentPage) {
      onPageChange(page);
    }
  };

  return (
    <nav className="flex justify-center my-4">
      <ul className="relative z-0 inline-flex rounded-md -space-x-px" aria-label="Pagination">
        {currentPage > 1 && (
          <li>
            <button onClick={() => handlePageChange(1)} className="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50">
              First
            </button>
          </li>
        )}
        {getPageNumbers().map((page) => (
          <li key={page}>
            <button onClick={() => handlePageChange(page)} className={`relative inline-flex items-center px-4 py-2 border ${currentPage === page ? 'bg-gray-900 border-gray-900 text-white hover:bg-gray-800' : 'bg-white text-gray-700 hover:bg-gray-50 border-gray-200 '} text-sm font-medium`}>
              {page}
            </button>
          </li>
        ))}
        {currentPage < totalPages && (
          <li>
            <button onClick={() => handlePageChange(totalPages)} className="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50">
              Last
            </button>
          </li>
        )}
      </ul>
    </nav>
  );
};

export default Pagination;
