import { createFeatureSelector, createSelector } from '@ngrx/store';
import * as fromPlugins from './plugins.reducer';

export const selectPluginState = createFeatureSelector<fromPlugins.PluginsState>(
  fromPlugins.PluginsFeatureKey
);

export const selectPlugins = createSelector(selectPluginState, s => s.plugins);

