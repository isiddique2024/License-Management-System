import React from "react";

interface CustomButtonProps {
  type?: "button" | "submit" | "reset";
  text: string;
  onClick: () => void;
  bgColor: string;
  hoverBgColor: string;
  focusBgColor: string;
}

const CustomButton: React.FC<CustomButtonProps> = ({
  type = "button",
  text,
  onClick,
  bgColor,
  hoverBgColor,
  focusBgColor,
}) => {
  return (
    <button
      type={type}
      onClick={onClick}
      className={`button text-center justify-center flex-shrink font-semibold text-sm overflow-clip outline-transparent ${bgColor} hover:shadow-xl text-white rounded-3xl transition-all duration-300 ease-linear cursor-pointer shadow-lg hover:${hoverBgColor} focus:${focusBgColor}`}
    >
      {text}
    </button>
  );
};

export default CustomButton;
