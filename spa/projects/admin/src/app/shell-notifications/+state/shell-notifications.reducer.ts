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
  popoverNotification: {
    opId: '42962d8f-56fa-4b3e-9fed-9b57cbc067a5',
    type: 'operation',
    title: 'Update configuration',
    message: 'Changed "Default Invoice Currency" to "BTC"',
    date: 1581988334610
  } as any,
  notifications: [
    {
      opId: '42962d8f-56fa-4b3e-9fed-9b57cbc067a5',
      type: 'operation',
      title: 'Update configuration',
      message: 'Changed "Default Invoice Currency" to "BTC"',
      date: 1581988334610
    } as any,
    {
      opId: 'ok-56fa-4b3e-9fed-9b57cbc067a5',
      type: 'operation',
      title: 'Update configuration',
      message: 'Changed "Default Invoice Currency" to "BTC"',
      date: 1281988324210
    },
    {
      opId: 'fail-56fa-4b3e-9fed-9b57cbc067a5',
      type: 'operation',
      title: 'Update configuration',
      message: 'Changed "Default Invoice Currency" to "BTC"',
      date: 1581988334410
    },
  ],
  // popoverNotification: null,
  // notifications: [],
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
);

export function reducer(state: ShellNotificationsState | undefined, action: Action) {
  return ShellNotificationsReducer(state, action);
}
