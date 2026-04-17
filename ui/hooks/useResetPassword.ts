import { AuthService } from "@/services";
import showToast, { ToastStatus } from "@/utils/toastHelper";
import { useState } from "react";

export default function useResetPassword() {
  const [loading, setLoading] = useState(false);

  const resetPassword = async (password: string, confirmPassword: string, token: string) => {
    // normalize
    password = password.trim();
    confirmPassword = confirmPassword.trim();
    token = token.trim();

    if (!token) {
      showToast(ToastStatus.Failure, "invalid token");
    }

    if (!password || !confirmPassword) {
      showToast(ToastStatus.Failure, "all fields are required");
      return false;
    }

    if (password !== confirmPassword) {
      showToast(ToastStatus.Failure, "passwords do not match");
      return false;
    }

    setLoading(true);
    try {
      await AuthService.resetPassword({ password, token });
      return true;
    } catch {
      return false;
    } finally {
      setLoading(false);
    }
  };

  return { resetPassword, loading };
}
