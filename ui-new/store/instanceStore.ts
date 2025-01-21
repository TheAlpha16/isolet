import showToast, { ToastStatus } from "@/utils/toastHelper";
import { create } from "zustand";
import fetchTimeout from "@/utils/fetchTimeOut";

export interface InstanceType {
    chall_id: number;
    password: string;
    port: number;
    hostname: string;
    deadline: number;
    deployment: string;
    connString: string;
    active: boolean;
};

export interface InstanceData {
    [chall_id: number]: InstanceType;
};

interface InstanceStore {
    instances: InstanceData;
    loading: boolean;
    fetchInstances: () => void;
    startInstance: (chall_id: number) => void;
    stopInstance: (chall_id: number) => void;
    extendInstance: (chall_id: number) => void;
    updateInstance: (chall_id: number, instance: Partial<InstanceType>) => void;
};

export const useInstanceStore = create<InstanceStore>((set) => ({
    instances: {},
    loading: false,

    fetchInstances: async () => {},

    startInstance: async (chall_id: number) => {
        set({ loading: true });
        
        let formData = new FormData();
        formData.append("chall_id", chall_id.toString());

        try {
            const res = await fetch(`/api/launch`, {
                method: "POST",
                body: formData,
            });

            if (res.ok) {
                const instanceJSON = await res.json();

                set((state) => {
                    const instance = instanceJSON as InstanceType;

                    return {
                        instances: { ...state.instances, [chall_id]: instance }
                    };
                })

            } else if (res.status === 401) {
                showToast(ToastStatus.Warning, "login to continue");
            } else {
                const response = await res.json();
                showToast(ToastStatus.Failure, response.message);
            }
        } catch (error: any) {
            if (error.name === "AbortError") {
                showToast(ToastStatus.Failure, "request timed out!");
            } else {
                showToast(ToastStatus.Warning, "seems offline");
            }
        } finally {
            set({ loading: false });
        }
    },

    stopInstance: async (chall_id: number) => {
        set({ loading: true });

        let formData = new FormData();
        formData.append("chall_id", chall_id.toString());

        try {
            const res = await fetch(`/api/stop`, {
                method: "POST",
                body: formData,
            })

            if (res.ok) {
                const response = await res.json();
                set((state) => {
                    const updatedInstances = { ...state.instances };
                    
                    updatedInstances[chall_id].active = false;

                    return { instances: updatedInstances };
                });

                showToast(ToastStatus.Success, response.message);
            } else if (res.status === 401) {
                showToast(ToastStatus.Warning, "login to continue");
            } else {
                const response = await res.json();
                showToast(ToastStatus.Failure, response.message);
            }
        } catch (error: any) {
            if (error.name === "AbortError") {
                showToast(ToastStatus.Failure, "request timed out!");
            } else {
                showToast(ToastStatus.Warning, "seems offline");
            }
        } finally {
            set({ loading: false });
        }
    },

    extendInstance: async (chall_id: number) => {
        set({ loading: true });

        let formData = new FormData();
        formData.append("chall_id", chall_id.toString());

        try {
            const res = await fetchTimeout("/api/extend", 7000, new AbortController().signal, {
                method: "POST",
                body: formData,
            });

            if (res.ok) {
                const response = await res.json();
                set((state) => {
                    const updatedInstances = { ...state.instances };

                    updatedInstances[chall_id].deadline = response.deadline;

                    return { instances: updatedInstances };
                });

                showToast(ToastStatus.Success, response.message);
            } else if (res.status === 401) {
                showToast(ToastStatus.Warning, "login to continue");
            } else {
                const response = await res.json();
                showToast(ToastStatus.Failure, response.message);
            }
        } catch (error: any) {
            if (error.name === "AbortError") {
                showToast(ToastStatus.Failure, "request timed out!");
            } else {
                showToast(ToastStatus.Warning, "seems offline");
            }
        } finally {
            set({ loading: false });
        }
    },

    updateInstance: (chall_id: number, instance: Partial<InstanceType>) => {
        set((state) => {
            const updatedInstances = { ...state.instances };

            updatedInstances[chall_id] = { ...updatedInstances[chall_id], ...instance };

            return { instances: updatedInstances };
        });
    }
}));