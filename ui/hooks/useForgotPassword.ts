import { AuthService } from "@/services";
import showToast, { ToastStatus } from "@/utils/toastHelper";
import { useState } from "react";

export default function useForgotPassword() {
  const [loading, setLoading] = useState(false);

  const forgotPassword = async (email: string) => {
    // normalize
    email = email.trim();

    if (!email) {
      showToast(ToastStatus.Failure, "email is required");
      return false;
    }

    setLoading(true);
    try {
      await AuthService.forgotPassword({ email });
      return true;
    } catch {
      return false;
    } finally {
      setLoading(false);
    }
  };

  return { forgotPassword, loading };
}
