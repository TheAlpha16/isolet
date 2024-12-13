import fetchTimeout from '@/utils/fetchTimeOut';
import showToast, { ToastStatus } from '@/utils/toastHelper';
import { useAuthStore } from '@/store/authStore';
import { useState } from 'react';

function useJoinTeam() {
    const [loading, setLoading] = useState(false);
    const { user, setUser, logout } = useAuthStore();

    const teamJoin = async (teamname: string, password: string, action: string) => {
        setLoading(true);

        teamname = teamname.trim();
        password = password.trim();
        action = action.trim();

        if (!teamname || !password) {
            showToast(ToastStatus.Failure, 'teamname and password are required');
            setLoading(false);
            return;
        }

        if (action !== "join" && action !== "create") {
            showToast(ToastStatus.Failure, 'invalid action!');
            setLoading(false);
            return;
        }

        let url = `/onboard/team/${action}`;
        let formData = new FormData();
        formData.append('teamname', teamname);
        formData.append('password', password);

        try {
            const res = await fetchTimeout(url, 5000, new AbortController().signal, {
                method: 'POST',
                body: formData
            });

            switch (res.status) {
                case 200:
                    const teamIDJSON = await res.json();
                    let teamid = teamIDJSON.teamid;
                    let userCopy = { 
                        ...user, 
                        teamid, 
                        teamname, 
                        userid: user?.userid ?? -1, 
                        email: user?.email ?? "",
                        username: user?.username ?? "",
                        rank: user?.rank ?? 3
                    };
                    setUser(userCopy);
                    showToast(ToastStatus.Success, `${action}${action == "create" ? "": "e"}d ${teamname}`);
                    return;
                case 401:
                    logout();
                    showToast(ToastStatus.Warning, "please login to continue");
                    break;
                default:
                    const resJSON = await res.json();
                    showToast(ToastStatus.Failure, resJSON.message);
                    break;
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

        return;
    };

    return { loading, teamJoin };
}

export default useJoinTeam;