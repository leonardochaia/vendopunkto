import { Action, createReducer, on } from '@ngrx/store';
import * as ShellOperationsActions from './shell-operations.actions';
import { ShellOperationInstance, VPOperationAction } from '../model';

export const ShellOperationsFeatureKey = 'vpShellOperations';

export interface ShellOperationsState {
  operations: { [opId: string]: ShellOperationInstance };
  error: string;
}

export const initialState: ShellOperationsState = {
  operations: {
    '42962d8f-56fa-4b3e-9fed-9b57cbc067a5': {
      id: '42962d8f-56fa-4b3e-9fed-9b57cbc067a5',
      operation: {
        opKey: 'config-update',
        title: 'Update Configuration',
        description: 'Changes in configuration values',
        successAction: '[Config] Update config success',
        failureAction: '[Config] Update config failure'
      },
      status: 'pending'
    },
    'ok-56fa-4b3e-9fed-9b57cbc067a5': {
      id: 'ok-56fa-4b3e-9fed-9b57cbc067a5',
      operation: {
        opKey: 'config-update',
        title: 'Update Configuration',
        description: 'Changes in configuration values',
        successAction: '[Config] Update config success',
        failureAction: '[Config] Update config failure'
      },
      status: 'success'
    },
    'fail-56fa-4b3e-9fed-9b57cbc067a5': {
      id: 'fail-56fa-4b3e-9fed-9b57cbc067a5',
      operation: {
        opKey: 'config-update',
        title: 'Update Configuration',
        description: 'Changes in configuration values',
        successAction: '[Config] Update config success',
        failureAction: '[Config] Update config failure'
      },
      status: 'failure',
      error: 'algo salio muy malche'
    },
  },
  error: null
};

const ShellOperationsReducer = createReducer(
  initialState,
);

export function reducer(state: ShellOperationsState | undefined, action: Action) {

  let newState: ShellOperationsState;
  if (state) {
    newState = { ...state };
    const toDispatch = (action as VPOperationAction).vpDispatchOperationInstance;
    if (toDispatch) {
      newState = {
        ...state,
        operations: {
          ...state.operations,
          [toDispatch.id]: toDispatch,
        }
      };
    } else {
      for (const opId in state.operations) {
        if (state.operations.hasOwnProperty(opId)) {
          const instance = state.operations[opId];
          const op = instance.operation;
          const newOp: ShellOperationInstance = {
            ...newState.operations[opId],
          };
          switch (action.type) {
            case op.successAction:
              newOp.status = 'success';
              newState = {
                ...state,
                operations: {
                  ...state.operations,
                  [opId]: newOp
                }
              };
              break;
            case op.failureAction:
              newOp.status = 'failure';
              newOp.error = (action as any).error;
              newState = {
                ...state,
                operations: {
                  ...state.operations,
                  [opId]: newOp
                }
              };
              break;
          }
        }
      }
    }
  } else {
    newState = state;
  }

  return ShellOperationsReducer(newState, action);
}
