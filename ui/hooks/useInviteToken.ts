import fetchTimeout from "@/utils/fetchTimeOut";
import { useState } from "react";
import showToast, { ToastStatus } from "@/utils/toastHelper";

function useInvite() {
    const [inviteToken, setInviteToken] = useState("");
    const [loading, setLoading] = useState(false);

    const generateInviteTokenAPI = async () => {
        setLoading(true); 

        try {
            const response = await fetchTimeout(
                "api/profile/team/invite",
                60000,
                new AbortController().signal,
                {},
            );

            if (response.ok) {
                const { token } = await response.json();
                setInviteToken(token);
            }
        } catch (error: any) {
            if (error.name === "AbortError") {
                showToast(
                    ToastStatus.Failure,
                    "verification timed out, reload!"
                );
            } else {
                showToast(ToastStatus.Warning, "server seems offline");
            }
        } finally {
            setLoading(false);
        }
    }

    return { inviteToken, loading, generateInviteTokenAPI };
}

export default useInvite;