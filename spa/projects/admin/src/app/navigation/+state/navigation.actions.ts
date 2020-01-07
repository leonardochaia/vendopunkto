import { createAction, props } from '@ngrx/store';

export const startMobileView = createAction(
  '[Navigation] Start Mobile View'
);

export const startDesktopView = createAction(
  '[Navigation] Start Desktop View'
);
