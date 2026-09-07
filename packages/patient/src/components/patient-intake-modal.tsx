import React, { useState, useEffect } from "react";
import {
  UserPlus,
  AlertTriangle,
  CheckCircle2,
  Phone,
  Mail,
  ExternalLink,
  Loader2,
  X,
} from "lucide-react";
import { useResolveDuplicates, useRegisterPatient } from "../hooks/use-patients";
import type {
  DuplicateEvaluationResponse,
  RegisterCanonicalPatientPayload,
} from "@curexal/contracts";

export interface PatientIntakeModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess?: (patient?: any) => void;
}

export const PatientIntakeModal: React.FC<PatientIntakeModalProps> = ({
  isOpen,
  onClose,
  onSuccess,
}) => {
  const [firstName, setFirstName] = useState("");
  const [middleName, setMiddleName] = useState("");
  const [lastName, setLastName] = useState("");
  const [gender, setGender] = useState("MALE");
  const [dob, setDob] = useState("");
  const [phone, setPhone] = useState("");
  const [email, setEmail] = useState("");
  const [nin, setNin] = useState("");
  const [bloodGroup, setBloodGroup] = useState("");
  const [genotype, setGenotype] = useState("");
  const [address, setAddress] = useState("");
  const [forceOverride, setForceOverride] = useState(false);

  const resolveDuplicates = useResolveDuplicates();
  const registerPatient = useRegisterPatient();

  const [mpiResults, setMpiResults] = useState<DuplicateEvaluationResponse | null>(null);

  // Debounced duplicate evaluation
  useEffect(() => {
    if ((phone.length >= 8 || nin.length >= 6) && lastName.length >= 2) {
      const timer = setTimeout(() => {
        resolveDuplicates.mutate(
          {
            firstName,
            middleName: middleName || undefined,
            lastName,
            dateOfBirth: dob || undefined,
            phone,
            nin: nin || undefined,
          },
          {
            onSuccess: (data) => setMpiResults(data),
          }
        );
      }, 500);
      return () => clearTimeout(timer);
    } else {
      setMpiResults(null);
    }
  }, [phone, nin, lastName, dob, firstName]);

  if (!isOpen) return null;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!firstName || !lastName || !dob || !phone) return;

    const payload: RegisterCanonicalPatientPayload = {
      firstName,
      middleName: middleName || undefined,
      lastName,
      gender,
      dateOfBirth: dob,
      phone,
      email: email || undefined,
      nin: nin || undefined,
      bloodGroup: bloodGroup || undefined,
      genotype: genotype || undefined,
      address: address || undefined,
      registrationChannel: "RECEPTION",
      forceRegistration: forceOverride,
    };

    registerPatient.mutate(payload, {
      onSuccess: (created: any) => {
        if (onSuccess) onSuccess(created?.patient || created);
        onClose();
      },
    });
  };

  const hasDuplicateWarning =
    mpiResults &&
    (mpiResults.matchStatus === "EXACT_MATCH" ||
      mpiResults.matchStatus === "PROBABLE_DUPLICATE");

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/70 backdrop-blur-sm animate-in fade-in duration-200">
      <div className="relative w-full max-w-3xl max-h-[90vh] overflow-y-auto bg-slate-900 border border-slate-800 rounded-2xl shadow-2xl text-slate-100 flex flex-col">
        {/* Header */}
        <div className="sticky top-0 z-10 flex items-center justify-between px-6 py-4 bg-slate-900/90 backdrop-blur-md border-b border-slate-800">
          <div className="flex items-center gap-3">
            <div className="p-2.5 rounded-xl bg-cyan-500/10 text-cyan-400 border border-cyan-500/20">
              <UserPlus className="w-5 h-5" />
            </div>
            <div>
              <h2 className="text-lg font-semibold tracking-tight text-white">
                New Patient Intake & Registration
              </h2>
              <p className="text-xs text-slate-400">
                Master Patient Index duplicate resolution & canonical record creation
              </p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="p-2 rounded-lg text-slate-400 hover:text-white hover:bg-slate-800 transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Content Form */}
        <form onSubmit={handleSubmit} className="p-6 space-y-6 flex-1">
          {/* MPI Duplicate Alert Banner */}
          {hasDuplicateWarning && (
            <div className="p-4 rounded-xl bg-amber-500/10 border border-amber-500/30 text-amber-200 space-y-3 animate-in slide-in-from-top-2">
              <div className="flex items-start gap-3">
                <AlertTriangle className="w-5 h-5 text-amber-400 shrink-0 mt-0.5" />
                <div className="flex-1 text-xs space-y-1">
                  <div className="font-semibold text-amber-300">
                    Master Patient Index: Probable Duplicate Found ({mpiResults.matchStatus})
                  </div>
                  <p className="text-amber-200/80">
                    The following existing registered patient(s) share identical phone, national ID, or demographic signals:
                  </p>
                  <div className="mt-2 space-y-2">
                    {mpiResults.candidates.map((cand) => (
                      <div
                        key={cand.patientId}
                        className="flex items-center justify-between p-2.5 rounded-lg bg-slate-950/60 border border-amber-500/20 text-xs"
                      >
                        <div>
                          <span className="font-medium text-white">
                            {cand.firstName} {cand.lastName}
                          </span>{" "}
                          <span className="text-cyan-400 font-mono">({cand.mrn})</span>
                          <div className="text-[10px] text-slate-400 mt-0.5">
                            DOB: {cand.dateOfBirth} • Signals: {cand.matchedSignals.join(", ")} ({cand.confidenceScore}% confidence)
                          </div>
                        </div>
                        <button
                          type="button"
                          onClick={() => {
                            window.open(`/:branchSlug/patients/${cand.patientId}`, "_blank");
                          }}
                          className="flex items-center gap-1.5 px-2.5 py-1 rounded-md bg-amber-500/20 text-amber-300 hover:bg-amber-500/30 font-medium text-xs transition-colors"
                        >
                          <span>Open Record</span>
                          <ExternalLink className="w-3 h-3" />
                        </button>
                      </div>
                    ))}
                  </div>
                  <div className="pt-2 flex items-center gap-2">
                    <input
                      type="checkbox"
                      id="forceOverride"
                      checked={forceOverride}
                      onChange={(e) => setForceOverride(e.target.checked)}
                      className="rounded border-slate-700 text-cyan-500 focus:ring-cyan-500/30 bg-slate-800"
                    />
                    <label htmlFor="forceOverride" className="text-xs text-amber-300 font-medium cursor-pointer">
                      Confirm this is a different individual (Force create distinct canonical record)
                    </label>
                  </div>
                </div>
              </div>
            </div>
          )}

          {/* Demographics Section */}
          <div className="space-y-4">
            <h3 className="text-xs font-semibold uppercase tracking-wider text-slate-400 flex items-center gap-2">
              <span className="w-1.5 h-1.5 rounded-full bg-cyan-400"></span>
              Primary Demographics
            </h3>
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
              <div>
                <label className="block text-xs font-medium text-slate-300 mb-1">
                  First Name <span className="text-rose-400">*</span>
                </label>
                <input
                  type="text"
                  required
                  value={firstName}
                  onChange={(e) => setFirstName(e.target.value)}
                  placeholder="e.g. Amina"
                  className="w-full px-3 py-2 bg-slate-800/80 border border-slate-700 rounded-lg text-sm text-white placeholder-slate-500 focus:outline-none focus:border-cyan-500 focus:ring-1 focus:ring-cyan-500 transition-colors"
                />
              </div>
              <div>
                <label className="block text-xs font-medium text-slate-300 mb-1">
                  Middle Name
                </label>
                <input
                  type="text"
                  value={middleName}
                  onChange={(e) => setMiddleName(e.target.value)}
                  placeholder="e.g. Joy"
                  className="w-full px-3 py-2 bg-slate-800/80 border border-slate-700 rounded-lg text-sm text-white placeholder-slate-500 focus:outline-none focus:border-cyan-500 focus:ring-1 focus:ring-cyan-500 transition-colors"
                />
              </div>
              <div>
                <label className="block text-xs font-medium text-slate-300 mb-1">
                  Last Name <span className="text-rose-400">*</span>
                </label>
                <input
                  type="text"
                  required
                  value={lastName}
                  onChange={(e) => setLastName(e.target.value)}
                  placeholder="e.g. Yusuf"
                  className="w-full px-3 py-2 bg-slate-800/80 border border-slate-700 rounded-lg text-sm text-white placeholder-slate-500 focus:outline-none focus:border-cyan-500 focus:ring-1 focus:ring-cyan-500 transition-colors"
                />
              </div>
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label className="block text-xs font-medium text-slate-300 mb-1">
                  Gender <span className="text-rose-400">*</span>
                </label>
                <select
                  value={gender}
                  onChange={(e) => setGender(e.target.value)}
                  className="w-full px-3 py-2 bg-slate-800/80 border border-slate-700 rounded-lg text-sm text-white focus:outline-none focus:border-cyan-500 focus:ring-1 focus:ring-cyan-500 transition-colors"
                >
                  <option value="MALE">Male</option>
                  <option value="FEMALE">Female</option>
                  <option value="OTHER">Other</option>
                </select>
              </div>
              <div>
                <label className="block text-xs font-medium text-slate-300 mb-1">
                  Date of Birth <span className="text-rose-400">*</span>
                </label>
                <div className="relative">
                  <input
                    type="date"
                    required
                    value={dob}
                    onChange={(e) => setDob(e.target.value)}
                    className="w-full px-3 py-2 bg-slate-800/80 border border-slate-700 rounded-lg text-sm text-white focus:outline-none focus:border-cyan-500 focus:ring-1 focus:ring-cyan-500 transition-colors"
                  />
                </div>
              </div>
            </div>
          </div>

          {/* Telecoms & Identification Section */}
          <div className="space-y-4 pt-2">
            <h3 className="text-xs font-semibold uppercase tracking-wider text-slate-400 flex items-center gap-2">
              <span className="w-1.5 h-1.5 rounded-full bg-cyan-400"></span>
              Contact & National Identification
            </h3>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label className="block text-xs font-medium text-slate-300 mb-1">
                  Primary Mobile Phone <span className="text-rose-400">*</span>
                </label>
                <div className="relative">
                  <Phone className="w-4 h-4 text-slate-400 absolute left-3 top-2.5" />
                  <input
                    type="tel"
                    required
                    value={phone}
                    onChange={(e) => setPhone(e.target.value)}
                    placeholder="+234 801 234 5678"
                    className="w-full pl-9 pr-3 py-2 bg-slate-800/80 border border-slate-700 rounded-lg text-sm text-white placeholder-slate-500 focus:outline-none focus:border-cyan-500 focus:ring-1 focus:ring-cyan-500 transition-colors"
                  />
                </div>
              </div>
              <div>
                <label className="block text-xs font-medium text-slate-300 mb-1">
                  Email Address
                </label>
                <div className="relative">
                  <Mail className="w-4 h-4 text-slate-400 absolute left-3 top-2.5" />
                  <input
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    placeholder="patient@example.com"
                    className="w-full pl-9 pr-3 py-2 bg-slate-800/80 border border-slate-700 rounded-lg text-sm text-white placeholder-slate-500 focus:outline-none focus:border-cyan-500 focus:ring-1 focus:ring-cyan-500 transition-colors"
                  />
                </div>
              </div>
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
              <div>
                <label className="block text-xs font-medium text-slate-300 mb-1">
                  National ID (NIN)
                </label>
                <input
                  type="text"
                  value={nin}
                  onChange={(e) => setNin(e.target.value)}
                  placeholder="11-digit NIN"
                  className="w-full px-3 py-2 bg-slate-800/80 border border-slate-700 rounded-lg text-sm text-white placeholder-slate-500 focus:outline-none focus:border-cyan-500 focus:ring-1 focus:ring-cyan-500 transition-colors"
                />
              </div>
              <div>
                <label className="block text-xs font-medium text-slate-300 mb-1">
                  Blood Group
                </label>
                <select
                  value={bloodGroup}
                  onChange={(e) => setBloodGroup(e.target.value)}
                  className="w-full px-3 py-2 bg-slate-800/80 border border-slate-700 rounded-lg text-sm text-white focus:outline-none focus:border-cyan-500 focus:ring-1 focus:ring-cyan-500 transition-colors"
                >
                  <option value="">Select Blood Group</option>
                  <option value="A+">A+</option>
                  <option value="A-">A-</option>
                  <option value="B+">B+</option>
                  <option value="B-">B-</option>
                  <option value="AB+">AB+</option>
                  <option value="AB-">AB-</option>
                  <option value="O+">O+</option>
                  <option value="O-">O-</option>
                </select>
              </div>
              <div>
                <label className="block text-xs font-medium text-slate-300 mb-1">
                  Genotype
                </label>
                <select
                  value={genotype}
                  onChange={(e) => setGenotype(e.target.value)}
                  className="w-full px-3 py-2 bg-slate-800/80 border border-slate-700 rounded-lg text-sm text-white focus:outline-none focus:border-cyan-500 focus:ring-1 focus:ring-cyan-500 transition-colors"
                >
                  <option value="">Select Genotype</option>
                  <option value="AA">AA</option>
                  <option value="AS">AS</option>
                  <option value="SS">SS</option>
                  <option value="AC">AC</option>
                </select>
              </div>
            </div>

            <div>
              <label className="block text-xs font-medium text-slate-300 mb-1">
                Residential Address
              </label>
              <input
                type="text"
                value={address}
                onChange={(e) => setAddress(e.target.value)}
                placeholder="Street name, City, State"
                className="w-full px-3 py-2 bg-slate-800/80 border border-slate-700 rounded-lg text-sm text-white placeholder-slate-500 focus:outline-none focus:border-cyan-500 focus:ring-1 focus:ring-cyan-500 transition-colors"
              />
            </div>
          </div>

          {/* Footer Actions */}
          <div className="pt-4 border-t border-slate-800 flex items-center justify-between">
            <div className="text-[11px] text-slate-400 flex items-center gap-1.5">
              <CheckCircle2 className="w-3.5 h-3.5 text-cyan-400" />
              <span>Auto-provisions zero-friction Portal Access token via SMS</span>
            </div>
            <div className="flex items-center gap-3">
              <button
                type="button"
                onClick={onClose}
                className="px-4 py-2 text-xs font-medium text-slate-300 hover:text-white bg-slate-800 hover:bg-slate-700 rounded-lg transition-colors"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={
                  registerPatient.isPending ||
                  Boolean(hasDuplicateWarning && !forceOverride)
                }
                className="flex items-center gap-2 px-5 py-2 text-xs font-semibold text-white bg-gradient-to-r from-cyan-500 to-blue-600 hover:from-cyan-400 hover:to-blue-500 disabled:opacity-50 disabled:cursor-not-allowed rounded-lg shadow-lg shadow-cyan-500/20 transition-all"
              >
                {registerPatient.isPending ? (
                  <>
                    <Loader2 className="w-4 h-4 animate-spin" />
                    <span>Registering Patient...</span>
                  </>
                ) : (
                  <>
                    <UserPlus className="w-4 h-4" />
                    <span>Complete Registration</span>
                  </>
                )}
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
};
