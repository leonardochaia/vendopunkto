import { createFeatureSelector, createSelector } from '@ngrx/store';
import * as fromNavigation from './navigation.reducer';

export const selectNavigationState = createFeatureSelector<fromNavigation.NavigationState>(
  fromNavigation.NavigationFeatureKey
);

export const selectIsMobile = createSelector(selectNavigationState, s => s.isMobile);

