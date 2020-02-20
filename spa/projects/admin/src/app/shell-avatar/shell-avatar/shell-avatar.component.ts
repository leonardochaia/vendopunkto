import { Component, OnInit, Input } from '@angular/core';

@Component({
  selector: 'adm-shell-avatar',
  templateUrl: './shell-avatar.component.html',
  styleUrls: ['./shell-avatar.component.scss']
})
export class ShellAvatarComponent {

  @Input()
  public imageURL: string;

  @Input()
  public initials: string;

  @Input()
  public name: string;

  public getInitials(name: string) {
    return name.split(' ').slice(0, 2).map(part => part.charAt(0)).join('');
  }

}
