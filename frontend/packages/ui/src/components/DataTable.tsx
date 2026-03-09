import React from 'react';

interface DataTableProps {
  children?: React.ReactNode;
  className?: string;
}

export function DataTable({ children, className = '' }: DataTableProps) {
  return (
    <div className={`overflow-x-auto ${className}`}>
      <p className="text-gray-500 text-sm p-4">DataTable coming soon</p>
      {children}
    </div>
  );
}
