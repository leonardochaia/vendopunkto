import { createAction, props } from '@ngrx/store';
import { GetPluginResult } from 'shared';

export const loadPlugins = createAction(
  '[Plugins] Load Plugins'
);

export const loadPluginsSuccess = createAction(
  '[Plugins] Load Plugins Success',
  props<{ plugins: GetPluginResult[] }>()
);

export const loadPluginsFailure = createAction(
  '[Plugins] Load Plugins',
  props<{ error: string }>()
);
