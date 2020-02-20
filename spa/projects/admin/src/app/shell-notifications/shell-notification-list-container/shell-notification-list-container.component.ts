import { Component } from '@angular/core';
import { ShellNotificationsFacade } from '../+state/shell-notifications.facade';
import { dismissNotification } from '../+state/shell-notifications.actions';

@Component({
  selector: 'adm-shell-notification-list-container',
  templateUrl: './shell-notification-list-container.component.html',
  styleUrls: ['./shell-notification-list-container.component.scss'],
  animations: [
    // TODO: Animations got bugged after some refactoring.
    // trigger('slideOut', [
    //   transition('* => *', [
    //     query(':leave', [
    //       stagger(100, [
    //         animate('250ms ease-in', style({ transform: 'translateX(100%)' })),
    //         animate('250ms ease-in', style({ height: '0' }))
    //       ])
    //     ], { optional: true }),
    //   ])
    // ])
  ]
})
export class ShellNotificationListContainerComponent {

  public readonly notifications$ = this.notificationsFacade.notifications$;

  constructor(
    private readonly notificationsFacade: ShellNotificationsFacade,
  ) { }

  public dismiss(id: string) {
    this.notificationsFacade.dispatch(dismissNotification({ id }));
  }
}
