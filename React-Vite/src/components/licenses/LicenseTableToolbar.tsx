import React from "react";
import CustomButton from "../licenses/CustomButton";
import CreateLicenseModal from "../licenses/CreateLicenseModal";
import {
  handleDeleteSelectedLicenses,
  handleDeleteAllLicenses,
} from "../licenses/utils";

const LicenseTableToolbar: React.FC = () => {
  return (
    <div className="flex flex-col md:flex-row justify-center gap-2 md:gap-4 mb-5 p-2 md:p-4">
      <CreateLicenseModal />
      <CustomButton
        text="Create License"
        onClick={() => {
          const modal = document.getElementById(
            "licenseModal"
          ) as HTMLDialogElement;
          modal?.showModal();
        }}
        bgColor="bg-sky-600"
        hoverBgColor="bg-sky-500"
        focusBgColor="bg-sky-500"
      />
      <CustomButton
        text="Delete Selected License(s)"
        onClick={handleDeleteSelectedLicenses}
        bgColor="bg-red-700"
        hoverBgColor="bg-red-500"
        focusBgColor="bg-red-500"
      />
      <CustomButton
        text="Delete All Licenses"
        onClick={handleDeleteAllLicenses}
        bgColor="bg-red-700"
        hoverBgColor="bg-red-500"
        focusBgColor="bg-red-500"
      />
    </div>
  );
};

export default LicenseTableToolbar;
