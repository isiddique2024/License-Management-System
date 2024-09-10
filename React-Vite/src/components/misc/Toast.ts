import { toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";

interface ToastProps {
  message: string;
  type: "success" | "error" | "info";
}
const showToast = ({ message, type }: ToastProps) => {
  const options = {
    autoClose: 5000,
    theme: "colored",
  };

  type === "success"
    ? toast.success(message, options)
    : toast.error(message, options);
};

export { showToast };
