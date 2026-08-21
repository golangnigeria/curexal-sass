import React, { useState } from "react";
import { Link } from "react-router-dom";
import { useBootstrap } from "@/api/hooks/use-bootstrap";
import { useOrgMembers, useInviteMember, useOrgBranches, useOrgRoles } from "@/api/hooks/use-organization";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { toast } from "sonner";
import {
  Users,
  UserPlus,
  Search,
  Mail,
  Building2,
  Shield,
  CheckCircle2,
  Sparkles,
  MoreHorizontal,
} from "lucide-react";

export default function OrganizationMembersPage() {
  const { data: bootstrap } = useBootstrap();
  const { data: members, isLoading } = useOrgMembers();
  const { data: branches } = useOrgBranches();
  const { data: roles } = useOrgRoles();
  const inviteMemberMutation = useInviteMember();

  const limits = bootstrap?.limits || { maxBranches: 1, maxMembers: 5, storageGb: 10 };
  const orgPlan = bootstrap?.organization?.subscription || "smart";

  const [searchQuery, setSearchQuery] = useState("");
  const [isInviteOpen, setIsInviteOpen] = useState(false);

  // Form State
  const [fullName, setFullName] = useState("");
  const [email, setEmail] = useState("");
  const [role, setRole] = useState("org_admin");
  const [selectedBranch, setSelectedBranch] = useState("");

  const membersCount = members?.length || 0;
  const maxMembers = limits.maxMembers || 5;
  const isSeatLimitReached = membersCount >= maxMembers;

  const handleInvite = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email.trim() || !fullName.trim()) {
      toast.error("Please enter both staff member full name and email.");
      return;
    }

    try {
      await inviteMemberMutation.mutateAsync({
        fullName: fullName.trim(),
        email: email.trim().toLowerCase(),
        role,
        tenantId: selectedBranch || undefined,
      });

      toast.success("Staff Member Invitation Sent!", {
        description: `An activation link was emailed to ${email}.`,
      });

      setIsInviteOpen(false);
      setFullName("");
      setEmail("");
    } catch (err: any) {
      toast.error("Failed to send invitation: " + (err.message || "Network error"));
    }
  };

  const filteredMembers = (members || []).filter((m) =>
    m.fullName.toLowerCase().includes(searchQuery.toLowerCase()) ||
    m.email.toLowerCase().includes(searchQuery.toLowerCase()) ||
    m.role.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="space-y-8 animate-fade-in">
      {/* Header & Controls */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border pb-6">
        <div>
          <div className="flex items-center gap-2.5 mb-1">
            <h1 className="text-2xl font-bold tracking-tight text-foreground">
              Staff Roster & Access Control
            </h1>
            <Badge variant="outline" className="border-primary/40 text-primary bg-primary/10 text-[10px] font-mono uppercase font-bold">
              {membersCount} of {maxMembers} Seats Allocated
            </Badge>
          </div>
          <p className="text-xs text-muted-foreground">
            Manage clinicians, laboratory scientists, radiologists, nurses, and billing officers across your branch facilities.
          </p>
        </div>

        <div className="flex items-center gap-3">
          <Dialog open={isInviteOpen} onOpenChange={setIsInviteOpen}>
            <DialogTrigger asChild>
              <Button
                size="sm"
                className="text-xs h-9 gap-1.5 bg-primary text-primary-foreground shadow"
                disabled={isSeatLimitReached}
              >
                <UserPlus className="w-3.5 h-3.5" />
                Invite Staff Member
              </Button>
            </DialogTrigger>
            <DialogContent className="max-w-md bg-card border-border">
              <DialogHeader>
                <DialogTitle className="text-base font-bold">Invite Staff Member</DialogTitle>
                <DialogDescription className="text-xs">
                  Send an official email invitation with assigned organizational role and branch facility.
                </DialogDescription>
              </DialogHeader>

              <form onSubmit={handleInvite} className="space-y-4 py-2">
                <div className="space-y-1.5">
                  <Label htmlFor="fullName" className="text-xs font-medium">Full Name</Label>
                  <Input
                    id="fullName"
                    placeholder="e.g. Dr. Amina Yusuf"
                    value={fullName}
                    onChange={(e) => setFullName(e.target.value)}
                    required
                    className="text-xs h-9"
                  />
                </div>

                <div className="space-y-1.5">
                  <Label htmlFor="email" className="text-xs font-medium">Email Address</Label>
                  <Input
                    id="email"
                    type="email"
                    placeholder="amina.yusuf@hospital.org"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                    className="text-xs h-9"
                  />
                </div>

                <div className="space-y-1.5">
                  <Label htmlFor="role" className="text-xs font-medium">Access Role</Label>
                  <select
                    id="role"
                    value={role}
                    onChange={(e) => setRole(e.target.value)}
                    className="w-full h-9 rounded-md border border-input bg-background px-3 text-xs focus:outline-none focus:ring-1 focus:ring-primary"
                  >
                    <option value="org_admin">Organization Administrator</option>
                    <option value="pathologist">Chief / Consultant Pathologist</option>
                    <option value="radiologist">Consultant Radiologist</option>
                    <option value="doctor">Medical Doctor / Clinician</option>
                    <option value="lab_scientist">Medical Laboratory Scientist</option>
                    <option value="pharmacist">Pharmacist</option>
                    <option value="nurse">Nurse / Triage Officer</option>
                    <option value="cashier">Billing / Cashier Officer</option>
                  </select>
                </div>

                <div className="space-y-1.5">
                  <Label htmlFor="branch" className="text-xs font-medium">Assigned Branch Facility</Label>
                  <select
                    id="branch"
                    value={selectedBranch}
                    onChange={(e) => setSelectedBranch(e.target.value)}
                    className="w-full h-9 rounded-md border border-input bg-background px-3 text-xs focus:outline-none focus:ring-1 focus:ring-primary"
                  >
                    <option value="">All Branch Facilities (Global HQ)</option>
                    {branches?.map((b) => (
                      <option key={b.id} value={b.id}>{b.name} ({b.code})</option>
                    ))}
                  </select>
                </div>

                <DialogFooter className="pt-3">
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={() => setIsInviteOpen(false)}
                    className="text-xs h-9"
                  >
                    Cancel
                  </Button>
                  <Button
                    type="submit"
                    size="sm"
                    disabled={inviteMemberMutation.isPending}
                    className="text-xs h-9 bg-primary text-primary-foreground gap-1.5"
                  >
                    {inviteMemberMutation.isPending ? "Sending..." : "Send Invitation"}
                  </Button>
                </DialogFooter>
              </form>
            </DialogContent>
          </Dialog>
        </div>
      </div>

      {/* Limit Alert */}
      {isSeatLimitReached && (
        <div className="p-4 rounded-xl border border-amber-500/30 bg-amber-500/10 text-amber-900 dark:text-amber-200 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <Sparkles className="w-5 h-5 text-amber-500 shrink-0" />
            <div>
              <p className="text-xs font-bold">Staff Seat Limit Reached ({membersCount} / {maxMembers})</p>
              <p className="text-[11px] text-amber-800/80 dark:text-amber-300/80">
                You have reached your staff seat capacity on the {orgPlan} plan. Upgrade to invite additional staff.
              </p>
            </div>
          </div>
          <Button asChild size="sm" className="text-xs h-8 bg-amber-600 hover:bg-amber-700 text-white">
            <Link to="/organization/billing">Add Seats</Link>
          </Button>
        </div>
      )}

      {/* Filter */}
      <div className="max-w-sm relative">
        <Search className="w-3.5 h-3.5 absolute left-3 top-2.5 text-muted-foreground" />
        <Input
          placeholder="Filter staff by name, email, or role..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="text-xs h-8 pl-9 bg-secondary/30"
        />
      </div>

      {/* Members Table */}
      <Card className="border-border shadow-sm overflow-hidden">
        <CardContent className="p-0">
          {isLoading ? (
            <div className="py-16 text-center text-xs text-muted-foreground animate-pulse">
              Loading staff roster...
            </div>
          ) : filteredMembers.length > 0 ? (
            <div className="overflow-x-auto">
              <table className="w-full text-xs text-left">
                <thead className="bg-secondary/40 border-b border-border text-[11px] font-bold text-muted-foreground uppercase tracking-wider">
                  <tr>
                    <th className="py-3 px-4">Staff Member</th>
                    <th className="py-3 px-4">Role Title</th>
                    <th className="py-3 px-4">Branch Access</th>
                    <th className="py-3 px-4">Status</th>
                    <th className="py-3 px-4">Joined Date</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border">
                  {filteredMembers.map((m) => (
                    <tr key={m.id} className="hover:bg-secondary/20 transition-colors">
                      <td className="py-3 px-4">
                        <div className="flex items-center gap-3">
                          <div className="w-8 h-8 rounded-full bg-primary/10 text-primary font-bold flex items-center justify-center text-xs">
                            {m.fullName.split(" ").map((n) => n[0]).join("").slice(0, 2).toUpperCase()}
                          </div>
                          <div>
                            <p className="font-semibold text-foreground">{m.fullName}</p>
                            <p className="text-[11px] text-muted-foreground font-mono">{m.email}</p>
                          </div>
                        </div>
                      </td>
                      <td className="py-3 px-4">
                        <Badge variant="outline" className="text-[10px] font-medium border-border capitalize">
                          {m.roleTitle || m.role.replace("_", " ")}
                        </Badge>
                      </td>
                      <td className="py-3 px-4">
                        <span className="flex items-center gap-1 text-muted-foreground">
                          <Building2 className="w-3.5 h-3.5" />
                          {m.tenantName || "Global HQ (All Branches)"}
                        </span>
                      </td>
                      <td className="py-3 px-4">
                        <span className="inline-flex items-center gap-1 text-[11px] font-semibold text-emerald-600 dark:text-emerald-400">
                          <span className="w-1.5 h-1.5 rounded-full bg-emerald-500" />
                          Active
                        </span>
                      </td>
                      <td className="py-3 px-4 text-muted-foreground font-mono text-[11px]">
                        {new Date(m.joinedAt).toLocaleDateString()}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <div className="p-12 text-center space-y-3">
              <Users className="w-10 h-10 mx-auto text-muted-foreground opacity-50" />
              <h3 className="text-sm font-semibold text-foreground">No staff members found</h3>
              <p className="text-xs text-muted-foreground">
                Invite your first healthcare team member to begin collaborating.
              </p>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
