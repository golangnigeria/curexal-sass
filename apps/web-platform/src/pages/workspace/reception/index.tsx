import React, { useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import {
  Users,
  UserPlus,
  Search,
  Filter,
  Phone,
  Calendar,
  CreditCard,
  ShieldCheck,
  Send,
  Stethoscope,
  MoreVertical,
  Activity,
  ArrowRight,
  Sparkles,
  CheckCircle,
} from "lucide-react";
import { usePatients, useSendPortalOTP } from "../../../api/hooks/use-patients";
import { PatientIntakeModal } from "../../../components/patients/patient-intake-modal";
import { PatientDrawer } from "@/features/patients/patient-drawer";
import type { CanonicalPatient } from "../../../api/contracts";

export const ReceptionWorkspacePage: React.FC = () => {
  const { branchSlug } = useParams<{ branchSlug: string }>();
  const navigate = useNavigate();

  const [searchQuery, setSearchQuery] = useState("");
  const [statusFilter, setStatusFilter] = useState("");
  const [isIntakeOpen, setIsIntakeOpen] = useState(false);
  const [otpSentNotice, setOtpSentNotice] = useState<string | null>(null);
  const [selectedPatient, setSelectedPatient] = useState<CanonicalPatient | null>(null);

  const { data: patientList, isLoading, refetch } = usePatients({
    query: searchQuery || undefined,
    status: statusFilter || undefined,
    limit: 50,
  });

  const sendPortalOTP = useSendPortalOTP();

  const handleSendOTP = (patient: CanonicalPatient) => {
    const primaryPhone =
      patient.contacts?.find((c) => c.system === "PHONE")?.value || patient.mrn;
    sendPortalOTP.mutate(
      { identifier: primaryPhone },
      {
        onSuccess: (data) => {
          setOtpSentNotice(
            `Portal login OTP sent to ${patient.firstName} (${primaryPhone})`
          );
          setTimeout(() => setOtpSentNotice(null), 5000);
        },
      }
    );
  };

  const totalPatients = patientList?.total || 0;

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 p-6 space-y-6">
      {/* Header Banner */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 p-6 rounded-2xl bg-gradient-to-r from-slate-900 via-slate-900 to-cyan-950/40 border border-slate-800 shadow-xl">
        <div className="flex items-center gap-4">
          <div className="p-3.5 rounded-2xl bg-cyan-500/10 text-cyan-400 border border-cyan-500/20 shadow-inner">
            <Users className="w-7 h-7" />
          </div>
          <div>
            <div className="flex items-center gap-2">
              <h1 className="text-2xl font-bold tracking-tight text-white">
                Patient Access & Intake Desk
              </h1>
              <span className="px-2.5 py-0.5 rounded-full text-[11px] font-semibold bg-cyan-500/10 text-cyan-400 border border-cyan-500/20 uppercase tracking-wider">
                MPI Active
              </span>
            </div>
            <p className="text-xs text-slate-400 mt-1">
              Canonical patient directory, duplicate resolution, and portal onboarding for branch{" "}
              <span className="text-cyan-300 font-mono">/{branchSlug}</span>
            </p>
          </div>
        </div>

        <button
          onClick={() => setIsIntakeOpen(true)}
          className="flex items-center justify-center gap-2 px-5 py-2.5 rounded-xl bg-gradient-to-r from-cyan-500 to-blue-600 hover:from-cyan-400 hover:to-blue-500 text-white text-sm font-semibold shadow-lg shadow-cyan-500/20 transition-all hover:scale-[1.02] active:scale-[0.98]"
        >
          <UserPlus className="w-4 h-4" />
          <span>New Patient Intake</span>
        </button>
      </div>

      {/* OTP Dispatch Feedback Notice */}
      {otpSentNotice && (
        <div className="p-4 rounded-xl bg-emerald-500/10 border border-emerald-500/30 text-emerald-300 text-xs flex items-center justify-between animate-in slide-in-from-top-2">
          <div className="flex items-center gap-2">
            <CheckCircle className="w-4 h-4 text-emerald-400" />
            <span>{otpSentNotice}</span>
          </div>
        </div>
      )}

      {/* Metrics Row */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div className="p-4 rounded-xl bg-slate-900 border border-slate-800 flex items-center justify-between">
          <div>
            <div className="text-xs text-slate-400 font-medium">Canonical Patients</div>
            <div className="text-2xl font-bold text-white mt-1">{totalPatients}</div>
          </div>
          <div className="p-2.5 rounded-lg bg-slate-800 text-cyan-400">
            <Users className="w-5 h-5" />
          </div>
        </div>

        <div className="p-4 rounded-xl bg-slate-900 border border-slate-800 flex items-center justify-between">
          <div>
            <div className="text-xs text-slate-400 font-medium">MPI Duplicates Prevented</div>
            <div className="text-2xl font-bold text-emerald-400 mt-1">100%</div>
          </div>
          <div className="p-2.5 rounded-lg bg-slate-800 text-emerald-400">
            <ShieldCheck className="w-5 h-5" />
          </div>
        </div>

        <div className="p-4 rounded-xl bg-slate-900 border border-slate-800 flex items-center justify-between">
          <div>
            <div className="text-xs text-slate-400 font-medium">Active Care Orchestrations</div>
            <div className="text-2xl font-bold text-blue-400 mt-1">Ready</div>
          </div>
          <div className="p-2.5 rounded-lg bg-slate-800 text-blue-400">
            <Activity className="w-5 h-5" />
          </div>
        </div>
      </div>

      {/* Search & Filter Bar */}
      <div className="flex flex-col sm:flex-row items-center gap-3 p-3 bg-slate-900 border border-slate-800 rounded-xl">
        <div className="relative flex-1 w-full">
          <Search className="w-4 h-4 text-slate-400 absolute left-3 top-2.5" />
          <input
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            placeholder="Search by Patient Name, MRN, Phone Number, or NIN..."
            className="w-full pl-9 pr-3 py-1.5 bg-slate-950 border border-slate-800 rounded-lg text-sm text-white placeholder-slate-500 focus:outline-none focus:border-cyan-500"
          />
        </div>

        <div className="flex items-center gap-2 w-full sm:w-auto">
          <select
            value={statusFilter}
            onChange={(e) => setStatusFilter(e.target.value)}
            className="px-3 py-1.5 bg-slate-950 border border-slate-800 rounded-lg text-xs text-white focus:outline-none focus:border-cyan-500"
          >
            <option value="">All Statuses</option>
            <option value="REGISTERED">Registered</option>
            <option value="ACTIVE">Active</option>
            <option value="IDENTIFIED">Identified</option>
          </select>
        </div>
      </div>

      {/* Patient Directory Table / Grid */}
      <div className="bg-slate-900 border border-slate-800 rounded-2xl overflow-hidden shadow-xl">
        <div className="px-6 py-4 border-b border-slate-800 flex items-center justify-between">
          <h2 className="text-sm font-semibold text-white">Registered Patient Directory</h2>
          <span className="text-xs text-slate-400">{patientList?.items?.length || 0} visible records</span>
        </div>

        {isLoading ? (
          <div className="p-12 text-center text-slate-400 text-xs animate-pulse">
            Loading canonical patient records...
          </div>
        ) : !patientList?.items || patientList.items.length === 0 ? (
          <div className="p-12 text-center space-y-3">
            <div className="w-12 h-12 rounded-full bg-slate-800 text-slate-500 flex items-center justify-center mx-auto">
              <Users className="w-6 h-6" />
            </div>
            <div className="text-sm font-medium text-white">No patients found</div>
            <p className="text-xs text-slate-400 max-w-sm mx-auto">
              No matching records found. Use the "New Patient Intake" button to register a patient.
            </p>
          </div>
        ) : (
          <div className="divide-y divide-slate-800">
            {patientList.items.map((patient) => {
              const primaryPhone =
                patient.contacts?.find((c) => c.system === "PHONE")?.value || "N/A";
              return (
                <div
                  key={patient.id}
                  className="p-5 hover:bg-slate-800/40 transition-colors flex flex-col md:flex-row md:items-center justify-between gap-4"
                >
                  <div className="flex items-start gap-4">
                    <div className="w-10 h-10 rounded-xl bg-gradient-to-tr from-cyan-500/20 to-blue-500/20 text-cyan-300 font-bold flex items-center justify-center border border-cyan-500/30 text-sm">
                      {patient.firstName[0]}
                      {patient.lastName[0]}
                    </div>
                    <div className="space-y-1">
                      <div className="flex items-center gap-2">
                        <span className="font-semibold text-white text-sm">
                          {patient.firstName} {patient.middleName ? `${patient.middleName} ` : ""}
                          {patient.lastName}
                        </span>
                        <span className="px-2 py-0.5 rounded-md bg-slate-800 text-cyan-400 font-mono text-[11px] font-semibold border border-slate-700">
                          {patient.mrn}
                        </span>
                        <span className="px-2 py-0.5 rounded-md bg-emerald-500/10 text-emerald-400 text-[10px] font-medium border border-emerald-500/20">
                          {patient.status}
                        </span>
                      </div>
                      <div className="flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-slate-400">
                        <span className="flex items-center gap-1">
                          <Calendar className="w-3.5 h-3.5" />
                          <span>DOB: {new Date(patient.dateOfBirth).toLocaleDateString()}</span>
                        </span>
                        <span className="flex items-center gap-1">
                          <Phone className="w-3.5 h-3.5" />
                          <span>{primaryPhone}</span>
                        </span>
                        {patient.bloodGroup && (
                          <span className="text-rose-300">
                            Blood: <strong>{patient.bloodGroup}</strong>
                          </span>
                        )}
                        {patient.genotype && (
                          <span className="text-amber-300">
                            Genotype: <strong>{patient.genotype}</strong>
                          </span>
                        )}
                      </div>
                    </div>
                  </div>

                  {/* Actions */}
                  <div className="flex items-center gap-2 self-end md:self-center">
                    <button
                      onClick={() => setSelectedPatient(patient)}
                      title="Inspect Patient 360 Profile"
                      className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white text-xs font-medium border border-slate-700 transition-colors"
                    >
                      <Users className="w-3.5 h-3.5 text-cyan-400" />
                      <span>Patient 360</span>
                    </button>
                    <button
                      onClick={() => handleSendOTP(patient)}
                      title="Send Portal Passwordless OTP"
                      className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white text-xs font-medium border border-slate-700 transition-colors"
                    >
                      <Send className="w-3.5 h-3.5 text-cyan-400" />
                      <span>Send Portal OTP</span>
                    </button>
                    <button
                      onClick={() => navigate(`/${branchSlug}/clinical?patientId=${patient.id}`)}
                      className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-cyan-500/10 hover:bg-cyan-500/20 text-cyan-300 text-xs font-medium border border-cyan-500/30 transition-colors"
                    >
                      <Stethoscope className="w-3.5 h-3.5" />
                      <span>Start Encounter</span>
                    </button>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>

      {/* Patient Intake Modal */}
      <PatientIntakeModal
        isOpen={isIntakeOpen}
        onClose={() => setIsIntakeOpen(false)}
        onSuccess={() => refetch()}
      />

      {/* Patient 360 Drawer */}
      <PatientDrawer
        isOpen={Boolean(selectedPatient)}
        onClose={() => setSelectedPatient(null)}
        patient={selectedPatient}
      />
    </div>
  );
};
