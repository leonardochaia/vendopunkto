import { Injectable } from '@angular/core';
import { Subject, Observable } from 'rxjs';
import { MatDialog } from '@angular/material';
import { ShellDialogComponent } from './shell-dialog/shell-dialog.component';
import { ShellDialogConfig } from './models';

@Injectable({
  providedIn: 'root'
})
export class ShellDialogService {

  public readonly opened: Observable<ShellDialogConfig>;

  private readonly subject = new Subject<ShellDialogConfig>();

  constructor(private readonly matDialog: MatDialog) {
    this.opened = this.subject.asObservable();
  }

  public open(config: ShellDialogConfig) {
    this.matDialog.open(ShellDialogComponent, {
      position: {
        top: '64px',
        right: '0',
      },
      panelClass: 'shell-dialog-panel',
      data: config,
      autoFocus: false,
    });

    // this.subject.next(config);
  }

}
