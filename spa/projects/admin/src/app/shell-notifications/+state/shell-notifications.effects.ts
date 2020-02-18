import { Injectable } from '@angular/core';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { delay, switchMap } from 'rxjs/operators';

import * as shellNotificationsActions from './shell-notifications.actions';
import { of } from 'rxjs';

@Injectable()
export class ShellNotificationsEffects {

  onNotification$ = createEffect(() => this.actions$.pipe(
    ofType(shellNotificationsActions.notificationAdded),
    switchMap(() => of(shellNotificationsActions.notificationPopoverCleared())
      .pipe(delay(3000))),
  ));

  constructor(
    private readonly actions$: Actions) { }

}
