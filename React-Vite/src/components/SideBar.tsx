import React from "react";
import { Link } from "react-router-dom";
import { MdSpaceDashboard } from "react-icons/md";
import { FaKey } from "react-icons/fa";
import { IoSettingsSharp, IoCloudDownloadSharp } from "react-icons/io5";
import { ImBlocked } from "react-icons/im";
import { RiInboxArchiveFill } from "react-icons/ri";
import { IoIosTime } from "react-icons/io";
import { signal, effect } from "@preact/signals-react";
import useKeycloakData from "../hooks/useKeycloakData";

interface SideBarIconProps {
  icon: React.ReactNode;
  text: string;
  delay: number;
}

const SideBarIcon: React.FC<SideBarIconProps> = ({ icon, text, delay }) => (
  <div
    className="sidebar-icon group hover:shadow-xl"
    style={{ transitionDelay: `${delay}ms` }}
  >
    {icon}
    <span className="sidebar-tooltip group-hover:scale-100 hover:shadow-md">
      {text}
    </span>
  </div>
);

const showContent = signal<boolean>(false);

const SideBar: React.FC = () => {
  const { data, isLoading } = useKeycloakData();
  const initialized = Boolean(data);

  effect(() => {
    if (initialized && !isLoading) {
      setTimeout(() => {
        showContent.value = true;
      }, 1000);
    }
  });

  effect(() => {
    if (showContent.value) {
      const icons = document.querySelectorAll(".sidebar-icon");
      icons.forEach((icon, index) => {
        setTimeout(() => {
          icon.classList.add("show");
        }, index * 100);
      });
    }
  });

  if (isLoading) {
    return (
      <div className="fixed top-0 left-0 h-full w-16 flex items-center justify-center bg-[#2A323C] shadow-lg z-10">
        {/* <span className="loading loading-spinner"></span> */}
      </div>
    );
  }

  const icons = [
    {
      to: "/dashboard",
      icon: <MdSpaceDashboard size="25" />,
      text: "Dashboard",
    },
    { to: "/dashboard/licenses", icon: <FaKey size="25" />, text: "Licenses" },
    {
      to: "/dashboard/sessions",
      icon: <IoIosTime size="25" />,
      text: "Sessions",
    },
    // { to: "/dashboard/users", icon: <FaUser size="25" />, text: "Users" },
    {
      to: "/dashboard/downloads",
      icon: <IoCloudDownloadSharp size="25" />,
      text: "Downloads",
    },
    {
      to: "/dashboard/logs",
      icon: <RiInboxArchiveFill size="25" />,
      text: "Logs",
    },
    {
      to: "/dashboard/blacklists",
      icon: <ImBlocked size="25" />,
      text: "Blacklists",
    },
    {
      to: "/dashboard/settings",
      icon: <IoSettingsSharp size="25" />,
      text: "Settings",
    },
  ];

  return (
    <div className="fixed top-0 left-0 h-full w-16 flex-auto flex-col bg-[#2A323C] shadow-lg z-10">
      {initialized &&
        icons.map((icon, index) => (
          <Link to={icon.to} key={icon.text}>
            <SideBarIcon
              icon={icon.icon}
              text={icon.text}
              delay={index * 100}
            />
          </Link>
        ))}
    </div>
  );
};

export default SideBar;
