import showToast, { ToastStatus } from "@/utils/toastHelper";
import { create } from "zustand";
import fetchTimeout from "@/utils/fetchTimeOut";
import { InstanceType } from "@/utils/types";
import { io, Socket } from "socket.io-client";

interface InstanceData {
    [chall_id: number]: InstanceType;
};

interface InstanceStore {
    instances: InstanceData;
    socket: Socket | null;
    loading: boolean;
    fetchInstances: () => void;
    startInstance: (chall_id: number) => void;
    stopInstance: (chall_id: number) => void;
    extendInstance: (chall_id: number) => void;
    updateInstance: (chall_id: number, instance: Partial<InstanceType>) => void;
    setLoading: (valueToSet: boolean) => void;
};

export const useInstanceStore = create<InstanceStore>((set) => ({
    instances: {},
    loading: false,
    socket: null,

    setLoading: (valueToSet: boolean) => {
        set({ loading: valueToSet });
    },

    fetchInstances: async () => {
        if (useInstanceStore.getState().socket !== null) {
            return;
        };

        const socket = io("/", {
            path: "/ws",
        });

        socket.on("instanceUpdate", (payload) => {
            useInstanceStore.getState().updateInstance(payload.chall_id, {
                password: payload.password,
                port: payload.port,
                hostname: payload.hostname,
                deadline: payload.deadline,
                deployment: payload.deployment,
                active: true,
            });
        });

        socket.on("instanceStop", (payload) => {
            useInstanceStore.getState().updateInstance(payload.chall_id, { active: false });
        });

        socket.on("instances", (instances) => {
            instances.forEach((instance: InstanceType) => {
                useInstanceStore.getState().updateInstance(instance.chall_id, {
                    password: instance.password,
                    port: instance.port,
                    hostname: instance.hostname,
                    deadline: instance.deadline,
                    deployment: instance.deployment,
                    active: true,
                });
            });
        });

        set({ socket });
    },

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

                useInstanceStore.getState().updateInstance(chall_id, {
                    password: instanceJSON.message.password,
                    port: instanceJSON.message.port,
                    hostname: instanceJSON.message.hostname,
                    deadline: instanceJSON.message.deadline,
                    deployment: instanceJSON.message.deployment,
                    active: true,
                });

                showToast(ToastStatus.Success, "instance started successfully");
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
                useInstanceStore.getState().updateInstance(chall_id, { active: false });

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

                useInstanceStore.getState().updateInstance(chall_id, { deadline: response.message.deadline });

                showToast(ToastStatus.Success, "instance extended successfully");
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
