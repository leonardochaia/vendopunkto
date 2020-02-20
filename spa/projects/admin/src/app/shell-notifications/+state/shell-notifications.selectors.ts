import { createFeatureSelector, createSelector } from '@ngrx/store';
import * as fromShellNotifications from './shell-notifications.reducer';
import { selectOperations } from '../../shell-operations/+state/shell-operations.selectors';
import { isOperationNotification, OperationStartShellNotification, ShellNotification } from '../model';
import { ShellOperationInstance } from '../../shell-operations/model';

export const selectConfigState = createFeatureSelector<fromShellNotifications.ShellNotificationsState>(
  fromShellNotifications.ShellNotificationsFeatureKey
);

function appendOperationIfNeeded(n: ShellNotification, ops: { [k: string]: ShellOperationInstance }) {
  if (isOperationNotification(n)) {
    return {
      ...n,
      operation: ops[n.opId]
    };
  }
  return n;
}

export const selectNotifications = createSelector(
  selectConfigState,
  selectOperations,
  (s, ops) => s.notifications.map(n => appendOperationIfNeeded(n, ops)));

export const selectPopoverNotification = createSelector(
  selectConfigState,
  selectOperations,
  (s, ops) => s.popoverNotification ? appendOperationIfNeeded(s.popoverNotification, ops) : null);

export const selectHasPendingNotifications =
  createSelector(selectConfigState,
    selectOperations,
    (s, ops) => s.notifications.filter(n => isOperationNotification(n))
      .some((n: OperationStartShellNotification) => ops[n.opId].status === 'pending'));

