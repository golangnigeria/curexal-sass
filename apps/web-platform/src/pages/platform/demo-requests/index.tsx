import React, { useState } from "react";
import {
  Inbox,
  Search,
  CheckCircle2,
  Clock,
  Calendar,
  Phone,
  Mail,
  Building,
  Edit2,
} from "lucide-react";
import {
  useDemoRequests,
  useUpdateDemoRequest,
} from "@/api/hooks/use-demo-requests";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { formatDate } from "@/lib/utils";
import type { DemoRequest } from "@/api/contracts";
import { toast } from "sonner";

export default function DemoRequestsPage() {
  const { data: requests, isLoading, refetch } = useDemoRequests();

  const [search, setSearch] = useState("");
  const [statusFilter, setStatusFilter] = useState("all");

  const [editingRequest, setEditingRequest] = useState<DemoRequest | null>(null);
  const [newStatus, setNewStatus] = useState("pending");
  const [notes, setNotes] = useState("");

  const updateMutation = useUpdateDemoRequest(editingRequest?.id || "");

  const filteredRequests = (requests || []).filter((req) => {
    const matchesSearch =
      (req.laboratoryName || "").toLowerCase().includes(search.toLowerCase()) ||
      (req.contactName || "").toLowerCase().includes(search.toLowerCase()) ||
      (req.email || "").toLowerCase().includes(search.toLowerCase());
    const matchesStatus = statusFilter === "all" || (req.status || "").toLowerCase() === statusFilter.toLowerCase();
    return matchesSearch && matchesStatus;
  });

  const handleOpenEdit = (req: DemoRequest) => {
    setEditingRequest(req);
    setNewStatus(req.status || "pending");
    setNotes(req.notes || "");
  };

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingRequest) return;
    try {
      await updateMutation.mutateAsync({
        status: newStatus,
        notes,
      });
      toast.success("Demo request updated successfully!");
      setEditingRequest(null);
      refetch();
    } catch (err: any) {
      toast.error(err?.response?.data?.message || "Failed to update demo request");
    }
  };

  const getStatusBadge = (status: string) => {
    switch (status?.toLowerCase()) {
      case "completed":
        return (
          <Badge variant="outline" className="border-emerald-500/30 text-emerald-600 bg-emerald-500/5 text-[10px]">
            <CheckCircle2 className="mr-1 h-3 w-3 inline" /> Completed
          </Badge>
        );
      case "scheduled":
        return (
          <Badge variant="outline" className="border-blue-500/30 text-blue-600 bg-blue-500/5 text-[10px]">
            <Calendar className="mr-1 h-3 w-3 inline" /> Scheduled
          </Badge>
        );
      default:
        return (
          <Badge variant="outline" className="border-amber-500/30 text-amber-600 bg-amber-500/5 text-[10px]">
            <Clock className="mr-1 h-3 w-3 inline" /> Pending
          </Badge>
        );
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold tracking-tight text-foreground">
          Inbound Healthcare Demo Requests
        </h1>
        <p className="text-sm text-muted-foreground mt-1">
          Prospective hospitals, diagnostic centers, and clinics requesting platform demonstrations.
        </p>
      </div>

      {/* Filter Bar */}
      <Card className="card-enterprise p-4">
        <div className="flex flex-col sm:flex-row items-center gap-3">
          <div className="relative flex-1 w-full">
            <Search className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search by facility name, contact person, email..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="pl-9 text-xs h-9"
            />
          </div>

          <Select value={statusFilter} onValueChange={setStatusFilter}>
            <SelectTrigger className="w-[140px] text-xs h-9">
              <SelectValue placeholder="All Status" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Status</SelectItem>
              <SelectItem value="pending">Pending</SelectItem>
              <SelectItem value="scheduled">Scheduled</SelectItem>
              <SelectItem value="completed">Completed</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </Card>

      {/* Requests Table */}
      <Card className="card-enterprise overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="border-b border-border bg-muted/40 font-semibold uppercase tracking-wider text-muted-foreground">
              <tr>
                <th className="py-3 px-4">Facility & Contact</th>
                <th className="py-3 px-4">Contact Channels</th>
                <th className="py-3 px-4">Est. Daily Volume</th>
                <th className="py-3 px-4">Status</th>
                <th className="py-3 px-4">Submitted Date</th>
                <th className="py-3 px-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {isLoading ? (
                <tr>
                  <td colSpan={6} className="py-8 text-center text-muted-foreground">
                    Loading demo requests...
                  </td>
                </tr>
              ) : filteredRequests.length === 0 ? (
                <tr>
                  <td colSpan={6} className="py-8 text-center text-muted-foreground">
                    No demo requests found.
                  </td>
                </tr>
              ) : (
                filteredRequests.map((req) => (
                  <tr key={req.id} className="hover:bg-muted/20">
                    <td className="py-3 px-4">
                      <div className="flex flex-col">
                        <span className="font-semibold text-foreground">{req.laboratoryName}</span>
                        <span className="text-muted-foreground text-[11px]">{req.contactName}</span>
                      </div>
                    </td>
                    <td className="py-3 px-4">
                      <div className="flex flex-col gap-0.5">
                        <span className="flex items-center gap-1 text-muted-foreground text-[11px]">
                          <Mail className="h-3 w-3" /> {req.email}
                        </span>
                        {req.phone && (
                          <span className="flex items-center gap-1 text-muted-foreground text-[11px]">
                            <Phone className="h-3 w-3" /> {req.phone}
                          </span>
                        )}
                      </div>
                    </td>
                    <td className="py-3 px-4 font-mono">{req.dailyVolume || "—"} tests/day</td>
                    <td className="py-3 px-4">{getStatusBadge(req.status)}</td>
                    <td className="py-3 px-4 text-muted-foreground">{formatDate(req.createdAt)}</td>
                    <td className="py-3 px-4 text-right">
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleOpenEdit(req)}
                        className="h-7 text-xs"
                      >
                        <Edit2 className="h-3 w-3" />
                      </Button>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </Card>

      {/* Edit Modal */}
      <Dialog open={!!editingRequest} onOpenChange={(open) => !open && setEditingRequest(null)}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="text-base font-semibold">Update Demo Request</DialogTitle>
            <DialogDescription className="text-xs">
              Manage follow-up status for {editingRequest?.laboratoryName}.
            </DialogDescription>
          </DialogHeader>

          <form onSubmit={handleUpdate} className="space-y-3 py-2 text-xs">
            <div className="space-y-1">
              <Label htmlFor="reqStatus">Processing Status</Label>
              <Select value={newStatus} onValueChange={setNewStatus}>
                <SelectTrigger id="reqStatus" className="text-xs">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="pending">Pending</SelectItem>
                  <SelectItem value="scheduled">Scheduled</SelectItem>
                  <SelectItem value="completed">Completed</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-1">
              <Label htmlFor="reqNotes">Follow-up Notes</Label>
              <Input
                id="reqNotes"
                placeholder="e.g. Scheduled walkthrough for Tuesday 2 PM"
                value={notes}
                onChange={(e) => setNotes(e.target.value)}
                className="text-xs"
              />
            </div>

            <DialogFooter className="pt-3">
              <Button type="button" variant="outline" size="sm" onClick={() => setEditingRequest(null)} className="text-xs">
                Cancel
              </Button>
              <Button type="submit" size="sm" className="bg-primary text-primary-foreground text-xs">
                Save Status
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  );
}
