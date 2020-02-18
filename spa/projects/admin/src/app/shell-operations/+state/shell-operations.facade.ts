import { Injectable } from '@angular/core';
import { Store, Action } from '@ngrx/store';
import { selectOperations } from './shell-operations.selectors';
import { ShellOperationsState } from './shell-operations.reducer';
import { ShellOperation, createOperationInstance, ShellOperationInstance, VPOperationAction } from '../model';
import { OperationStartShellNotification, createNotification } from '../../shell-notifications/model';
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
    notify = false) {

    (action as VPOperationAction).vpDispatchOperationInstance = instance;
    this.store.dispatch(action);

    if (notify) {
      const notification: OperationStartShellNotification = {
        ...createNotification('operation'),
        opId: instance.id,
        title: instance.title || instance.operation.title,
      };

      this.store.dispatch(notificationAdded({ notification }));
    }

    return instance;
  }
}
