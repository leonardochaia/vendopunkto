import { Component, Input } from '@angular/core';

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

  @Input()
  public size = '32px';

  public forceInitials = false;

  public getInitials(name: string) {
    return name.split(' ').slice(0, 2).map(part => part.charAt(0)).join('');
  }

  public onError(event) {
    // when an error happens fetching the image, fallback to initials
    this.forceInitials = true;
  }
}
