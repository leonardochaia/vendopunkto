import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { StoreModule } from '@ngrx/store';
import { EffectsModule } from '@ngrx/effects';

import * as fromShellNotifications from './+state/shell-notifications.reducer';
import { ShellNotificationsEffects } from './+state/shell-notifications.effects';
import { ShellNotificationsDropdownComponent } from './shell-notifications-dropdown/shell-notifications-dropdown.component';
import { MatButtonModule, MatIconModule, MatMenuModule, MatSidenavModule, MatCardModule, MatProgressBarModule } from '@angular/material';
import { FlexLayoutModule } from '@angular/flex-layout';

@NgModule({
  imports: [
    CommonModule,

    FlexLayoutModule,
    MatButtonModule,
    MatIconModule,
    MatMenuModule,

    MatCardModule,
    MatProgressBarModule,

    MatSidenavModule,

    StoreModule.forFeature(fromShellNotifications.ShellNotificationsFeatureKey,
      fromShellNotifications.reducer),
    EffectsModule.forFeature([ShellNotificationsEffects]),
  ],
  declarations: [
    ShellNotificationsDropdownComponent
  ],
  exports: [
    ShellNotificationsDropdownComponent
  ],
})
export class ShellNotificationsModule { }
