import { Action, createReducer, on } from '@ngrx/store';
import * as NavigationActions from './navigation.actions';

export const NavigationFeatureKey = 'vpNavigation';

export interface NavigationState {
  isMobile: boolean;
}

export const initialState: NavigationState = {
  isMobile: true
};

const NavigationReducer = createReducer(
  initialState,

  on(NavigationActions.startMobileView, state => ({
    ...state,
    isMobile: true
  })),

  on(NavigationActions.startDesktopView, state => ({
    ...state,
    isMobile: false
  })),

);

export function reducer(state: NavigationState | undefined, action: Action) {
  return NavigationReducer(state, action);
}
