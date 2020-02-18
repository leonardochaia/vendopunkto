import { createFeatureSelector, createSelector } from '@ngrx/store';
import * as fromShellNotifications from './shell-notifications.reducer';
import { selectOperations } from '../../shell-operations/+state/shell-operations.selectors';
import { isOperationNotification, OperationStartShellNotification } from '../model';

export const selectConfigState = createFeatureSelector<fromShellNotifications.ShellNotificationsState>(
  fromShellNotifications.ShellNotificationsFeatureKey
);

export const selectNotifications = createSelector(selectConfigState, s => s.notifications);

export const selectHasPendingNotifications =
  createSelector(selectConfigState,
    selectOperations,
    (s, ops) => s.notifications.filter(n => isOperationNotification(n))
      .some((n: OperationStartShellNotification) => ops[n.opId].status === 'pending'));

export const selectPopoverNotification = createSelector(selectConfigState,
  s => s.popoverNotification);
