import { Injectable } from '@angular/core';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { catchError, map, concatMap } from 'rxjs/operators';
import { of } from 'rxjs';

import * as pluginActions from './plugins.actions';
import { Action } from '@ngrx/store';
import { VendopunktoApiService } from '../../vendopunkto-api.service';

@Injectable()
export class PluginEffects {

  onLoadConfig$ = createEffect(() => this.actions$
    .pipe(
      ofType(pluginActions.loadPlugins),

      concatMap(() => this.api.getPlugins()
        .pipe(
          map(plugins => pluginActions.loadPluginsSuccess({ plugins })),
          catchError(e => of(pluginActions.loadPluginsFailure({ error: e.message })))
        ))
    ));

  ngrxOnInitEffects(): Action {
    return pluginActions.loadPlugins();
  }

  constructor(
    private readonly actions$: Actions,
    private readonly api: VendopunktoApiService) { }

}
