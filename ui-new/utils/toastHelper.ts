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
            toast.success(message)
            break;
        case ToastStatus.Failure:
            toast.error(message)
            break;
        case ToastStatus.Warning:
            toast.warn(message)
            break;
        default:
            toast.warn(message)
    }
}

export default showToast;