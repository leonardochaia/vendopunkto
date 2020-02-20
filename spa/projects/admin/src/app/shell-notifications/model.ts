import { generateUniqueId } from 'shared';
import { ShellOperationInstance } from '../shell-operations/model';

export type ShellNotificationType = 'message' | 'operation';

export interface ShellNotification {
    id: string;
    title: string;
    message?: string;
    type: ShellNotificationType;
    date: number;
}

export interface OperationStartShellNotification extends ShellNotification {
    opId: string;
    type: 'operation';
    operation?: ShellOperationInstance;
}

export function isOperationNotification(n: ShellNotification): n is OperationStartShellNotification {
    return n.type === 'operation';
}

export function createNotification(type: 'message'): ShellNotification;
export function createNotification(type: 'operation'): OperationStartShellNotification;
export function createNotification(type: ShellNotificationType): ShellNotification {
    return {
        id: generateUniqueId(),
        date: Date.now(),
        type,
        title: null
    };
}
