import { createAction, props } from '@ngrx/store';
import { UpdateConfigParams } from 'shared';

export const loadConfig = createAction(
  '[Config] Load Config'
);

export const loadConfigSuccess = createAction(
  '[Config] Load Config Success',
  props<{ config: any }>()
);

export const loadConfigFailure = createAction(
  '[Config] Load Config',
  props<{ error: string }>()
);

export const updateConfigStart = createAction(
  '[Config] Update config start',
  props<{
    config: UpdateConfigParams,
  }>()
);

export const updateConfigSuccess = createAction(
  '[Config] Update config success',
  props<{
    config: UpdateConfigParams,
  }>()
);

export const updateConfigFailure = createAction(
  '[Config] Update config failure',
  props<{ error: string }>()
);
