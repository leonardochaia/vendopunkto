import { Injectable } from '@angular/core';
import { PluginsState } from './plugins.reducer';
import { Store, Action } from '@ngrx/store';
import { selectPlugins } from './plugins.selectors';

@Injectable({
  providedIn: 'root'
})
export class PluginsFacade {

  public readonly plugins$ = this.store.select(selectPlugins);

  constructor(private readonly store: Store<PluginsState>) { }

  public dispatch(action: Action) {
    this.store.dispatch(action);
  }
}
