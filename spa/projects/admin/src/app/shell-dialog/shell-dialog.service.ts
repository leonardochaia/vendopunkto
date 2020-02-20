import { Injectable } from '@angular/core';
import { MatDialog } from '@angular/material';
import { ShellDialogComponent } from './shell-dialog/shell-dialog.component';
import { ShellDialogConfig } from './models';

@Injectable({
  providedIn: 'root'
})
export class ShellDialogService {

  constructor(private readonly matDialog: MatDialog) { }

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
  }

}
