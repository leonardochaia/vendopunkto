import { Component, OnInit } from '@angular/core';
import { ShellNotificationsFacade } from '../+state/shell-notifications.facade';
import { ShellOperationsFacade } from '../../shell-operations/+state/shell-operations.facade';
import { trigger, transition, style, animate } from '@angular/animations';

@Component({
  selector: 'adm-shell-notifications-dropdown',
  templateUrl: './shell-notifications-dropdown.component.html',
  styleUrls: ['./shell-notifications-dropdown.component.scss'],
  animations: [
    trigger('slideInOut', [
      transition(':enter', [
        style({ transform: 'translateX(100%)' }),
        animate('200ms ease-in', style({ transform: 'translateX(0%)' }))
      ]),
      transition(':leave', [
        animate('200ms ease-in', style({ transform: 'translateX(100%)' }))
      ])
    ])
  ]
})
export class ShellNotificationsDropdownComponent implements OnInit {

  public readonly notifications$ = this.notificationsFacade.notifications$;
  public readonly popoverNotification$ = this.notificationsFacade.popoverNotification$;
  public readonly hasPending$ = this.notificationsFacade.hasPendingNotifications$;
  public readonly operations$ = this.operationsFacade.operations$;

  constructor(
    private readonly notificationsFacade: ShellNotificationsFacade,
    private readonly operationsFacade: ShellOperationsFacade,
  ) { }

  ngOnInit() {
  }

}
