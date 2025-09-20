import { AuthService } from "@/services";
import showToast, { ToastStatus } from "@/utils/toastHelper";
import { useState } from "react";

export default function useUserOnboard() {
  const [loading, setLoading] = useState(false);

  const register = async (username: string, email: string, password: string, confirm: string) => {
    // normalize
    username = username.trim();
    email = email.trim();
    password = password.trim();
    confirm = confirm.trim();

    if (!username || !email || !password || !confirm) {
      showToast(ToastStatus.Failure, "All fields are required");
      return false;
    }

    if (password !== confirm) {
      showToast(ToastStatus.Failure, "Passwords do not match");
      return false;
    }

    if (password.length < 8) {
      showToast(ToastStatus.Failure, "Password must be at least 8 characters long");
    }

    setLoading(true);
    try {
      await AuthService.register({ username, email, password });
      return true;
    } catch (err: any) {
      return false;
    } finally {
      setLoading(false);
    }
  };

  return { loading, register };
}
