import { createFeatureSelector, createSelector } from '@ngrx/store';
import * as fromConfig from './config.reducer';

export const selectConfigState = createFeatureSelector<fromConfig.ConfigState>(
  fromConfig.ConfigFeatureKey
);

export const selectCurrentConfig = createSelector(selectConfigState, s => s.current);

