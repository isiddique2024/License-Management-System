import axios from "axios";
import keycloak from "../../Keycloak";
import { mutate } from "swr";
import { showToast } from "../misc/Toast";
import { signal } from "@preact/signals-react";
import { selectedApplicationID } from "../misc/globalState";

export const ITEMS_PER_PAGE = 10;
export const currentPage = signal(1);
export const selectedLicenses = signal<string[]>([]);

export interface LicenseItem {
  key: string;
  created_on: string;
  generated_by: string;
  duration: string;
  note?: string;
  used_on?: string;
  expires_on?: string;
  status?: string;
  ip?: string;
  hwid?: string;
}

export interface GenerateLicenseData {
  license_amount: number;
  license_mask: string;
  prefix: string;
  license_note: string;
  license_expiry_unit: string;
  license_duration: number;
}

const showSuccessToast = (message: string) => {
  showToast({ message, type: "success" });
};

const showErrorToast = (message: string) => {
  showToast({ message, type: "error" });
};

const mutateKeycloakData = async () => {
  await mutate("keycloakData");
};

const handleApiRequest = async (
  apiCall: () => Promise<any>,
  successMessage: string,
  errorMessage: string
) => {
  try {
    const response = await apiCall();
    showSuccessToast(successMessage);
    await mutateKeycloakData();
    return response;
  } catch (error) {
    console.error(errorMessage, error);
    showErrorToast(errorMessage);
    throw error;
  }
};

export const generateLicense = (data: GenerateLicenseData) => {
  return handleApiRequest(
    () =>
      axios.post(
        `${import.meta.env.VITE_PUBLIC_URL}api/v1/private/applications/${
          selectedApplicationID.value
        }/licenses`,
        data,
        { headers: { Authorization: `Bearer ${keycloak.token}` } }
      ),
    "License(s) generated",
    "Failed to generate license(s)"
  );
};

export const handleDeleteAllLicenses = () => {
  return handleApiRequest(
    () =>
      axios.delete(
        `${import.meta.env.VITE_PUBLIC_URL}api/v1/private/applications/${
          selectedApplicationID.value
        }/licenses-all`,
        {
          headers: {
            Accept: "application/json",
            Authorization: `Bearer ${keycloak.token}`,
          },
        }
      ),
    "All licenses deleted successfully",
    "Failed to delete licenses"
  );
};

export const handleDeleteSelectedLicenses = () => {
  return handleApiRequest(
    () =>
      axios.delete(
        `${import.meta.env.VITE_PUBLIC_URL}api/v1/private/applications/${
          selectedApplicationID.value
        }/licenses`,
        {
          data: { keys: selectedLicenses.value },
          headers: {
            Accept: "application/json",
            Authorization: `Bearer ${keycloak.token}`,
          },
        }
      ),
    "Selected licenses deleted successfully",
    "Failed to delete selected licenses"
  );
};

export const handleDeleteLicense = (key: string) => {
  return handleApiRequest(
    () =>
      axios.delete(
        `${import.meta.env.VITE_PUBLIC_URL}api/v1/private/applications/${
          selectedApplicationID.value
        }/licenses/${key}`,
        {
          headers: {
            Accept: "application/json",
            Authorization: `Bearer ${keycloak.token}`,
          },
        }
      ),
    "License deleted successfully",
    "Failed to delete license"
  );
};

export const handleBanLicense = (key: string) => {
  return handleApiRequest(
    () =>
      axios.patch(
        `${import.meta.env.VITE_PUBLIC_URL}api/v1/private/applications/${
          selectedApplicationID.value
        }/licenses/${key}/ban`,
        { key },
        {
          headers: {
            Accept: "application/json",
            Authorization: `Bearer ${keycloak.token}`,
          },
        }
      ),
    "License banned successfully",
    "Failed to ban license"
  );
};

export const handleSelectLicense = (key: string, selected: boolean) => {
  selectedLicenses.value = selected
    ? [...selectedLicenses.value, key]
    : selectedLicenses.value.filter((item) => item !== key);
};

export const handleSelectAll = (
  selected: boolean,
  currentData: LicenseItem[]
) => {
  if (selected) {
    const currentKeys = currentData.map((item) => item.key);
    selectedLicenses.value = [
      ...new Set([...selectedLicenses.value, ...currentKeys]),
    ];
  } else {
    const currentKeys = currentData.map((item) => item.key);
    selectedLicenses.value = selectedLicenses.value.filter(
      (key) => !currentKeys.includes(key)
    );
  }
};

export const handleGenerateLicense = async () => {
  const licenseAmount = parseInt(
    (document.getElementById("licenseAmount") as HTMLInputElement).value
  );
  const licenseMask = (
    document.getElementById("licenseMask") as HTMLInputElement
  ).value;
  const prefix = (document.getElementById("keyPrefix") as HTMLInputElement)
    .value;
  const licenseNote = (
    document.getElementById("licenseNote") as HTMLInputElement
  ).value;
  const licenseExpiryUnit = (
    document.getElementById("licenseExpiryUnit") as HTMLSelectElement
  ).value;
  const licenseDuration = parseInt(
    (document.getElementById("licenseDuration") as HTMLInputElement).value
  );

  const data: GenerateLicenseData = {
    license_amount: licenseAmount,
    license_mask: licenseMask,
    prefix: prefix,
    license_note: licenseNote,
    license_expiry_unit: licenseExpiryUnit,
    license_duration: licenseDuration,
  };

  await generateLicense(data);
};
