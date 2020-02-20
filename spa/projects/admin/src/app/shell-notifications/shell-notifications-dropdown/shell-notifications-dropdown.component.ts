import { Component, OnInit } from '@angular/core';
import { ShellNotificationsFacade } from '../+state/shell-notifications.facade';
import { ShellOperationsFacade } from '../../shell-operations/+state/shell-operations.facade';
import { trigger, transition, style, animate, query, stagger } from '@angular/animations';
import { dismissNotification } from '../+state/shell-notifications.actions';
import { ShellDialogService } from '../../shell-dialog/shell-dialog.service';
import { ShellNotificationListContainerComponent } from '../shell-notification-list-container/shell-notification-list-container.component';
import { map } from 'rxjs/operators';

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
  public readonly hasPending$ = this.notificationsFacade.notifications$
    .pipe(map(n => n && n.length));

  constructor(
    private readonly notificationsFacade: ShellNotificationsFacade,
    private readonly shellDialog: ShellDialogService,
  ) { }

  ngOnInit() {
  }

  public openNotifications() {
    this.shellDialog.open({
      title: 'Notifications',
      component: ShellNotificationListContainerComponent
    });
  }

  public dismiss(id: string) {
    this.notificationsFacade.dispatch(dismissNotification({ id }));
  }

}
