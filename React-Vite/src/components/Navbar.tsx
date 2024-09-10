import { signal, effect } from "@preact/signals-react";
import { FaSignOutAlt } from "react-icons/fa";
import useKeycloakData from "../hooks/useKeycloakData";
import keycloak from "../Keycloak";

const initialized = signal(false);
const visible = signal(false);

interface NavbarIconProps {
  icon: React.ReactNode;
  text: string;
  delay: number;
  onClick: () => void;
}

const NavbarIcon: React.FC<NavbarIconProps> = ({
  icon,
  text,
  delay,
  onClick,
}) => (
  <button
    className={`navbar-icon flex items-center px-4 py-2 group hover:shadow-xl hover:text-white rounded-3xl transition-all duration-500 ease-linear cursor-pointer shadow-lg hover:bg-red-500 focus:bg-red-500 ${
      visible.value ? "show" : ""
    }`}
    style={{ transitionDelay: `${delay}ms` }}
    onClick={onClick}
    onKeyDown={(e) => {
      if (e.key === "Enter" || e.key === " ") {
        onClick();
      }
    }}
  >
    {icon}
    <span className="ml-2 md:block font-bold text-white group-hover:scale-100 hover:shadow-md">
      {text}
    </span>
  </button>
);

const Navbar: React.FC = () => {
  const { data, isLoading } = useKeycloakData();

  effect(() => {
    if (data && !isLoading) {
      initialized.value = true;
      setTimeout(() => {
        visible.value = true;
      }, 1000);
    }
  });

  return (
    <div className="fixed top-0 left-0 right-0 h-16 flex justify-end items-center bg-[#2A323C] shadow-lg z-10">
      {isLoading || !initialized.value ? (
        <div className="flex items-center justify-end w-full h-full mr-10">
          {/* <span className="loading loading-bars"></span> */}
        </div>
      ) : (
        <>
          <div
            className={`welcome-message flex items-center px-4 py-2 ${
              visible.value ? "show" : ""
            }`}
            style={{
              transition:
                "opacity 1s ease-in-out, transform 1s ease-in-out, color 1s ease-in-out",
              transitionDelay: "200ms",
            }}
          >
            <span className="font-bold">Welcome {data?.userName}</span>
          </div>
          <NavbarIcon
            icon={<FaSignOutAlt size="25" />}
            text="Logout"
            delay={400}
            onClick={() => keycloak.logout()}
          />
        </>
      )}
    </div>
  );
};

export default Navbar;
