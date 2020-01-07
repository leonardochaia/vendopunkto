import { Injectable } from '@angular/core';
import { Actions, createEffect, ofType, OnInitEffects } from '@ngrx/effects';
import { catchError, map, concatMap, tap } from 'rxjs/operators';
import { EMPTY, of, fromEvent } from 'rxjs';

import * as NavigationActions from './navigation.actions';
import { Action } from '@ngrx/store';
import { MediaMatcher } from '@angular/cdk/layout';

@Injectable()
export class NavigationEffects {

  protected readonly query = this.media.matchMedia('(max-width: 600px)');

  onMediaChange$ = createEffect(() => fromEvent(this.query, 'change')
    .pipe(
      map(() => this.query.matches ?
        NavigationActions.startMobileView() : NavigationActions.startDesktopView())
    ));

  ngrxOnInitEffects(): Action {
    if (this.query.matches) {
      return NavigationActions.startMobileView();
    } else {
      return NavigationActions.startDesktopView();
    }
  }

  constructor(
    private readonly actions$: Actions,
    private readonly media: MediaMatcher) { }

}
