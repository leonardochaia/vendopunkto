import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ShellSelectorDialogComponent } from './shell-selector-dialog/shell-selector-dialog.component';
import { MatInputModule, MatCardModule, MatButtonModule, MatProgressSpinnerModule, MatIconModule } from '@angular/material';
import { ReactiveFormsModule } from '@angular/forms';
import { FlexLayoutModule } from '@angular/flex-layout';
import { ShellAvatarModule } from '../shell-avatar/shell-avatar.module';
import { ShellSelectorComponent } from './shell-selector/shell-selector.component';
import { ScrollingModule } from '@angular/cdk/scrolling';

@NgModule({
  imports: [
    CommonModule,
    ReactiveFormsModule,

    FlexLayoutModule,

    MatInputModule,
    MatButtonModule,
    MatIconModule,
    MatCardModule,
    MatProgressSpinnerModule,

    ScrollingModule,

    ShellAvatarModule,
  ],
  declarations: [
    ShellSelectorDialogComponent,
    ShellSelectorComponent,
  ],
  exports: [
    ShellSelectorComponent,
  ],
  entryComponents: [
    ShellSelectorDialogComponent,
  ]
})
export class ShellSelectorModule { }
