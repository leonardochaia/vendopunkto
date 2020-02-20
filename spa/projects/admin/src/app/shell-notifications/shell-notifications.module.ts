import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { StoreModule } from '@ngrx/store';
import { EffectsModule } from '@ngrx/effects';

import * as fromShellNotifications from './+state/shell-notifications.reducer';
import { ShellNotificationsEffects } from './+state/shell-notifications.effects';
import { ShellNotificationsDropdownComponent } from './shell-notifications-dropdown/shell-notifications-dropdown.component';
import { MatButtonModule, MatIconModule, MatMenuModule, MatSidenavModule, MatCardModule, MatProgressBarModule } from '@angular/material';
import { FlexLayoutModule } from '@angular/flex-layout';
import { ShellNotificationCardComponent } from './shell-notification-card/shell-notification-card.component';
import { ShellNotificationListContainerComponent } from './shell-notification-list-container/shell-notification-list-container.component';

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
    ShellNotificationsDropdownComponent,
    ShellNotificationCardComponent,
    ShellNotificationListContainerComponent
  ],
  exports: [
    ShellNotificationsDropdownComponent
  ],
  entryComponents: [
    ShellNotificationListContainerComponent,
  ]
})
export class ShellNotificationsModule { }
