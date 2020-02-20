import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ShellAvatarComponent } from './shell-avatar/shell-avatar.component';

@NgModule({
  imports: [
    CommonModule
  ],
  declarations: [
    ShellAvatarComponent,
  ],
  exports: [
    ShellAvatarComponent
  ]
})
export class ShellAvatarModule { }
