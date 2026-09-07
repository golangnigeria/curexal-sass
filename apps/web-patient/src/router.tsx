import React from "react";
import { createBrowserRouter, Navigate } from "react-router-dom";
import { PatientPortalLayout } from "./layouts/portal-layout";
import { PatientPortalLoginPage } from "./pages/login";
import { PatientPortalDashboardPage } from "./pages/dashboard";
import { NewConsultationBookingPage } from "./pages/consultations/new";
import { PatientPortalTelehealthRoomPage } from "./pages/consultations/room";

function ProtectedPatientRoute({ children }: { children: React.ReactNode }) {
  const token = typeof window !== "undefined" ? localStorage.getItem("curexal_portal_token") : null;
  if (!token) {
    return <Navigate to="/login" replace />;
  }
  return <>{children}</>;
}

export const router = createBrowserRouter([
  {
    path: "/login",
    element: <PatientPortalLoginPage />,
  },
  {
    path: "/",
    element: (
      <ProtectedPatientRoute>
        <PatientPortalLayout />
      </ProtectedPatientRoute>
    ),
    children: [
      {
        index: true,
        element: <Navigate to="/dashboard" replace />,
      },
      {
        path: "dashboard",
        element: <PatientPortalDashboardPage />,
      },
      {
        path: "consultations/new",
        element: <NewConsultationBookingPage />,
      },
      {
        path: "consultations/:requestId/room",
        element: <PatientPortalTelehealthRoomPage />,
      },
    ],
  },
  {
    path: "*",
    element: <Navigate to="/dashboard" replace />,
  },
]);
