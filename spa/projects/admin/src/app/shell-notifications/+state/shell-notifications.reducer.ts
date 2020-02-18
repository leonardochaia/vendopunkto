import { Action, createReducer, on } from '@ngrx/store';
import * as ShellNotificationsActions from './shell-notifications.actions';
import { ShellNotification } from '../model';

export const ShellNotificationsFeatureKey = 'vpShellNotifications';

export interface ShellNotificationsState {
  notifications: ShellNotification[];
  popoverNotification: ShellNotification;
  error: string;
}

export const initialState: ShellNotificationsState = {
  popoverNotification: null,
  notifications: [],
  error: null
};

const ShellNotificationsReducer = createReducer(
  initialState,

  on(ShellNotificationsActions.notificationAdded, (state, action) => ({
    ...state,
    notifications: [action.notification, ...state.notifications],
    popoverNotification: action.notification,
  })),

  on(ShellNotificationsActions.notificationPopoverCleared, (state) => ({
    ...state,
    popoverNotification: undefined
  })),

  on(ShellNotificationsActions.dismissNotification, (state, action) => ({
    ...state,
    notifications: state.notifications.filter(n => n.id !== action.id),
    popoverNotification: state.popoverNotification
      && state.popoverNotification.id === action.id ? null : state.popoverNotification
  }))
);

export function reducer(state: ShellNotificationsState | undefined, action: Action) {
  return ShellNotificationsReducer(state, action);
}
