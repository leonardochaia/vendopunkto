import { createAction, props } from '@ngrx/store';

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
