import { Component, OnInit, OnDestroy } from '@angular/core';
import { NavigationFacade } from '../+state/navigation.facade';
import { takeUntil } from 'rxjs/operators';
import { Subject } from 'rxjs';

@Component({
  selector: 'adm-nav-container',
  templateUrl: './nav-container.component.html',
  styleUrls: ['./nav-container.component.scss']
})
export class NavContainerComponent implements OnDestroy {

  public isMobile = true;

  protected destroyedSubject = new Subject();

  constructor(private readonly navigation: NavigationFacade) {
    this.navigation.isMobile$
      .pipe(
        takeUntil(this.destroyedSubject)
      )
      .subscribe((isMobile) => {
        this.isMobile = isMobile;
      });
  }

  ngOnDestroy() {
    this.destroyedSubject.next();
    this.destroyedSubject.complete();
  }

}
