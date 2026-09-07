import React, { useState, useMemo } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  Search,
  ChevronLeft,
  ChevronRight,
  Download,
  ArrowUpDown,
  SlidersHorizontal,
  CheckSquare,
  Square,
  Sparkles,
} from "lucide-react";

export interface ColumnDef<T> {
  id?: string;
  header: string | React.ReactNode;
  accessorKey?: keyof T;
  cell?: (item: T, index: number) => React.ReactNode;
  sortable?: boolean;
  className?: string;
}

interface DataTableProps<T> {
  data: T[];
  columns: ColumnDef<T>[];
  searchPlaceholder?: string;
  searchKey?: keyof T;
  pageSize?: number;
  onExportCsv?: () => void;
  isLoading?: boolean;
  emptyState?: React.ReactNode;
  onRowClick?: (item: T) => void;
  enableRowSelection?: boolean;
  onSelectionChange?: (selectedItems: T[]) => void;
  bulkActions?: React.ReactNode;
  statusFilterKey?: keyof T;
  statusOptions?: Array<{ label: string; value: string }>;
}

export function DataTable<T extends Record<string, any>>({
  data,
  columns,
  searchPlaceholder = "Search worklist records...",
  searchKey,
  pageSize = 10,
  onExportCsv,
  isLoading,
  emptyState,
  onRowClick,
  enableRowSelection = false,
  onSelectionChange,
  bulkActions,
  statusFilterKey,
  statusOptions,
}: DataTableProps<T>) {
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedStatus, setSelectedStatus] = useState<string>("ALL");
  const [sortKey, setSortKey] = useState<keyof T | null>(null);
  const [sortDirection, setSortDirection] = useState<"asc" | "desc">("asc");
  const [currentPage, setCurrentPage] = useState(1);
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());

  const handleSort = (key?: keyof T) => {
    if (!key) return;
    if (sortKey === key) {
      setSortDirection((prev) => (prev === "asc" ? "desc" : "asc"));
    } else {
      setSortKey(key);
      setSortDirection("asc");
    }
  };

  const filteredData = useMemo(() => {
    let result = [...data];

    // Status filter
    if (statusFilterKey && selectedStatus !== "ALL") {
      result = result.filter((item) => String(item[statusFilterKey]).toLowerCase() === selectedStatus.toLowerCase());
    }

    // Search query
    if (searchQuery.trim()) {
      const q = searchQuery.toLowerCase();
      result = result.filter((item) => {
        if (searchKey && item[searchKey]) {
          return String(item[searchKey]).toLowerCase().includes(q);
        }
        return Object.values(item).some((val) =>
          String(val).toLowerCase().includes(q)
        );
      });
    }

    // Sorting
    if (sortKey) {
      result.sort((a, b) => {
        const valA = a[sortKey];
        const valB = b[sortKey];
        if (valA === valB) return 0;
        if (valA == null) return 1;
        if (valB == null) return -1;
        const res = valA < valB ? -1 : 1;
        return sortDirection === "asc" ? res : -res;
      });
    }

    return result;
  }, [data, searchQuery, searchKey, selectedStatus, statusFilterKey, sortKey, sortDirection]);

  const totalPages = Math.ceil(filteredData.length / pageSize) || 1;
  const paginatedData = filteredData.slice((currentPage - 1) * pageSize, currentPage * pageSize);

  const toggleSelectAll = () => {
    if (selectedIds.size === paginatedData.length) {
      setSelectedIds(new Set());
      if (onSelectionChange) onSelectionChange([]);
    } else {
      const newSet = new Set(paginatedData.map((item) => String(item.id || item.code || item.mrn)));
      setSelectedIds(newSet);
      if (onSelectionChange) onSelectionChange(paginatedData);
    }
  };

  const toggleSelectRow = (item: T, e: React.MouseEvent) => {
    e.stopPropagation();
    const id = String(item.id || item.code || item.mrn);
    const newSet = new Set(selectedIds);
    if (newSet.has(id)) {
      newSet.delete(id);
    } else {
      newSet.add(id);
    }
    setSelectedIds(newSet);
    if (onSelectionChange) {
      const selected = data.filter((d) => newSet.has(String(d.id || d.code || d.mrn)));
      onSelectionChange(selected);
    }
  };

  return (
    <div className="space-y-3.5">
      {/* Table Toolbar & Status Tabs */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-3">
        <div className="flex flex-wrap items-center gap-2 w-full md:w-auto">
          {/* Search Box */}
          <div className="relative w-full sm:w-72">
            <Search className="w-3.5 h-3.5 absolute left-3 top-2.5 text-muted-foreground" />
            <Input
              placeholder={searchPlaceholder}
              value={searchQuery}
              onChange={(e) => {
                setSearchQuery(e.target.value);
                setCurrentPage(1);
              }}
              className="text-xs h-8 pl-9 bg-secondary/30 border-border"
            />
          </div>

          {/* Status Quick Filter Options */}
          {statusOptions && (
            <div className="flex items-center gap-1 overflow-x-auto pb-1 sm:pb-0">
              <button
                type="button"
                onClick={() => {
                  setSelectedStatus("ALL");
                  setCurrentPage(1);
                }}
                className={`px-2.5 py-1 rounded-md text-[11px] font-medium transition-colors cursor-pointer ${
                  selectedStatus === "ALL"
                    ? "bg-primary text-primary-foreground font-bold shadow-xs"
                    : "bg-secondary text-muted-foreground hover:text-foreground"
                }`}
              >
                All
              </button>
              {statusOptions.map((opt) => (
                <button
                  key={opt.value}
                  type="button"
                  onClick={() => {
                    setSelectedStatus(opt.value);
                    setCurrentPage(1);
                  }}
                  className={`px-2.5 py-1 rounded-md text-[11px] font-medium transition-colors cursor-pointer ${
                    selectedStatus === opt.value
                      ? "bg-primary text-primary-foreground font-bold shadow-xs"
                      : "bg-secondary text-muted-foreground hover:text-foreground"
                  }`}
                >
                  {opt.label}
                </button>
              ))}
            </div>
          )}
        </div>

        {/* Right Toolbar Actions */}
        <div className="flex items-center gap-2 self-end md:self-auto">
          {selectedIds.size > 0 && bulkActions}

          {onExportCsv && (
            <Button
              size="sm"
              variant="outline"
              onClick={onExportCsv}
              className="text-xs h-8 gap-1.5"
            >
              <Download className="w-3.5 h-3.5" />
              <span>Export CSV</span>
            </Button>
          )}
        </div>
      </div>

      {/* Table Canvas */}
      <div className="rounded-xl border border-border bg-card overflow-hidden shadow-sm">
        <div className="overflow-x-auto">
          <table className="w-full text-xs text-left border-collapse">
            <thead className="bg-secondary/40 border-b border-border text-[11px] font-bold text-muted-foreground uppercase tracking-wider">
              <tr>
                {enableRowSelection && (
                  <th className="py-3 px-3 w-8">
                    <button
                      type="button"
                      onClick={toggleSelectAll}
                      className="text-muted-foreground hover:text-foreground cursor-pointer"
                    >
                      {selectedIds.size === paginatedData.length && paginatedData.length > 0 ? (
                        <CheckSquare className="w-3.5 h-3.5 text-primary" />
                      ) : (
                        <Square className="w-3.5 h-3.5" />
                      )}
                    </button>
                  </th>
                )}
                {columns.map((col, idx) => (
                  <th
                    key={idx}
                    onClick={() => col.sortable && handleSort(col.accessorKey)}
                    className={`py-3 px-4 ${col.sortable ? "cursor-pointer select-none hover:text-foreground" : ""} ${col.className || ""}`}
                  >
                    <div className="flex items-center gap-1">
                      <span>{col.header}</span>
                      {col.sortable && <ArrowUpDown className="w-3 h-3 opacity-60" />}
                    </div>
                  </th>
                ))}
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {isLoading ? (
                Array.from({ length: 5 }).map((_, rIdx) => (
                  <tr key={rIdx} className="animate-pulse">
                    {enableRowSelection && <td className="py-3.5 px-3"><div className="w-3.5 h-3.5 bg-secondary rounded" /></td>}
                    {columns.map((_, cIdx) => (
                      <td key={cIdx} className="py-3.5 px-4">
                        <div className="h-3.5 bg-secondary/80 rounded w-3/4" />
                      </td>
                    ))}
                  </tr>
                ))
              ) : paginatedData.length > 0 ? (
                paginatedData.map((item, rowIdx) => {
                  const isSelected = selectedIds.has(String(item.id || item.code || item.mrn));
                  return (
                    <tr
                      key={rowIdx}
                      onClick={() => onRowClick && onRowClick(item)}
                      className={`transition-colors ${
                        onRowClick ? "cursor-pointer hover:bg-secondary/30" : "hover:bg-secondary/15"
                      } ${isSelected ? "bg-primary/5" : ""}`}
                    >
                      {enableRowSelection && (
                        <td className="py-3 px-3">
                          <button
                            type="button"
                            onClick={(e) => toggleSelectRow(item, e)}
                            className="text-muted-foreground hover:text-foreground cursor-pointer"
                          >
                            {isSelected ? (
                              <CheckSquare className="w-3.5 h-3.5 text-primary" />
                            ) : (
                              <Square className="w-3.5 h-3.5" />
                            )}
                          </button>
                        </td>
                      )}
                      {columns.map((col, colIdx) => (
                        <td key={colIdx} className={`py-3 px-4 ${col.className || ""}`}>
                          {col.cell ? col.cell(item, rowIdx) : col.accessorKey ? String(item[col.accessorKey] ?? "") : null}
                        </td>
                      ))}
                    </tr>
                  );
                })
              ) : (
                <tr>
                  <td colSpan={columns.length + (enableRowSelection ? 1 : 0)} className="py-12 text-center text-muted-foreground">
                    {emptyState || "No records matching your search or filters."}
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
