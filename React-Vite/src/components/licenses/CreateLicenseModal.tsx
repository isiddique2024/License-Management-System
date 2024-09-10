import React from "react";
import { handleGenerateLicense } from "./utils";

const CreateLicenseModal: React.FC = () => {
  return (
    <dialog id="licenseModal" className="modal">
      <div className="modal-box bg-[#333b45] rounded-3xl shadow-lg">
        <h3 className="heading text-white text-2xl mb-4">
          Create A New License
        </h3>
        <div className="w-full">
          <div className="form-group flex flex-col items-center">
            <label className="label text-white" htmlFor="licenseAmount">
              License Amount
            </label>
            <input
              type="number"
              id="licenseAmount"
              className="input p-1 rounded-3xl text-center"
              defaultValue="1"
              min="1"
              max="100"
            />
          </div>
          <div className="form-group flex flex-col items-center">
            <label className="label text-white" htmlFor="licenseMask">
              License Mask
            </label>
            <input
              type="text"
              id="licenseMask"
              className="input p-1 rounded-3xl text-center"
              defaultValue="XXXXXX-XXXXXX"
            />
          </div>
          <div className="form-group flex flex-col items-center">
            <label className="label text-white" htmlFor="keyPrefix">
              Prefix
            </label>
            <input
              type="text"
              id="keyPrefix"
              className="input p-1 rounded-3xl text-center"
              defaultValue="LICENSE"
            />
          </div>
          <div className="form-group flex flex-col items-center">
            <label className="label text-white" htmlFor="licenseNote">
              License Note
            </label>
            <input
              type="text"
              id="licenseNote"
              className="input p-1 rounded-3xl text-center"
            />
          </div>
          <div className="form-group flex flex-col items-center">
            <label className="label text-white" htmlFor="licenseExpiryUnit">
              License Expiry Unit
            </label>
            <select
              id="licenseExpiryUnit"
              className="input p-1 rounded-3xl w-5/12 text-center"
            >
              <option value="Day">Days</option>
              <option value="Week">Weeks</option>
              <option value="Month">Months</option>
              <option value="Year">Years</option>
            </select>
          </div>
          <div className="form-group flex flex-col items-center">
            <label className="label text-white" htmlFor="licenseDuration">
              License Duration
            </label>
            <input
              type="number"
              id="licenseDuration"
              className="input p-1 rounded-3xl text-center"
              defaultValue="1"
              min="1"
              max="10"
            />
          </div>
          <button
            onClick={handleGenerateLicense}
            className="button mt-4 p-2 rounded-3xl bg-blue-500 text-white w-full text-center"
          >
            Generate Keys
          </button>
        </div>
      </div>
      <form method="dialog" className="modal-backdrop">
        <button>Close</button>
      </form>
    </dialog>
  );
};

export default CreateLicenseModal;
