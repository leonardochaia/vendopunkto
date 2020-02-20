import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatSidenavModule, MatButtonModule, MatIconModule, MatDialogModule } from '@angular/material';
import { FlexLayoutModule } from '@angular/flex-layout';
import { ShellDialogComponent } from './shell-dialog/shell-dialog.component';
import { ShellDialogAnchorDirective } from './shell-dialog-anchor.directive';

@NgModule({
  imports: [
    CommonModule,

    FlexLayoutModule,
    MatSidenavModule,
    MatButtonModule,
    MatIconModule,
    MatDialogModule,
  ],
  declarations: [
    ShellDialogComponent,
    ShellDialogAnchorDirective,
  ],
  entryComponents: [
    ShellDialogComponent
  ]
})
export class ShellDialogModule { }
