import React, { useRef } from "react";
import ReactDOM from "react-dom/client";
import {
  BrowserRouter,
  Routes,
  Route,
  Navigate,
  useLocation,
} from "react-router-dom";
import { CSSTransition, SwitchTransition } from "react-transition-group";
import SideBar from "./components/SideBar";
import LoginPage from "./pages/Login";
import RegisterPage from "./pages/Register";
import Navbar from "./components/Navbar";
import DashboardPage from "./pages/protected/Dashboard";
import LicensesPage from "./pages/protected/License";
import Health from "./pages/Health"; // Import the HealthCheck component
import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import "./stylesheets/App.css";

const App: React.FC = () => {
  const location = useLocation();
  const nodeRef = useRef(null);

  return (
    <div>
      <ToastContainer />
      <Navbar />
      <SideBar />
      <SwitchTransition>
        <CSSTransition
          key={location.key}
          timeout={500}
          classNames="fade"
          unmountOnExit
          nodeRef={nodeRef}
        >
          <div ref={nodeRef} style={{ position: "relative" }}>
            <Routes location={location}>
              <Route path="/login" element={<LoginPage />} />
              <Route path="/register" element={<RegisterPage />} />
              <Route path="/dashboard" element={<DashboardPage />} />
              <Route path="/dashboard/licenses" element={<LicensesPage />} />
              <Route path="/health" element={<Health />} />
              <Route
                path="*"
                element={<Navigate to="/dashboard/licenses" replace />}
              />
            </Routes>
          </div>
        </CSSTransition>
      </SwitchTransition>
    </div>
  );
};

const Root: React.FC = () => (
  <BrowserRouter>
    <App />
  </BrowserRouter>
);

const rootElement = document.getElementById("app");
if (rootElement) {
  const root = ReactDOM.createRoot(rootElement);
  root.render(
    <React.StrictMode>
      <Root />
    </React.StrictMode>
  );
}

export default Root;
