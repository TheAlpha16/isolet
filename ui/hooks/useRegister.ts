import fetchTimeout from '@/utils/fetchTimeOut';
import showToast, { ToastStatus } from '@/utils/toastHelper';
import { useState } from 'react';

function useRegister() {
    const [loading, setLoading] = useState(false);

    const registerAPI = async (username: string, email: string, password: string, confirm: string) => {
        setLoading(true);

        username = username.trim();
        email = email.trim();
        password = password.trim();
        confirm = confirm.trim();

        if (!username || !email || !password || !confirm) {
            showToast(ToastStatus.Failure, "all fields are required");
            setLoading(false);
            return false;
        }

        if (password !== confirm) {
            showToast(ToastStatus.Failure, "passwords mismatch");
            setLoading(false);
            return false;
        }

        if (!(password.length >= 8)) {
            showToast(ToastStatus.Failure, "password length should be 8 or more");
            setLoading(false);
            return false;
        }

        let formData = new FormData();
        formData.append("username", username);
        formData.append("email", email);
        formData.append("password", password);
        formData.append("confirm", confirm);

        try {
            const res = await fetchTimeout('/auth/register', 60000, new AbortController().signal, {
                method: 'POST',
                body: formData
            });

            if (res.ok) {
                const jsonResp = await res.json();
                showToast(ToastStatus.Success, jsonResp.message);
                return true;
            } else {
                const resJSON = await res.json();
                showToast(ToastStatus.Failure, resJSON.message);
            }
        } catch (error: any) {
            if (error.name === 'AbortError') {
                showToast(ToastStatus.Failure, 'request timed out, reload!');
            } else {
                showToast(ToastStatus.Warning, 'server seems offline');
            }
        } finally {
            setLoading(false);
        }

        return false;
    };

    return { loading, registerAPI };
}

export default useRegister;