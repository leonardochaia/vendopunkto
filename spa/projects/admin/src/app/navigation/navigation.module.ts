import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { StoreModule } from '@ngrx/store';
import { EffectsModule } from '@ngrx/effects';
import * as fromNavigation from './+state/navigation.reducer';
import { NavigationEffects } from './+state/navigation.effects';
import { NavContainerComponent } from './nav-container/nav-container.component';
import { FlexLayoutModule } from '@angular/flex-layout';
import {
  MatToolbarModule,
  MatSidenavModule,
  MatListModule,
  MatIconModule,
  MatButtonModule
} from '@angular/material';
import { RouterModule } from '@angular/router';
import { ShellNotificationsModule } from '../shell-notifications/shell-notifications.module';

@NgModule({
  imports: [
    CommonModule,

    RouterModule,

    FlexLayoutModule,
    MatToolbarModule,
    MatSidenavModule,
    MatListModule,
    MatIconModule,
    MatButtonModule,

    StoreModule.forFeature(fromNavigation.NavigationFeatureKey, fromNavigation.reducer),
    EffectsModule.forFeature([NavigationEffects]),

    ShellNotificationsModule,
  ],
  declarations: [
    NavContainerComponent
  ],
  exports: [
    NavContainerComponent
  ]
})
export class NavigationModule { }
