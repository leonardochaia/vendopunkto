import { Injectable } from '@angular/core';
import { Store, Action } from '@ngrx/store';
import { selectOperations } from './shell-operations.selectors';
import { ShellOperationsState } from './shell-operations.reducer';
import { ShellOperation, createOperationInstance, ShellOperationInstance, VPOperationAction } from '../model';
import { OperationStartShellNotification } from '../../shell-notifications/model';
import { notificationAdded } from '../../shell-notifications/+state/shell-notifications.actions';

@Injectable({
  providedIn: 'root'
})
export class ShellOperationsFacade {

  public readonly operations$ = this.store.select(selectOperations);

  constructor(private readonly store: Store<ShellOperationsState>) { }

  public dispatch(action: Action) {
    this.store.dispatch(action);
  }

  public dispatchOperation(
    action: Action,
    instance: ShellOperationInstance,
    notification: Partial<OperationStartShellNotification> = null) {

    (action as VPOperationAction).vpDispatchOperationInstance = instance;
    this.store.dispatch(action);

    if (notification) {
      const notif: OperationStartShellNotification = {
        ...notification,
        opId: instance.id,
        type: 'operation',
        title: notification.title || instance.title || instance.operation.title,
        date: Date.now(),
      };

      this.store.dispatch(notificationAdded({ notification: notif }));
    }

    return instance;
  }
}
