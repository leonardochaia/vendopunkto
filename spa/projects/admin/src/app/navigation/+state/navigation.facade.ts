import { Injectable } from '@angular/core';
import { NavigationState } from './navigation.reducer';
import { Store, Action } from '@ngrx/store';
import { selectIsMobile } from './navigation.selectors';

@Injectable({
  providedIn: 'root'
})
export class NavigationFacade {

  public readonly isMobile$ = this.store.select(selectIsMobile);

  constructor(private readonly store: Store<NavigationState>) { }

  public dispatch(action: Action) {
    this.store.dispatch(action);
  }
}
