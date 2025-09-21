import { AuthService } from "@/services";
import showToast, { ToastStatus } from "@/utils/toastHelper";
import { useState } from "react";

export default function useLogin() {
  const [loading, setLoading] = useState(false);

  const login = async (identifier: string, password: string) => {
    // normalize
    identifier = identifier.trim();
    password = password.trim();

    if (!identifier || !password) {
      showToast(ToastStatus.Failure, "email/username and password are required");
      return false;
    }

    setLoading(true);
    try {
      const session = await AuthService.login({ identifier, password });
      if (!session) {
        return false;
      }
      return true;
    } catch (err: any) {
      return false;
    } finally {
      setLoading(false);
    }
  };

  return { login, loading };
}
