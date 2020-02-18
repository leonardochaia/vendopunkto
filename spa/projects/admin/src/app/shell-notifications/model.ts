import { Observable } from 'rxjs';

export interface ShellNotification {
    title: string;
    message?: string;
    type: 'message' | 'operation';
    date: number;
}

export interface OperationStartShellNotification extends ShellNotification {
    opId: string;
    type: 'operation';
}

export function isOperationNotification(n: ShellNotification): n is OperationStartShellNotification {
    return n.type === 'operation';
}