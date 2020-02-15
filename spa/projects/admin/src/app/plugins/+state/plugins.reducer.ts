import { Action, createReducer, on } from '@ngrx/store';
import * as PluginActions from './plugins.actions';
import { GetPluginResult } from 'shared';

export const PluginsFeatureKey = 'vpPlugins';

export interface PluginsState {
  plugins: GetPluginResult[];
  loading: boolean;
  error: string;
}

export const initialState: PluginsState = {
  plugins: [],
  loading: false,
  error: null
};

const NavigationReducer = createReducer(
  initialState,

  on(PluginActions.loadPlugins, (state) => ({
    ...state,
    loading: true
  })),

  on(PluginActions.loadPluginsSuccess, (state, action) => ({
    ...state,
    loading: false,
    plugins: action.plugins
  })),

  on(PluginActions.loadPluginsFailure, (state, action) => ({
    ...state,
    loading: false,
    error: action.error
  })),

);

export function reducer(state: PluginsState | undefined, action: Action) {
  return NavigationReducer(state, action);
}
