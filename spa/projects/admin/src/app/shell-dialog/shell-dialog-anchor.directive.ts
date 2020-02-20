import { Directive, ViewContainerRef } from '@angular/core';

@Directive({
  selector: '[admShellDialogAnchor]'
})
export class ShellDialogAnchorDirective {

  constructor(public viewContainerRef: ViewContainerRef) { }

}
