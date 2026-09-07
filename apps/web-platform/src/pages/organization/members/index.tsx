import React, { useState } from "react";
import { Link } from "react-router-dom";
import { useBootstrap } from "@/api/hooks/use-bootstrap";
import {
  useOrgMembers,
  useInviteMember,
  useCreateMember,
  useOrgBranches,
  useOrgRoles,
} from "@/api/hooks/use-organization";
import { Card, CardContent } from "@/components/ui/card";
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
import { TableSkeleton } from "@/components/loading";
import { toast } from "sonner";
import {
  Users,
  UserPlus,
  Mail,
  Search,
  Building2,
  Sparkles,
  KeyRound,
  CheckCircle2,
  ShieldAlert,
} from "lucide-react";

export default function OrganizationMembersPage() {
  const { data: bootstrap } = useBootstrap();
  const { data: members, isLoading } = useOrgMembers();
  const { data: branches } = useOrgBranches();
  const { data: roles } = useOrgRoles();

  const inviteMemberMutation = useInviteMember();
  const createMemberMutation = useCreateMember();

  const limits = bootstrap?.limits || { maxBranches: 1, maxMembers: 5, storageGb: 10 };
  const orgPlan = bootstrap?.organization?.subscription || "smart";

  const [searchQuery, setSearchQuery] = useState("");
  const [isInviteOpen, setIsInviteOpen] = useState(false);
  const [isQuickAddOpen, setIsQuickAddOpen] = useState(false);

  // Invite Form State
  const [inviteFullName, setInviteFullName] = useState("");
  const [inviteEmail, setInviteEmail] = useState("");
  const [inviteRole, setInviteRole] = useState("doctor");
  const [inviteBranch, setInviteBranch] = useState("");

  // Quick Add Form State
  const [quickFullName, setQuickFullName] = useState("");
  const [quickEmail, setQuickEmail] = useState("");
  const [quickPassword, setQuickPassword] = useState("password");
  const [quickRole, setQuickRole] = useState("doctor");
  const [quickPrimaryBranch, setQuickPrimaryBranch] = useState("");
  const [quickSelectedBranches, setQuickSelectedBranches] = useState<string[]>([]);

  const membersCount = members?.length || 0;
  const maxMembers = limits.maxMembers || 5;
  const isSeatLimitReached = membersCount >= maxMembers;

  const handleInvite = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!inviteEmail.trim() || !inviteFullName.trim()) {
      toast.error("Please enter both staff member full name and email.");
      return;
    }

    try {
      await inviteMemberMutation.mutateAsync({
        fullName: inviteFullName.trim(),
        email: inviteEmail.trim().toLowerCase(),
        role: inviteRole,
        tenantId: inviteBranch || undefined,
      });

      toast.success("Staff Member Invitation Sent!", {
        description: `An activation link was emailed to ${inviteEmail}.`,
      });

      setIsInviteOpen(false);
      setInviteFullName("");
      setInviteEmail("");
    } catch (err: any) {
      toast.error("Failed to send invitation: " + (err.message || "Network error"));
    }
  };

  const handleQuickAdd = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!quickFullName.trim() || !quickEmail.trim()) {
      toast.error("Please enter both staff member full name and email.");
      return;
    }

    const branchIds = Array.from(
      new Set(
        [quickPrimaryBranch, ...quickSelectedBranches].filter(Boolean)
      )
    );

    try {
      await createMemberMutation.mutateAsync({
        fullName: quickFullName.trim(),
        email: quickEmail.trim().toLowerCase(),
        password: quickPassword || "password",
        role: quickRole,
        roleTitle: quickRole.replace(/_/g, " "),
        facilityBranchId: quickPrimaryBranch || undefined,
        branchIds: branchIds.length > 0 ? branchIds : undefined,
      });

      toast.success("Staff Member Provisioned Successfully!", {
        description: `${quickFullName} can now sign in immediately with email & password.`,
      });

      setIsQuickAddOpen(false);
      setQuickFullName("");
      setQuickEmail("");
      setQuickPassword("password");
      setQuickSelectedBranches([]);
    } catch (err: any) {
      toast.error("Failed to provision staff member: " + (err.message || "Server error"));
    }
  };

  const toggleBranchSelection = (branchId: string) => {
    setQuickSelectedBranches((prev) =>
      prev.includes(branchId) ? prev.filter((id) => id !== branchId) : [...prev, branchId]
    );
  };

  const filteredMembers = (members || []).filter((m: any) => {
    const name = (m.fullName || m.full_name || m.name || "").toLowerCase();
    const memberEmail = (m.email || "").toLowerCase();
    const memberRole = (m.roleTitle || m.role_title || m.role || "").toLowerCase();
    const query = searchQuery.toLowerCase().trim();

    if (!query) return true;
    return name.includes(query) || memberEmail.includes(query) || memberRole.includes(query);
  });

  return (
    <div className="space-y-8 animate-fade-in">
      {/* Header & Controls */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border pb-6">
        <div>
          <div className="flex items-center gap-2.5 mb-1">
            <h1 className="text-2xl font-bold tracking-tight text-foreground">
              Staff Roster & Access Control
            </h1>
            <Badge
              variant="outline"
              className="border-primary/40 text-primary bg-primary/10 text-[10px] font-mono uppercase font-bold"
            >
              {membersCount} of {maxMembers} Seats Allocated
            </Badge>
          </div>
          <p className="text-xs text-muted-foreground">
            Provision clinicians, laboratory scientists, radiologists, nurses, and billing officers across your branch facilities.
          </p>
        </div>

        <div className="flex items-center gap-2.5">
          {/* Direct Quick-Add Staff Dialog */}
          <Dialog open={isQuickAddOpen} onOpenChange={setIsQuickAddOpen}>
            <DialogTrigger asChild>
              <Button
                size="sm"
                className="text-xs h-9 gap-1.5 bg-primary text-primary-foreground shadow"
                disabled={isSeatLimitReached}
              >
                <UserPlus className="w-3.5 h-3.5" />
                Quick Add Staff
              </Button>
            </DialogTrigger>
            <DialogContent className="max-w-md bg-card border-border max-h-[90vh] overflow-y-auto">
              <DialogHeader>
                <DialogTitle className="text-base font-bold flex items-center gap-2">
                  <UserPlus className="w-4 h-4 text-primary" />
                  Quick Add Staff Member
                </DialogTitle>
                <DialogDescription className="text-xs">
                  Instantly provision a staff member account with direct login credentials and branch facility scoping.
                </DialogDescription>
              </DialogHeader>

              <form onSubmit={handleQuickAdd} className="space-y-4 py-2">
                <div className="space-y-1.5">
                  <Label htmlFor="quickFullName" className="text-xs font-medium">Full Name</Label>
                  <Input
                    id="quickFullName"
                    placeholder="e.g. Dr. Sarah Jenkins"
                    value={quickFullName}
                    onChange={(e) => setQuickFullName(e.target.value)}
                    required
                    className="text-xs h-9"
                  />
                </div>

                <div className="space-y-1.5">
                  <Label htmlFor="quickEmail" className="text-xs font-medium">Work Email Address</Label>
                  <Input
                    id="quickEmail"
                    type="email"
                    placeholder="sarah.jenkins@hospital.org"
                    value={quickEmail}
                    onChange={(e) => setQuickEmail(e.target.value)}
                    required
                    className="text-xs h-9"
                  />
                </div>

                <div className="space-y-1.5">
                  <Label htmlFor="quickPassword" className="text-xs font-medium flex items-center justify-between">
                    <span>Initial Password</span>
                    <span className="text-[10px] text-muted-foreground font-normal">Default: password</span>
                  </Label>
                  <div className="relative">
                    <KeyRound className="w-3.5 h-3.5 absolute left-3 top-3 text-muted-foreground" />
                    <Input
                      id="quickPassword"
                      type="text"
                      placeholder="Initial password (default: password)"
                      value={quickPassword}
                      onChange={(e) => setQuickPassword(e.target.value)}
                      className="text-xs h-9 pl-9 font-mono"
                    />
                  </div>
                </div>

                <div className="space-y-1.5">
                  <Label htmlFor="quickRole" className="text-xs font-medium">Clinical / Governance Role</Label>
                  <select
                    id="quickRole"
                    value={quickRole}
                    onChange={(e) => setQuickRole(e.target.value)}
                    className="w-full h-9 rounded-md border border-input bg-background px-3 text-xs focus:outline-none focus:ring-1 focus:ring-primary"
                  >
                    <option value="doctor">Medical Doctor / Clinician</option>
                    <option value="lab_scientist">Medical Laboratory Scientist</option>
                    <option value="radiologist">Consultant Radiologist</option>
                    <option value="pharmacist">Pharmacist</option>
                    <option value="nurse">Nurse / Triage Officer</option>
                    <option value="cashier">Billing / Cashier Officer</option>
                    <option value="receptionist">Front Desk / Receptionist</option>
                    <option value="org_admin">Organization Administrator</option>
                    <option value="pathologist">Chief / Consultant Pathologist</option>
                    {roles?.map((r) => (
                      <option key={r.id} value={r.code || r.name}>{r.name}</option>
                    ))}
                  </select>
                </div>

                <div className="space-y-1.5">
                  <Label htmlFor="quickPrimaryBranch" className="text-xs font-medium">Primary Branch Facility</Label>
                  <select
                    id="quickPrimaryBranch"
                    value={quickPrimaryBranch}
                    onChange={(e) => setQuickPrimaryBranch(e.target.value)}
                    className="w-full h-9 rounded-md border border-input bg-background px-3 text-xs focus:outline-none focus:ring-1 focus:ring-primary"
                  >
                    <option value="">Global HQ (All Branch Facilities)</option>
                    {branches?.map((b) => (
                      <option key={b.id} value={b.id}>
                        {b.name} ({b.code})
                      </option>
                    ))}
                  </select>
                </div>

                {/* Multi-Branch Switching Access */}
                {branches && branches.length > 1 && (
                  <div className="space-y-2 pt-1 border-t border-border">
                    <Label className="text-xs font-medium">
                      Secondary Branch Access (Allows Branch Switching)
                    </Label>
                    <div className="grid grid-cols-1 gap-1.5 max-h-32 overflow-y-auto p-2 rounded-lg bg-secondary/30 border border-border">
                      {branches.map((b) => {
                        const isChecked =
                          quickSelectedBranches.includes(b.id) || quickPrimaryBranch === b.id;
                        return (
                          <label
                            key={b.id}
                            className="flex items-center gap-2 text-xs text-foreground cursor-pointer hover:bg-secondary/50 p-1.5 rounded transition-colors"
                          >
                            <input
                              type="checkbox"
                              checked={isChecked}
                              disabled={quickPrimaryBranch === b.id}
                              onChange={() => toggleBranchSelection(b.id)}
                              className="rounded border-input text-primary focus:ring-primary h-3.5 w-3.5"
                            />
                            <span className="font-medium">{b.name}</span>
                            <span className="text-[10px] text-muted-foreground font-mono">({b.code})</span>
                          </label>
                        );
                      })}
                    </div>
                  </div>
                )}

                <DialogFooter className="pt-3">
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={() => setIsQuickAddOpen(false)}
                    className="text-xs h-9"
                  >
                    Cancel
                  </Button>
                  <Button
                    type="submit"
                    size="sm"
                    disabled={createMemberMutation.isPending}
                    className="text-xs h-9 bg-primary text-primary-foreground gap-1.5"
                  >
                    {createMemberMutation.isPending ? "Provisioning..." : "Create Staff Account"}
                  </Button>
                </DialogFooter>
              </form>
            </DialogContent>
          </Dialog>

          {/* Email Invitation Dialog */}
          <Dialog open={isInviteOpen} onOpenChange={setIsInviteOpen}>
            <DialogTrigger asChild>
              <Button
                variant="outline"
                size="sm"
                className="text-xs h-9 gap-1.5 bg-card border-border"
                disabled={isSeatLimitReached}
              >
                <Mail className="w-3.5 h-3.5 text-muted-foreground" />
                Invite via Email
              </Button>
            </DialogTrigger>
            <DialogContent className="max-w-md bg-card border-border">
              <DialogHeader>
                <DialogTitle className="text-base font-bold flex items-center gap-2">
                  <Mail className="w-4 h-4 text-primary" />
                  Invite Staff Member
                </DialogTitle>
                <DialogDescription className="text-xs">
                  Send a secure email invitation link with 7-day token expiration.
                </DialogDescription>
              </DialogHeader>

              <form onSubmit={handleInvite} className="space-y-4 py-2">
                <div className="space-y-1.5">
                  <Label htmlFor="inviteFullName" className="text-xs font-medium">Full Name</Label>
                  <Input
                    id="inviteFullName"
                    placeholder="e.g. Dr. Amina Yusuf"
                    value={inviteFullName}
                    onChange={(e) => setInviteFullName(e.target.value)}
                    required
                    className="text-xs h-9"
                  />
                </div>

                <div className="space-y-1.5">
                  <Label htmlFor="inviteEmail" className="text-xs font-medium">Email Address</Label>
                  <Input
                    id="inviteEmail"
                    type="email"
                    placeholder="amina.yusuf@hospital.org"
                    value={inviteEmail}
                    onChange={(e) => setInviteEmail(e.target.value)}
                    required
                    className="text-xs h-9"
                  />
                </div>

                <div className="space-y-1.5">
                  <Label htmlFor="inviteRole" className="text-xs font-medium">Access Role</Label>
                  <select
                    id="inviteRole"
                    value={inviteRole}
                    onChange={(e) => setInviteRole(e.target.value)}
                    className="w-full h-9 rounded-md border border-input bg-background px-3 text-xs focus:outline-none focus:ring-1 focus:ring-primary"
                  >
                    <option value="doctor">Medical Doctor / Clinician</option>
                    <option value="lab_scientist">Medical Laboratory Scientist</option>
                    <option value="radiologist">Consultant Radiologist</option>
                    <option value="pharmacist">Pharmacist</option>
                    <option value="nurse">Nurse / Triage Officer</option>
                    <option value="cashier">Billing / Cashier Officer</option>
                    <option value="receptionist">Front Desk / Receptionist</option>
                    <option value="org_admin">Organization Administrator</option>
                    {roles?.map((r) => (
                      <option key={r.id} value={r.code || r.name}>{r.name}</option>
                    ))}
                  </select>
                </div>

                <div className="space-y-1.5">
                  <Label htmlFor="inviteBranch" className="text-xs font-medium">Assigned Branch Facility</Label>
                  <select
                    id="inviteBranch"
                    value={inviteBranch}
                    onChange={(e) => setInviteBranch(e.target.value)}
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

      {/* Seat Limit Alert */}
      {isSeatLimitReached && (
        <div className="p-4 rounded-xl border border-amber-500/30 bg-amber-500/10 text-amber-900 dark:text-amber-200 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <Sparkles className="w-5 h-5 text-amber-500 shrink-0" />
            <div>
              <p className="text-xs font-bold">Staff Seat Limit Reached ({membersCount} / {maxMembers})</p>
              <p className="text-[11px] text-amber-800/80 dark:text-amber-300/80">
                You have reached your staff seat capacity on the {orgPlan} plan. Upgrade to invite or provision additional staff.
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
            <div className="p-6">
              <TableSkeleton rows={5} columns={5} showHeader={false} />
            </div>
          ) : filteredMembers.length > 0 ? (
            <div className="overflow-x-auto">
              <table className="w-full text-xs text-left">
                <thead className="bg-secondary/40 border-b border-border text-[11px] font-bold text-muted-foreground uppercase tracking-wider">
                  <tr>
                    <th className="py-3 px-4">Staff Member</th>
                    <th className="py-3 px-4">Role Title</th>
                    <th className="py-3 px-4">Branch Facilities</th>
                    <th className="py-3 px-4">Status</th>
                    <th className="py-3 px-4">Joined Date</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border">
                  {filteredMembers.map((m: any) => {
                    const memberName = m.fullName || m.full_name || m.name || m.email || "Staff Member";
                    const memberInitials =
                      memberName
                        .split(" ")
                        .filter(Boolean)
                        .map((n: string) => n[0])
                        .join("")
                        .slice(0, 2)
                        .toUpperCase() || "ST";
                    const roleDisplay = (m.roleTitle || m.role_title || m.role || "Member").replace(/_/g, " ");
                    const joinedDate = m.joinedAt || m.joined_at || m.createdAt || m.created_at;

                    const assignedBranches: any[] = m.assignedBranches || [];

                    return (
                      <tr key={m.id || m.userId || m.membershipId || m.email} className="hover:bg-secondary/20 transition-colors">
                        <td className="py-3 px-4">
                          <div className="flex items-center gap-3">
                            <div className="w-8 h-8 rounded-full bg-primary/10 text-primary font-bold flex items-center justify-center text-xs">
                              {memberInitials}
                            </div>
                            <div>
                              <p className="font-semibold text-foreground">{memberName}</p>
                              <p className="text-[11px] text-muted-foreground font-mono">{m.email || ""}</p>
                            </div>
                          </div>
                        </td>
                        <td className="py-3 px-4">
                          <Badge variant="outline" className="text-[10px] font-medium border-border capitalize">
                            {roleDisplay}
                          </Badge>
                        </td>
                        <td className="py-3 px-4">
                          {assignedBranches.length > 0 ? (
                            <div className="flex flex-wrap gap-1">
                              {assignedBranches.map((b: any) => (
                                <span
                                  key={b.id}
                                  className="inline-flex items-center gap-1 text-[10px] font-medium px-2 py-0.5 rounded bg-secondary text-secondary-foreground"
                                >
                                  <Building2 className="w-3 h-3 text-primary" />
                                  {b.name || b.code}
                                </span>
                              ))}
                            </div>
                          ) : (
                            <span className="flex items-center gap-1 text-muted-foreground">
                              <Building2 className="w-3.5 h-3.5" />
                              Global HQ (All Branches)
                            </span>
                          )}
                        </td>
                        <td className="py-3 px-4">
                          <span className="inline-flex items-center gap-1 text-[11px] font-semibold text-emerald-600 dark:text-emerald-400">
                            <CheckCircle2 className="w-3 h-3" />
                            {m.status || (m.isActive !== false ? "Active" : "Inactive")}
                          </span>
                        </td>
                        <td className="py-3 px-4 text-muted-foreground font-mono text-[11px]">
                          {joinedDate ? new Date(joinedDate).toLocaleDateString() : "Active"}
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          ) : (
            <div className="p-12 text-center space-y-3">
              <Users className="w-10 h-10 mx-auto text-muted-foreground opacity-50" />
              <h3 className="text-sm font-semibold text-foreground">No staff members found</h3>
              <p className="text-xs text-muted-foreground">
                Provision or invite your first healthcare team member to begin collaborating.
              </p>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
