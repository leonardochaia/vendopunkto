import { createAction, props } from '@ngrx/store';
import { ShellNotification } from '../model';

export const notificationAdded = createAction(
  '[Notifications] Notification Added',
  props<{ notification: ShellNotification }>()
);

export const notificationPopoverCleared = createAction(
  '[Notifications] Notification Popover Cleared'
);

export const dismissNotification = createAction(
  '[Notifications] Dismiss Notification',
  props<{ id: string }>(),
);
