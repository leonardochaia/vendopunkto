import { Injectable } from '@angular/core';
import { ConfigState } from './config.reducer';
import { Store, Action } from '@ngrx/store';
import { selectCurrentConfig } from './config.selectors';

@Injectable({
  providedIn: 'root'
})
export class ConfigFacade {

  public readonly current$ = this.store.select(selectCurrentConfig);

  constructor(private readonly store: Store<ConfigState>) { }

  public dispatch(action: Action) {
    this.store.dispatch(action);
  }
}
