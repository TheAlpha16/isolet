'use client'

import { toast } from 'react-toastify';

export enum ToastStatus {
    Success,
    Warning,
    Failure
}

function showToast(status: ToastStatus, message: string) {
    switch (status) {
        case ToastStatus.Success:
            toast.success(message, { containerId: 'notification-toast' })
            break;
        case ToastStatus.Failure:
            toast.error(message, { containerId: 'notification-toast' })
            break;
        case ToastStatus.Warning:
            toast.warn(message, { containerId: 'notification-toast' })
            break;
        default:
            toast.warn(message, { containerId: 'notification-toast' })
    }
}

export default showToast;