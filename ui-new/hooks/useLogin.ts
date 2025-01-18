import fetchTimeout from '@/utils/fetchTimeOut';
import showToast, { ToastStatus } from '@/utils/toastHelper';
import { useAuthStore } from '@/store/authStore';
import { useState } from 'react';

function useLogin() {
    const [loading, setLoading] = useState(false);
    const { setUser } = useAuthStore();

    const loginAPI = async (email: string, password: string) => {
        setLoading(true);

        email = email.trim();
        password = password.trim();

        if (!email || !password) {
            showToast(ToastStatus.Failure, 'email and password are required');
            setLoading(false);
            return false;
        }

        let formData = new FormData();
        formData.append('email', email);
        formData.append('password', password);

        try {
            const res = await fetchTimeout('/auth/login', 5000, new AbortController().signal, {
                method: 'POST',
                body: formData
            });

            if (res.ok) {
                const user = await res.json();
                setUser(user);
                showToast(ToastStatus.Success, `Welcome back ${user.username}!`);
                return true;
            } else {
                const resJSON = await res.json();
                showToast(ToastStatus.Failure, resJSON.message);
            }
        } catch (error: any) {
            if (error.name === 'AbortError') {
                showToast(ToastStatus.Failure, 'verification timed out, reload!');
            } else {
                showToast(ToastStatus.Warning, 'server seems offline');
            }
        } finally {
            setLoading(false);
        }

        return false;
    };

    const forgotPasswordAPI = async (email: string) => {
        setLoading(true);

        email = email.trim();

        if (!email) {
            showToast(ToastStatus.Failure, 'email is required');
            setLoading(false);
            return false;
        }

        let formData = new FormData();
        formData.append('email', email);

        try {
            const res = await fetchTimeout('/auth/forgot-password', 5000, new AbortController().signal, {
                method: 'POST',
                body: formData
            });

            if (res.ok) {
                const resJSON = await res.json();
                showToast(ToastStatus.Success, resJSON.message);
                return true;
            } else {
                const resJSON = await res.json();
                showToast(ToastStatus.Failure, resJSON.message);
            }
        } catch (error: any) {
            if (error.name === 'AbortError') {
                showToast(ToastStatus.Failure, 'verification timed out, reload!');
            } else {
                showToast(ToastStatus.Warning, 'server seems offline');
            }
        } finally {
            setLoading(false);
        }

        return false;
    }

    return { loading, loginAPI, forgotPasswordAPI };
}

export default useLogin;