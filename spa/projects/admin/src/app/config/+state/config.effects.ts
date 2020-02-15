import { Injectable } from '@angular/core';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { catchError, map, concatMap } from 'rxjs/operators';
import { of } from 'rxjs';

import * as configActions from './config.actions';
import { Action } from '@ngrx/store';
import { VendopunktoApiService } from '../../vendopunkto-api.service';

@Injectable()
export class ConfigEffects {

  onLoadConfig$ = createEffect(() => this.actions$
    .pipe(
      ofType(configActions.loadConfig),

      concatMap(() => this.api.getConfig()
        .pipe(
          map(config => configActions.loadConfigSuccess({ config })),
          catchError(e => of(configActions.loadConfigFailure({ error: e.message })))
        ))
    ));

  ngrxOnInitEffects(): Action {
    return configActions.loadConfig();
  }

  constructor(
    private readonly actions$: Actions,
    private readonly api: VendopunktoApiService) { }

}
