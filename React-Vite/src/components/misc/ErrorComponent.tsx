import React from "react";

interface ErrorComponentProps {
  error: any;
}

const ErrorComponent: React.FC<ErrorComponentProps> = ({ error }) => {
  return (
    <div className="text-red-500 text-center mt-4">
      <p>Failed to load data: {error.message}</p>
    </div>
  );
};

export default ErrorComponent;
