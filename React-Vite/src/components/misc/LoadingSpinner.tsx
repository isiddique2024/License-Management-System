import React from "react";
import { signal } from "@preact/signals-react";

const isLoadingDone = signal<boolean>(false);

const LoadingSpinner = () => {
  return (
    <div
      className={`absolute transition-opacity duration-1000 ${
        isLoadingDone.value ? "opacity-0" : "opacity-100"
      } loading-spinner loading justify-center items-center w-16 h-16 sm:w-24 sm:h-24 md:w-32 md:h-32 bg-[#60A5FA] z-10 mb-10 sm:mb-16 md:mb-20`}
    ></div>
  );
};

export default LoadingSpinner;
export { isLoadingDone };
