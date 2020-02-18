import { Action, createReducer, on } from '@ngrx/store';
import * as ShellOperationsActions from './shell-operations.actions';
import { ShellOperationInstance, VPOperationAction } from '../model';

export const ShellOperationsFeatureKey = 'vpShellOperations';

export interface ShellOperationsState {
  operations: { [opId: string]: ShellOperationInstance };
  error: string;
}

export const initialState: ShellOperationsState = {
  operations: {},
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
