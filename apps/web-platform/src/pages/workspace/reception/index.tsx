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
  CalendarPlus,
  Clock,
  Video,
  Building2,
  X,
  Loader2,
} from "lucide-react";
import { usePatients, useSendPortalOTP } from "../../../api/hooks/use-patients";
import { useCreateCareRequest } from "../../../api/hooks/use-care-orchestration";
import { useCreateAppointment, useProviderProfiles } from "../../../api/hooks/use-appointments";
import { PatientIntakeModal } from "../../../components/patients/patient-intake-modal";
import { PatientDrawer } from "@/features/patients/patient-drawer";
import type { CanonicalPatient } from "../../../api/contracts";

export const ReceptionWorkspacePage: React.FC = () => {
  const { branchSlug } = useParams<{ branchSlug: string }>();
  const navigate = useNavigate();

  const [searchQuery, setSearchQuery] = useState("");
  const [statusFilter, setStatusFilter] = useState("");
  const [isIntakeOpen, setIsIntakeOpen] = useState(false);
  const [feedbackNotice, setFeedbackNotice] = useState<string | null>(null);
  const [selectedPatient, setSelectedPatient] = useState<CanonicalPatient | null>(null);

  // Quick Check-In State
  const [checkInPatient, setCheckInPatient] = useState<CanonicalPatient | null>(null);
  const [checkInService, setCheckInService] = useState("GENERAL_CONSULTATION");
  const [checkInChannel, setCheckInChannel] = useState("in_person");
  const [checkInComplaint, setCheckInComplaint] = useState("General consultation and physical check");

  // Book Appointment State
  const [schedulePatient, setSchedulePatient] = useState<CanonicalPatient | null>(null);
  const [selectedProviderId, setSelectedProviderId] = useState("");
  const [appointmentChannel, setAppointmentChannel] = useState("in_person");
  const [appointmentStartTime, setAppointmentStartTime] = useState("");
  const [appointmentReason, setAppointmentReason] = useState("");

  const { data: patientList, isLoading, refetch } = usePatients({
    query: searchQuery || undefined,
    status: statusFilter || undefined,
    limit: 50,
  });

  const { data: providers = [] } = useProviderProfiles("ON_DUTY");

  const sendPortalOTP = useSendPortalOTP();
  const createCareRequest = useCreateCareRequest();
  const createAppointment = useCreateAppointment();

  const handleSendOTP = (patient: CanonicalPatient) => {
    const primaryPhone =
      patient.contacts?.find((c) => c.system === "PHONE")?.value || patient.mrn;
    sendPortalOTP.mutate(
      { identifier: primaryPhone },
      {
        onSuccess: () => {
          setFeedbackNotice(
            `Portal login OTP sent to ${patient.firstName} (${primaryPhone})`
          );
          setTimeout(() => setFeedbackNotice(null), 5000);
        },
      }
    );
  };

  const handleConfirmCheckIn = (e: React.FormEvent) => {
    e.preventDefault();
    if (!checkInPatient) return;

    createCareRequest.mutate(
      {
        patientId: checkInPatient.id,
        serviceType: checkInService,
        preferredMode: checkInChannel === "video" ? "VIDEO" : "IN_PERSON",
        chiefComplaint: checkInComplaint || "General consultation",
      },
      {
        onSuccess: (res: any) => {
          const reqNum = res?.requestNumber || res?.data?.requestNumber || "REQ-LIVE";
          setFeedbackNotice(
            `Patient ${checkInPatient.firstName} ${checkInPatient.lastName} successfully queued on Care Desk live board (${reqNum})`
          );
          setCheckInPatient(null);
          setTimeout(() => setFeedbackNotice(null), 6000);
        },
      }
    );
  };

  const handleConfirmAppointment = (e: React.FormEvent) => {
    e.preventDefault();
    if (!schedulePatient || !selectedProviderId || !appointmentStartTime) return;

    const start = new Date(appointmentStartTime);
    const end = new Date(start.getTime() + 30 * 60000); // 30 min duration default

    createAppointment.mutate(
      {
        patientId: schedulePatient.id,
        providerId: selectedProviderId,
        serviceType: "CONSULTATION",
        deliveryChannel: appointmentChannel,
        startTime: start.toISOString(),
        endTime: end.toISOString(),
        reasonForVisit: appointmentReason || "Scheduled clinical consultation",
      },
      {
        onSuccess: (res: any) => {
          const aptNum = res?.appointmentNumber || res?.data?.appointmentNumber || "APT-BOOKED";
          setFeedbackNotice(
            `Appointment scheduled successfully! Reference ID: ${aptNum}`
          );
          setSchedulePatient(null);
          setAppointmentStartTime("");
          setAppointmentReason("");
          setTimeout(() => setFeedbackNotice(null), 6000);
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

      {/* Action Notice Banner */}
      {feedbackNotice && (
        <div className="p-4 rounded-xl bg-emerald-500/10 border border-emerald-500/30 text-emerald-300 text-xs flex items-center justify-between animate-in slide-in-from-top-2">
          <div className="flex items-center gap-2">
            <CheckCircle className="w-4 h-4 text-emerald-400" />
            <span>{feedbackNotice}</span>
          </div>
          <button
            onClick={() => navigate(`/${branchSlug}/care-desk`)}
            className="flex items-center gap-1 text-cyan-300 hover:text-cyan-200 underline text-xs font-semibold"
          >
            <span>Open Care Desk</span>
            <ArrowRight className="w-3.5 h-3.5" />
          </button>
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
            <div className="text-xs text-slate-400 font-medium">On-Duty Clinical Providers</div>
            <div className="text-2xl font-bold text-blue-400 mt-1">{providers.length} Available</div>
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
                  <div className="flex flex-wrap items-center gap-2 self-end md:self-center">
                    <button
                      onClick={() => setCheckInPatient(patient)}
                      title="Check In to Live Care Desk Queue"
                      className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-emerald-500/10 hover:bg-emerald-500/20 text-emerald-300 text-xs font-semibold border border-emerald-500/30 transition-colors shadow-sm"
                    >
                      <Clock className="w-3.5 h-3.5 text-emerald-400" />
                      <span>Check In</span>
                    </button>
                    <button
                      onClick={() => {
                        setSchedulePatient(patient);
                        if (providers.length > 0 && !selectedProviderId) {
                          setSelectedProviderId(providers[0].id);
                        }
                      }}
                      title="Schedule Calendar Appointment"
                      className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-indigo-500/10 hover:bg-indigo-500/20 text-indigo-300 text-xs font-semibold border border-indigo-500/30 transition-colors"
                    >
                      <CalendarPlus className="w-3.5 h-3.5 text-indigo-400" />
                      <span>Book Slot</span>
                    </button>
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
                      <span>Send OTP</span>
                    </button>
                    <button
                      onClick={() => navigate(`/${branchSlug}/clinical?patientId=${patient.id}`)}
                      className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-cyan-500/10 hover:bg-cyan-500/20 text-cyan-300 text-xs font-medium border border-cyan-500/30 transition-colors"
                    >
                      <Stethoscope className="w-3.5 h-3.5" />
                      <span>Encounter</span>
                    </button>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>

      {/* QUICK CHECK-IN MODAL */}
      {checkInPatient && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/70 backdrop-blur-sm animate-in fade-in duration-200">
          <div className="relative w-full max-w-lg bg-slate-900 border border-slate-800 rounded-2xl shadow-2xl p-6 text-slate-100 space-y-4">
            <div className="flex items-center justify-between pb-3 border-b border-slate-800">
              <div className="flex items-center gap-3">
                <div className="p-2.5 rounded-xl bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
                  <Clock className="w-5 h-5" />
                </div>
                <div>
                  <h3 className="text-base font-bold text-white">Queue Check-In</h3>
                  <p className="text-xs text-slate-400">
                    Routing {checkInPatient.firstName} {checkInPatient.lastName} ({checkInPatient.mrn}) to Live Desk
                  </p>
                </div>
              </div>
              <button
                onClick={() => setCheckInPatient(null)}
                className="p-2 rounded-lg text-slate-400 hover:text-white hover:bg-slate-800"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <form onSubmit={handleConfirmCheckIn} className="space-y-4">
              <div>
                <label className="block text-xs font-medium text-slate-400 mb-1">Service Required</label>
                <select
                  value={checkInService}
                  onChange={(e) => setCheckInService(e.target.value)}
                  className="w-full px-3 py-2 bg-slate-800 border border-slate-700 rounded-lg text-sm text-white focus:outline-none focus:border-cyan-500"
                >
                  <option value="GENERAL_CONSULTATION">General Consultation</option>
                  <option value="SPECIALIST_CONSULTATION">Specialist Review</option>
                  <option value="FOLLOW_UP">Follow-up Review</option>
                  <option value="EMERGENCY_TRIAGE">Emergency / Urgent Triage</option>
                </select>
              </div>

              <div>
                <label className="block text-xs font-medium text-slate-400 mb-1">Delivery Channel</label>
                <div className="grid grid-cols-2 gap-3">
                  <button
                    type="button"
                    onClick={() => setCheckInChannel("in_person")}
                    className={`flex items-center justify-center gap-2 p-3 rounded-xl border text-xs font-semibold transition-all ${
                      checkInChannel === "in_person"
                        ? "bg-blue-500/20 border-blue-500/50 text-white shadow-lg"
                        : "bg-slate-800/60 border-slate-700 text-slate-400 hover:bg-slate-800"
                    }`}
                  >
                    <Building2 className="w-4 h-4 text-blue-400" />
                    <span>In-Person (Room 104)</span>
                  </button>
                  <button
                    type="button"
                    onClick={() => setCheckInChannel("video")}
                    className={`flex items-center justify-center gap-2 p-3 rounded-xl border text-xs font-semibold transition-all ${
                      checkInChannel === "video"
                        ? "bg-purple-500/20 border-purple-500/50 text-white shadow-lg"
                        : "bg-slate-800/60 border-slate-700 text-slate-400 hover:bg-slate-800"
                    }`}
                  >
                    <Video className="w-4 h-4 text-purple-400" />
                    <span>Telehealth (WebRTC)</span>
                  </button>
                </div>
              </div>

              <div>
                <label className="block text-xs font-medium text-slate-400 mb-1">Chief Complaint / Notes</label>
                <textarea
                  rows={3}
                  value={checkInComplaint}
                  onChange={(e) => setCheckInComplaint(e.target.value)}
                  placeholder="Reason for visit or presenting symptoms..."
                  className="w-full px-3 py-2 bg-slate-800 border border-slate-700 rounded-lg text-sm text-white focus:outline-none focus:border-cyan-500"
                />
              </div>

              <div className="pt-3 border-t border-slate-800 flex justify-end gap-3">
                <button
                  type="button"
                  onClick={() => setCheckInPatient(null)}
                  className="px-4 py-2 rounded-lg text-xs font-semibold bg-slate-800 hover:bg-slate-700 text-slate-300"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={createCareRequest.isPending}
                  className="flex items-center gap-2 px-5 py-2 rounded-lg text-xs font-semibold bg-gradient-to-r from-emerald-500 to-teal-600 hover:from-emerald-400 hover:to-teal-500 text-white shadow-lg transition-all"
                >
                  {createCareRequest.isPending ? (
                    <>
                      <Loader2 className="w-4 h-4 animate-spin" />
                      <span>Checking In...</span>
                    </>
                  ) : (
                    <>
                      <Clock className="w-4 h-4" />
                      <span>Confirm Check-In</span>
                    </>
                  )}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* BOOK APPOINTMENT MODAL */}
      {schedulePatient && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/70 backdrop-blur-sm animate-in fade-in duration-200">
          <div className="relative w-full max-w-lg bg-slate-900 border border-slate-800 rounded-2xl shadow-2xl p-6 text-slate-100 space-y-4">
            <div className="flex items-center justify-between pb-3 border-b border-slate-800">
              <div className="flex items-center gap-3">
                <div className="p-2.5 rounded-xl bg-indigo-500/10 text-indigo-400 border border-indigo-500/20">
                  <CalendarPlus className="w-5 h-5" />
                </div>
                <div>
                  <h3 className="text-base font-bold text-white">Schedule Appointment</h3>
                  <p className="text-xs text-slate-400">
                    Book slot for {schedulePatient.firstName} {schedulePatient.lastName} ({schedulePatient.mrn})
                  </p>
                </div>
              </div>
              <button
                onClick={() => setSchedulePatient(null)}
                className="p-2 rounded-lg text-slate-400 hover:text-white hover:bg-slate-800"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <form onSubmit={handleConfirmAppointment} className="space-y-4">
              <div>
                <label className="block text-xs font-medium text-slate-400 mb-1">Attending Clinical Doctor</label>
                <select
                  value={selectedProviderId}
                  onChange={(e) => setSelectedProviderId(e.target.value)}
                  className="w-full px-3 py-2 bg-slate-800 border border-slate-700 rounded-lg text-sm text-white focus:outline-none focus:border-cyan-500"
                  required
                >
                  {providers.length === 0 ? (
                    <option value="">No on-duty doctors registered</option>
                  ) : (
                    providers.map((p) => (
                      <option key={p.id} value={p.id}>
                        {p.providerName || "Consulting Doctor"} — {p.specialtyCode} ({p.roomNumber || "Suite 104"})
                      </option>
                    ))
                  )}
                </select>
              </div>

              <div>
                <label className="block text-xs font-medium text-slate-400 mb-1">Delivery Channel</label>
                <div className="grid grid-cols-2 gap-3">
                  <button
                    type="button"
                    onClick={() => setAppointmentChannel("in_person")}
                    className={`flex items-center justify-center gap-2 p-3 rounded-xl border text-xs font-semibold transition-all ${
                      appointmentChannel === "in_person"
                        ? "bg-blue-500/20 border-blue-500/50 text-white shadow-lg"
                        : "bg-slate-800/60 border-slate-700 text-slate-400 hover:bg-slate-800"
                    }`}
                  >
                    <Building2 className="w-4 h-4 text-blue-400" />
                    <span>In-Person Consulting</span>
                  </button>
                  <button
                    type="button"
                    onClick={() => setAppointmentChannel("video")}
                    className={`flex items-center justify-center gap-2 p-3 rounded-xl border text-xs font-semibold transition-all ${
                      appointmentChannel === "video"
                        ? "bg-purple-500/20 border-purple-500/50 text-white shadow-lg"
                        : "bg-slate-800/60 border-slate-700 text-slate-400 hover:bg-slate-800"
                    }`}
                  >
                    <Video className="w-4 h-4 text-purple-400" />
                    <span>Telehealth Video</span>
                  </button>
                </div>
              </div>

              <div>
                <label className="block text-xs font-medium text-slate-400 mb-1">Appointment Date & Time</label>
                <input
                  type="datetime-local"
                  value={appointmentStartTime}
                  onChange={(e) => setAppointmentStartTime(e.target.value)}
                  className="w-full px-3 py-2 bg-slate-800 border border-slate-700 rounded-lg text-sm text-white focus:outline-none focus:border-cyan-500"
                  required
                />
              </div>

              <div>
                <label className="block text-xs font-medium text-slate-400 mb-1">Reason for Visit</label>
                <input
                  type="text"
                  value={appointmentReason}
                  onChange={(e) => setAppointmentReason(e.target.value)}
                  placeholder="e.g. Hypertension review, cardiology follow-up..."
                  className="w-full px-3 py-2 bg-slate-800 border border-slate-700 rounded-lg text-sm text-white focus:outline-none focus:border-cyan-500"
                />
              </div>

              <div className="pt-3 border-t border-slate-800 flex justify-end gap-3">
                <button
                  type="button"
                  onClick={() => setSchedulePatient(null)}
                  className="px-4 py-2 rounded-lg text-xs font-semibold bg-slate-800 hover:bg-slate-700 text-slate-300"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={createAppointment.isPending || !selectedProviderId || !appointmentStartTime}
                  className="flex items-center gap-2 px-5 py-2 rounded-lg text-xs font-semibold bg-gradient-to-r from-indigo-500 to-blue-600 hover:from-indigo-400 hover:to-blue-500 text-white shadow-lg transition-all disabled:opacity-50"
                >
                  {createAppointment.isPending ? (
                    <>
                      <Loader2 className="w-4 h-4 animate-spin" />
                      <span>Reserving Slot...</span>
                    </>
                  ) : (
                    <>
                      <CalendarPlus className="w-4 h-4" />
                      <span>Book Appointment</span>
                    </>
                  )}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

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
