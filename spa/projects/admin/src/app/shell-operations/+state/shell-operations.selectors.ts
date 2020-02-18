import { createFeatureSelector, createSelector } from '@ngrx/store';
import * as fromShellOperations from './shell-operations.reducer';

export const selectShellOperations = createFeatureSelector<fromShellOperations.ShellOperationsState>(
  fromShellOperations.ShellOperationsFeatureKey
);

export const selectOperations = createSelector(selectShellOperations, s => s.operations);

