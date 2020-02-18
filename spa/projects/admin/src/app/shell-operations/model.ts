import { Action } from '@ngrx/store';
import * as uuid from 'uuid/v4';

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
    const opId = uuid();

    return {
        id: opId,
        operation: op,
        status: 'pending'
    };
}


export interface VPOperationAction extends Action {
    vpDispatchOperationInstance: ShellOperationInstance;
}
