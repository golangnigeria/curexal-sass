import React, { useState } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Search, ChevronLeft, ChevronRight, Download } from "lucide-react";

export interface Column<T> {
  header: string;
  accessorKey?: keyof T;
  cell?: (item: T) => React.ReactNode;
  className?: string;
}

interface DataTableProps<T> {
  data: T[];
  columns: Column<T>[];
  searchPlaceholder?: string;
  searchKey?: keyof T;
  pageSize?: number;
  onExportCsv?: () => void;
  isLoading?: boolean;
  emptyState?: React.ReactNode;
}

export function DataTable<T extends Record<string, any>>({
  data,
  columns,
  searchPlaceholder = "Search records...",
  searchKey,
  pageSize = 10,
  onExportCsv,
  isLoading,
  emptyState,
}: DataTableProps<T>) {
  const [searchQuery, setSearchQuery] = useState("");
  const [currentPage, setCurrentPage] = useState(1);

  const filteredData = React.useMemo(() => {
    if (!searchQuery.trim()) return data;
    const q = searchQuery.toLowerCase();
    return data.filter((item) => {
      if (searchKey && item[searchKey]) {
        return String(item[searchKey]).toLowerCase().includes(q);
      }
      return Object.values(item).some((val) =>
        String(val).toLowerCase().includes(q)
      );
    });
  }, [data, searchQuery, searchKey]);

  const totalPages = Math.ceil(filteredData.length / pageSize) || 1;
  const paginatedData = filteredData.slice((currentPage - 1) * pageSize, currentPage * pageSize);

  return (
    <div className="space-y-4">
      {/* Table Toolbar */}
      <div className="flex flex-col sm:flex-row items-center justify-between gap-3">
        <div className="w-full sm:max-w-xs relative">
          <Search className="w-3.5 h-3.5 absolute left-3 top-2.5 text-muted-foreground" />
          <Input
            placeholder={searchPlaceholder}
            value={searchQuery}
            onChange={(e) => {
              setSearchQuery(e.target.value);
              setCurrentPage(1);
            }}
            className="text-xs h-8 pl-9 bg-secondary/30"
          />
        </div>

        {onExportCsv && (
          <Button
            size="sm"
            variant="outline"
            onClick={onExportCsv}
            className="text-xs h-8 gap-1.5"
          >
            <Download className="w-3.5 h-3.5" />
            Export CSV
          </Button>
        )}
      </div>

      {/* Table Canvas */}
      <div className="rounded-xl border border-border bg-card overflow-hidden shadow-sm">
        <div className="overflow-x-auto">
          <table className="w-full text-xs text-left">
            <thead className="bg-secondary/40 border-b border-border text-[11px] font-bold text-muted-foreground uppercase tracking-wider">
              <tr>
                {columns.map((col, idx) => (
                  <th key={idx} className={`py-3 px-4 ${col.className || ""}`}>
                    {col.header}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {isLoading ? (
                <tr>
                  <td colSpan={columns.length} className="py-12 text-center text-muted-foreground animate-pulse">
                    Loading records...
                  </td>
                </tr>
              ) : paginatedData.length > 0 ? (
                paginatedData.map((item, rowIdx) => (
                  <tr key={rowIdx} className="hover:bg-secondary/20 transition-colors">
                    {columns.map((col, colIdx) => (
                      <td key={colIdx} className={`py-3 px-4 ${col.className || ""}`}>
                        {col.cell ? col.cell(item) : col.accessorKey ? String(item[col.accessorKey] ?? "") : null}
                      </td>
                    ))}
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan={columns.length} className="py-12 text-center text-muted-foreground">
                    {emptyState || "No records found matching your criteria."}
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination Footer */}
        {filteredData.length > pageSize && (
          <div className="p-3 border-t border-border bg-secondary/20 flex items-center justify-between text-xs text-muted-foreground">
            <span>
              Showing {Math.min(filteredData.length, (currentPage - 1) * pageSize + 1)} to{" "}
              {Math.min(filteredData.length, currentPage * pageSize)} of {filteredData.length} records
            </span>
            <div className="flex items-center gap-1.5">
              <Button
                variant="outline"
                size="icon"
                disabled={currentPage === 1}
                onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
                className="h-7 w-7"
              >
                <ChevronLeft className="w-3.5 h-3.5" />
              </Button>
              <span className="font-mono px-2 text-xs">
                {currentPage} / {totalPages}
              </span>
              <Button
                variant="outline"
                size="icon"
                disabled={currentPage >= totalPages}
                onClick={() => setCurrentPage((p) => Math.min(totalPages, p + 1))}
                className="h-7 w-7"
              >
                <ChevronRight className="w-3.5 h-3.5" />
              </Button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
