import React, { useState } from "react";
import {
  Users,
  Search,
  Shield,
  Building2,
  CheckCircle2,
  XCircle,
  Mail,
  Calendar,
  Eye,
  RefreshCw,
} from "lucide-react";
import { usePlatformUsers } from "@/api/hooks/use-users";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
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
import type { DirectoryUser } from "@/api/contracts";

export default function UsersDirectoryPage() {
  const { data: users, isLoading, refetch } = usePlatformUsers();

  const [search, setSearch] = useState("");
  const [roleFilter, setRoleFilter] = useState("all");
  const [statusFilter, setStatusFilter] = useState("all");
  const [inspectUser, setInspectUser] = useState<DirectoryUser | null>(null);

  const filteredUsers = (users || []).filter((u) => {
    const name = u.name || u.userName || "";
    const email = u.email || u.userEmail || "";
    const role = u.role || u.roleName || "";
    const tenant = u.tenantName || u.organizationName || "";

    const matchesSearch =
      name.toLowerCase().includes(search.toLowerCase()) ||
      email.toLowerCase().includes(search.toLowerCase()) ||
      role.toLowerCase().includes(search.toLowerCase()) ||
      tenant.toLowerCase().includes(search.toLowerCase());

    const matchesRole = roleFilter === "all" || role.toLowerCase() === roleFilter.toLowerCase();
    const matchesStatus =
      statusFilter === "all" ||
      (statusFilter === "active" && u.isActive) ||
      (statusFilter === "inactive" && !u.isActive);

    return matchesSearch && matchesRole && matchesStatus;
  });

  const getRoleBadge = (role: string) => {
    const r = (role || "member").toLowerCase();
    if (r === "super_admin" || r === "platform_admin") {
      return (
        <Badge variant="outline" className="border-primary/40 text-primary bg-primary/10 text-[10px] uppercase font-mono">
          Super Admin
        </Badge>
      );
    }
    if (r === "owner" || r === "admin" || r === "org_admin") {
      return (
        <Badge variant="outline" className="border-blue-500/30 text-blue-600 bg-blue-500/5 text-[10px] capitalize">
          {r.replace("_", " ")}
        </Badge>
      );
    }
    if (r === "pathologist" || r === "radiologist" || r === "doctor" || r === "medical_officer") {
      return (
        <Badge variant="outline" className="border-purple-500/30 text-purple-600 bg-purple-500/5 text-[10px] capitalize">
          {r.replace("_", " ")}
        </Badge>
      );
    }
    return (
      <Badge variant="outline" className="border-border text-muted-foreground text-[10px] capitalize">
        {r.replace("_", " ")}
      </Badge>
    );
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">
            User Directory & Access Profiles
          </h1>
          <p className="text-sm text-muted-foreground mt-1">
            Global index of authenticated personnel, clinicians, laboratory scientists, and tenant administrators across all organizations.
          </p>
        </div>

        <Button
          variant="outline"
          size="sm"
          onClick={() => refetch()}
          className="h-9 gap-2 text-xs self-start sm:self-auto"
        >
          <RefreshCw className={isLoading ? "h-3.5 w-3.5 animate-spin" : "h-3.5 w-3.5"} />
          Refresh Directory
        </Button>
      </div>

      {/* Filter & Search Bar */}
      <Card className="card-enterprise p-4">
        <div className="flex flex-col sm:flex-row items-center gap-3">
          <div className="relative flex-1 w-full">
            <Search className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search users by name, email, role, or organization..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="pl-9 text-xs h-9"
            />
          </div>

          <div className="flex items-center gap-2 w-full sm:w-auto">
            <Select value={roleFilter} onValueChange={setRoleFilter}>
              <SelectTrigger className="w-[140px] text-xs h-9">
                <SelectValue placeholder="All Roles" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Roles</SelectItem>
                <SelectItem value="super_admin">Super Admin</SelectItem>
                <SelectItem value="owner">Owner</SelectItem>
                <SelectItem value="admin">Admin</SelectItem>
                <SelectItem value="pathologist">Pathologist</SelectItem>
                <SelectItem value="radiologist">Radiologist</SelectItem>
                <SelectItem value="medical_officer">Medical Officer</SelectItem>
                <SelectItem value="member">Member / Staff</SelectItem>
              </SelectContent>
            </Select>

            <Select value={statusFilter} onValueChange={setStatusFilter}>
              <SelectTrigger className="w-[130px] text-xs h-9">
                <SelectValue placeholder="All Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Status</SelectItem>
                <SelectItem value="active">Active</SelectItem>
                <SelectItem value="inactive">Inactive</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>
      </Card>

      {/* Users Table */}
      <Card className="card-enterprise overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="border-b border-border bg-muted/40 font-semibold uppercase tracking-wider text-muted-foreground">
              <tr>
                <th className="py-3 px-4">User & Email</th>
                <th className="py-3 px-4">Role Title</th>
                <th className="py-3 px-4">Assigned Organization / Branch</th>
                <th className="py-3 px-4">Status</th>
                <th className="py-3 px-4">Joined Date</th>
                <th className="py-3 px-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {isLoading ? (
                <tr>
                  <td colSpan={6} className="py-8 text-center text-muted-foreground">
                    Querying user directory...
                  </td>
                </tr>
              ) : filteredUsers.length === 0 ? (
                <tr>
                  <td colSpan={6} className="py-8 text-center text-muted-foreground">
                    No users matching query.
                  </td>
                </tr>
              ) : (
                filteredUsers.map((u) => {
                  const userName = u.name || u.userName || "Unnamed User";
                  const userEmail = u.email || u.userEmail || "—";
                  const orgName = u.organizationName || u.tenantName || "Global Platform";
                  const roleName = u.role || u.roleName || "Member";

                  return (
                    <tr key={u.id || u.userId} className="hover:bg-muted/20">
                      <td className="py-3 px-4">
                        <div className="flex flex-col">
                          <span className="font-semibold text-foreground">{userName}</span>
                          <span className="text-muted-foreground text-[11px] font-mono flex items-center gap-1">
                            <Mail className="h-3 w-3" /> {userEmail}
                          </span>
                        </div>
                      </td>
                      <td className="py-3 px-4">{getRoleBadge(roleName)}</td>
                      <td className="py-3 px-4">
                        <div className="flex items-center gap-1.5 font-medium text-foreground">
                          <Building2 className="h-3.5 w-3.5 text-muted-foreground" />
                          <span>{orgName}</span>
                        </div>
                      </td>
                      <td className="py-3 px-4">
                        {u.isActive ? (
                          <Badge variant="outline" className="border-emerald-500/30 text-emerald-600 bg-emerald-500/5 text-[10px]">
                            <CheckCircle2 className="mr-1 h-3 w-3 inline" /> Active
                          </Badge>
                        ) : (
                          <Badge variant="outline" className="border-border text-muted-foreground text-[10px]">
                            <XCircle className="mr-1 h-3 w-3 inline" /> Inactive
                          </Badge>
                        )}
                      </td>
                      <td className="py-3 px-4 text-muted-foreground font-mono">
                        {formatDate(u.joinedAt || u.createdAt)}
                      </td>
                      <td className="py-3 px-4 text-right">
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => setInspectUser(u)}
                          className="h-7 text-xs gap-1"
                        >
                          <Eye className="h-3 w-3" /> Details
                        </Button>
                      </td>
                    </tr>
                  );
                })
              )}
            </tbody>
          </table>
        </div>
      </Card>

      {/* User Details Modal */}
      <Dialog open={!!inspectUser} onOpenChange={(open) => !open && setInspectUser(null)}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="text-base font-semibold flex items-center gap-2">
              <Users className="h-4 w-4 text-primary" />
              User Profile & Membership Details
            </DialogTitle>
            <DialogDescription className="text-xs font-mono">
              User ID: {inspectUser?.userId || inspectUser?.id}
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-3 py-2 text-xs">
            <div className="rounded-lg border border-border bg-secondary/30 p-3 space-y-2">
              <div className="flex justify-between">
                <span className="text-muted-foreground">Full Name:</span>
                <span className="font-semibold text-foreground">{inspectUser?.name || inspectUser?.userName}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Email:</span>
                <span className="font-mono">{inspectUser?.email || inspectUser?.userEmail}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Role:</span>
                <span>{getRoleBadge(inspectUser?.role || inspectUser?.roleName || "Member")}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Assigned Organization:</span>
                <span className="font-medium">{inspectUser?.organizationName || inspectUser?.tenantName || "Global Platform"}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Membership State:</span>
                <span className={inspectUser?.isActive ? "text-emerald-600 font-semibold" : "text-muted-foreground"}>
                  {inspectUser?.isActive ? "Active Account" : "Inactive / Suspended"}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Joined Cluster:</span>
                <span className="font-mono">{formatDate(inspectUser?.joinedAt || inspectUser?.createdAt)}</span>
              </div>
            </div>
          </div>

          <DialogFooter>
            <Button size="sm" onClick={() => setInspectUser(null)} className="text-xs">
              Close Profile
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
