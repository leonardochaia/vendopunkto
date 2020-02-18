import { Action } from '@ngrx/store';
import { generateUniqueId } from 'shared';

export interface ShellOperation {
    opKey: string;
    title: string;
    description: string;

    successAction: string;
    failureAction: string;
}

export interface ShellOperationInstance<TOperation extends ShellOperation = ShellOperation> {
    id: string;
    title?: string;
    description?: string;

    operation: TOperation;
    status: 'pending' | 'success' | 'failure';
    error?: unknown;
}

export function createOperationInstance<T extends ShellOperation>(op: T): ShellOperationInstance<T> {
    return {
        id: generateUniqueId(),
        operation: op,
        status: 'pending'
    };
}


export interface VPOperationAction extends Action {
    vpDispatchOperationInstance: ShellOperationInstance;
}
