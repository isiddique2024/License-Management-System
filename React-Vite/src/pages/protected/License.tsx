import useKeycloakData from "../../hooks/useKeycloakData";
import LoadingSpinner, {
  isLoadingDone,
} from "../../components/misc/LoadingSpinner";
import useLoadingEffect from "../../hooks/useLoadingEffect";
import Table from "../../components/Table/Table";
import licenseColumn from "../../components/Table/ColumnHeaders";
import LicenseTableToolbar from "../../components/licenses/LicenseTableToolbar";
import { LicenseItem } from "../../components/licenses/utils";
import useFilteredData from "../../hooks/useFilteredData";
import { selectedApplication } from "../../components/misc/globalState";

interface Application {
  application_id: string;
  app_name: string;
  created_at: string;
  updated_at: string;
  licenses: LicenseItem[] | null;
}

const LicensesPage = () => {
  const { data, isLoading } = useKeycloakData();

  useLoadingEffect(data, isLoading, isLoadingDone);

  const licensesData: LicenseItem[] =
    data?.applicationsData?.flatMap((app: Application) => app.licenses ?? []) ??
    [];

  // Filter licenses based on selected application
  const filteredLicenses = selectedApplication.value
    ? licensesData.filter((license: LicenseItem) =>
        data?.applicationsData?.some(
          (app: Application) =>
            app.app_name === selectedApplication.value &&
            app.licenses?.some((lic: LicenseItem) => lic.key === license.key)
        )
      )
    : licensesData;

  const { filteredData, searchTerm } = useFilteredData(
    filteredLicenses,
    (license) => license.key
  );

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
        <h1 className="heading font-semibold text-center text-white text-4xl">
          Licenses
        </h1>
        <Table
          columns={licenseColumn}
          filteredData={filteredData}
          searchTerm={searchTerm}
          toolbar={<LicenseTableToolbar />}
        />
      </div>
    </div>
  );
};

export default LicensesPage;
