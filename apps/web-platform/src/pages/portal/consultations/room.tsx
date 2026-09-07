import React, { useState } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import {
  Video,
  VideoOff,
  Mic,
  MicOff,
  PhoneOff,
  Stethoscope,
  Clock,
  ShieldCheck,
  ArrowLeft,
  Sparkles,
  CheckCircle2,
} from "lucide-react";

export const PatientPortalTelehealthRoomPage: React.FC = () => {
  const { requestId } = useParams<{ requestId: string }>();
  const navigate = useNavigate();

  const [isVideoOn, setIsVideoOn] = useState(true);
  const [isMicOn, setIsMicOn] = useState(true);

  const handleLeaveRoom = () => {
    navigate("/portal/dashboard");
  };

  return (
    <div className="max-w-4xl mx-auto space-y-6 animate-in fade-in duration-300">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Link
            to="/portal/dashboard"
            className="p-2 rounded-xl text-slate-400 hover:text-white bg-slate-800 transition-colors"
          >
            <ArrowLeft className="w-4 h-4" />
          </Link>
          <div>
            <h1 className="text-xl font-bold tracking-tight text-white flex items-center gap-2">
              <Video className="w-5 h-5 text-cyan-400" />
              <span>Telehealth Virtual Room</span>
            </h1>
            <p className="text-xs text-slate-400">
              Encrypted live consultation with your attending physician
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2 px-3 py-1 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 text-xs font-mono">
          <span className="w-2 h-2 rounded-full bg-emerald-400 animate-ping" />
          <span>Connected (00:14:32)</span>
        </div>
      </div>

      {/* Video Stage */}
      <div className="relative aspect-video rounded-3xl bg-slate-900 border border-slate-800 overflow-hidden shadow-2xl flex flex-col justify-between p-6">
        <div className="absolute inset-0 bg-gradient-to-br from-slate-950 via-slate-900 to-blue-950 flex items-center justify-center">
          <div className="text-center space-y-3">
            <div className="w-20 h-20 rounded-3xl bg-cyan-500/20 border border-cyan-500/30 text-cyan-300 flex items-center justify-center mx-auto shadow-2xl animate-pulse">
              <Stethoscope className="w-10 h-10" />
            </div>
            <div>
              <div className="text-base font-bold text-white">Dr. Attending Physician</div>
              <div className="text-xs text-cyan-400 font-medium">General & Internal Medicine Specialist</div>
            </div>
            <div className="text-[11px] text-slate-400 font-mono">
              Secure WebRTC Audio/Video Stream Active
            </div>
          </div>
        </div>

        {/* Self PIP */}
        <div className="absolute right-6 bottom-20 w-32 h-20 rounded-2xl bg-slate-800/90 border border-slate-700 shadow-xl flex items-center justify-center text-xs text-slate-400 font-medium">
          {isVideoOn ? "Your Camera" : "Camera Off"}
        </div>

        {/* Bottom Control Bar */}
        <div className="relative z-10 flex items-center justify-center gap-4 py-2">
          <button
            type="button"
            onClick={() => setIsMicOn(!isMicOn)}
            className={`p-3.5 rounded-2xl transition-all shadow-lg ${
              isMicOn ? "bg-slate-800 text-white hover:bg-slate-700" : "bg-rose-500 text-white"
            }`}
            title={isMicOn ? "Mute Microphone" : "Unmute Microphone"}
          >
            {isMicOn ? <Mic className="w-5 h-5" /> : <MicOff className="w-5 h-5" />}
          </button>

          <button
            type="button"
            onClick={() => setIsVideoOn(!isVideoOn)}
            className={`p-3.5 rounded-2xl transition-all shadow-lg ${
              isVideoOn ? "bg-slate-800 text-white hover:bg-slate-700" : "bg-rose-500 text-white"
            }`}
            title={isVideoOn ? "Turn Camera Off" : "Turn Camera On"}
          >
            {isVideoOn ? <Video className="w-5 h-5" /> : <VideoOff className="w-5 h-5" />}
          </button>

          <button
            type="button"
            onClick={handleLeaveRoom}
            className="p-3.5 rounded-2xl bg-rose-500 hover:bg-rose-600 text-white transition-all shadow-lg shadow-rose-500/25"
            title="Leave Consultation"
          >
            <PhoneOff className="w-5 h-5" />
          </button>
        </div>
      </div>

      {/* Security notice */}
      <div className="p-4 rounded-2xl bg-slate-900 border border-slate-800 text-xs text-slate-400 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <ShieldCheck className="w-4 h-4 text-emerald-400" />
          <span>This virtual encounter is encrypted and compliant with medical privacy standards.</span>
        </div>
        <span className="text-[11px] text-slate-500">ID: {requestId || "REQ-2026-00412"}</span>
      </div>
    </div>
  );
};
