import fetchTimeout from '@/utils/fetchTimeOut';
import showToast, { ToastStatus } from '@/utils/toastHelper';
import { useAuthStore } from '@/store/authStore';
import { useState } from 'react';

function useLogin() {
    const [loading, setLoading] = useState(false);
    const { setLoggedIn, setUser } = useAuthStore();

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
                setLoggedIn(true);
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

    return { loading, loginAPI };
}

export default useLogin;