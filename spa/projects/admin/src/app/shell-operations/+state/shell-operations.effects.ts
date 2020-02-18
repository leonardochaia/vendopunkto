import { Injectable } from '@angular/core';
import { Actions, createEffect } from '@ngrx/effects';
import { map } from 'rxjs/operators';

import * as shellOperationsActions from './shell-operations.actions';
import { Action } from 'rxjs/internal/scheduler/Action';

@Injectable()
export class ShellOperationsEffects {

  constructor(
    private readonly actions$: Actions) { }

}
