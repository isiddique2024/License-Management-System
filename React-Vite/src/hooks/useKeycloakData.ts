import useSWR from "swr";
import axios from "axios";
import keycloak, { initKeycloak, keycloakInitialized } from "../Keycloak";

const fetchKeycloakData = async () => {
  console.log("Fetching Keycloak data");

  await initKeycloak();

  if (!keycloakInitialized.value) {
    throw new Error("Keycloak initialization failed");
  }

  if (!keycloak.token) {
    return { userName: "User", applicationsData: [] };
  }

  try {
    const response = await axios.get(
      `${import.meta.env.VITE_PUBLIC_URL}api/v1/private/applications/data`,
      {
        headers: {
          Authorization: `Bearer ${keycloak.token}`,
        },
      }
    );

    const result = response.data;

    const user = keycloak.tokenParsed;
    const userName =
      user?.preferred_username || user?.name || user?.email || "User";
    console.log("Fetched Keycloak data");
    return { userName, applicationsData: result.applications };
  } catch (error) {
    console.error("Error fetching global data:", error);
    throw new Error("Error fetching global data");
  }
};

const useKeycloakData = () => {
  const { data, error } = useSWR("keycloakData", fetchKeycloakData, {
    revalidateOnFocus: false,
    revalidateOnReconnect: false,
    dedupingInterval: 60000,
  });

  return {
    data,
    error,
    isLoading: !error && !data,
  };
};

export default useKeycloakData;
