import { Action, createReducer, on } from '@ngrx/store';
import * as ConfigActions from './config.actions';

export const ConfigFeatureKey = 'vpConfig';

export interface ConfigState {
  current: { [k: string]: unknown };
  loading: boolean;
  error: string;
}

export const initialState: ConfigState = {
  current: null,
  loading: false,
  error: null
};

const NavigationReducer = createReducer(
  initialState,

  on(ConfigActions.loadConfig, (state) => ({
    ...state,
    loading: true
  })),

  on(ConfigActions.loadConfigSuccess, (state, action) => ({
    ...state,
    loading: false,
    current: action.config
  })),

  on(ConfigActions.loadConfigFailure, (state, action) => ({
    ...state,
    loading: false,
    error: action.error
  })),

);

export function reducer(state: ConfigState | undefined, action: Action) {
  return NavigationReducer(state, action);
}
