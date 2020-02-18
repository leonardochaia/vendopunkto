import { Injectable } from '@angular/core';
import { ShellNotificationsState } from './shell-notifications.reducer';
import { Store, Action } from '@ngrx/store';
import { selectNotifications, selectHasPendingNotifications, selectPopoverNotification } from './shell-notifications.selectors';
import { ShellNotification } from '../model';
import * as shellNotificationsActions from './shell-notifications.actions';

@Injectable({
  providedIn: 'root'
})
export class ShellNotificationsFacade {

  public readonly notifications$ = this.store.select(selectNotifications);

  public readonly popoverNotification$ = this.store.select(selectPopoverNotification);

  public readonly hasPendingNotifications$ = this.store.select(selectHasPendingNotifications);

  constructor(private readonly store: Store<ShellNotificationsState>) { }

  public dispatch(action: Action) {
    this.store.dispatch(action);
  }

  public dispatchNotification(notification: ShellNotification) {
    this.store.dispatch(shellNotificationsActions.notificationAdded({ notification }));
  }
}
