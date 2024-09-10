import { useRef } from "react";
import { effect, signal } from "@preact/signals-react";
import axios from "axios";
import CustomButton from "../../components/licenses/CustomButton";
import useLoadingEffect from "../../hooks/useLoadingEffect";
import LoadingSpinner, {
  isLoadingDone,
} from "../../components/misc/LoadingSpinner";
import useKeycloakData from "../../hooks/useKeycloakData";
import keycloak from "../../Keycloak";
import {
  selectedApplication,
  selectedApplicationID,
} from "../../components/misc/globalState";
import { useSignals } from "@preact/signals-react/runtime";

const loading = signal<boolean>(false);
const error = signal<string>("");
const applicationName = signal<string>("");
const applications = signal<string[]>([]);

const DashboardPage = () => {
  useSignals();
  const { data, error, isLoading } = useKeycloakData();
  const initialLoad = useRef(true);
  useLoadingEffect(data, isLoading, isLoadingDone);

  effect(() => {
    if (data?.applicationsData) {
      const appNames = data.applicationsData.map((app: any) => app.app_name);
      applications.value = [...appNames];

      if (initialLoad.current) {
        selectedApplication.value =
          selectedApplication.value || appNames[0] || "";
        const selectedApp = data.applicationsData.find(
          (app: any) => app.app_name === selectedApplication.value
        );
        selectedApplicationID.value = selectedApp?.application_id || "";
        initialLoad.current = false;
      }
    }
  });

  const handleCreateApplication = async () => {
    loading.value = true;
    try {
      const token = keycloak.token;

      const response = await axios.post(
        `${import.meta.env.VITE_PUBLIC_URL}api/v1/private/applications`,
        { AppName: applicationName.value },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      applications.value.push(response.data.application.AppName);
      applicationName.value = "";
    } catch (err) {
      error.value = "Failed to create application.";
    } finally {
      loading.value = false;
    }
  };

  const handleDeleteApplication = async (appName: string) => {
    // Implement later
  };

  const handlePauseApplication = async (appName: string) => {
    // Implement later
  };

  if (isLoading) {
    return (
      <div className="fixed inset-0 flex items-center justify-center">
        <LoadingSpinner />
      </div>
    );
  }

  return (
    <div
      className="relative w-screen h-screen flex flex-wrap items-center justify-center"
      style={{ marginLeft: "64px", marginTop: "64px" }}
    >
      <LoadingSpinner />
      <div
        className={`transition-opacity duration-1000 ${
          isLoadingDone.value ? "opacity-100" : "opacity-0"
        } w-full max-w-4xl p-2 sm:p-4 md:p-8`}
      >
        <div className="flex flex-col items-center justify-center">
          <div className="min-w-max h-screen">
            <div className="mb-8">
              <h1 className="heading font-semibold text-center text-white text-4xl">
                Dashboard
              </h1>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="flex flex-col gap-y-3">
                  <input
                    type="text"
                    placeholder="Application Name"
                    className="input p-1 rounded-3xl text-center"
                    value={applicationName.value}
                    onChange={(e) => (applicationName.value = e.target.value)}
                  />
                  <CustomButton
                    text="Create Application"
                    onClick={handleCreateApplication}
                    bgColor="bg-sky-600"
                    hoverBgColor="bg-sky-500"
                    focusBgColor="bg-sky-500"
                  />
                </div>
                <div className="flex flex-col justify-start items-center">
                  <label
                    className="mb-2 font-semibold text-center"
                    htmlFor="applications"
                  >
                    Select Application
                  </label>
                  <select
                    id="applications"
                    className="input p-1 rounded-3xl w-8/12 text-center"
                    value={selectedApplication.value}
                    onChange={(e) => {
                      selectedApplication.value = e.target.value;
                      const selectedApp = data?.applicationsData.find(
                        (app: any) => app.app_name === e.target.value
                      );
                      selectedApplicationID.value =
                        selectedApp?.application_id ?? "";
                    }}
                  >
                    <option value="">-- Select Application --</option>
                    {applications.value.map((app) => (
                      <option key={app} value={app}>
                        {app}
                      </option>
                    ))}
                  </select>
                </div>
                {selectedApplication.value && (
                  <div className="col-span-2 flex flex-col justify-center items-center mt-4 gap-y-2">
                    <CustomButton
                      text="Pause Application"
                      onClick={() =>
                        handlePauseApplication(selectedApplication.value)
                      }
                      bgColor="bg-yellow-500"
                      hoverBgColor="bg-yellow-500"
                      focusBgColor="bg-yellow-500"
                    />
                    <CustomButton
                      text="Delete Application"
                      onClick={() =>
                        handleDeleteApplication(selectedApplication.value)
                      }
                      bgColor="bg-red-500"
                      hoverBgColor="bg-red-500"
                      focusBgColor="bg-red-500"
                    />
                  </div>
                )}
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
              <div className="bg-[#2A323C] p-4 rounded-lg shadow">
                <div className="text-lg font-semibold">...</div>
                <div className="text-2xl mt-2">...</div>
              </div>
              <div className="bg-[#2A323C] p-4 rounded-lg shadow">
                <div className="text-lg font-semibold">...</div>
                <div className="text-2xl mt-2">...</div>
              </div>
              <div className="bg-[#2A323C] p-4 rounded-lg shadow">
                <div className="text-lg font-semibold">...</div>
                <div className="text-2xl mt-2">...</div>
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="bg-[#2A323C] p-4 rounded-lg shadow">
                <h2 className="text-xl font-semibold">Active Users</h2>
                <div className="mt-4">
                  <canvas id="usersChart" className="w-full"></canvas>
                </div>
              </div>
              <div className="bg-[#2A323C] p-4 rounded-lg shadow">
                <h2 className="text-xl font-semibold">...</h2>
                <div className="mt-4">
                  <canvas id="salesChart" className="w-full"></canvas>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default DashboardPage;
