import { Component, Input, Output, EventEmitter } from '@angular/core';
import { ShellNotification, OperationStartShellNotification } from '../model';
import { ShellOperationInstance } from '../../shell-operations/model';

@Component({
  selector: 'adm-shell-notification-card',
  templateUrl: './shell-notification-card.component.html',
  styleUrls: ['./shell-notification-card.component.scss']
})
export class ShellNotificationCardComponent {

  @Input()
  public notification: ShellNotification;

  @Output()
  public readonly dismissed = new EventEmitter<ShellNotification>();

  public get operation() {
    return (this.notification as OperationStartShellNotification).operation;
  }

  public dismiss() {
    this.dismissed.emit(this.notification);
  }
}
