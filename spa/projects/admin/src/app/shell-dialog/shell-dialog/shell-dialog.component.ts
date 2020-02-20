import {
  Component, ViewChild, OnDestroy,
  ComponentFactoryResolver, Type, ComponentRef, Inject, AfterViewInit, ChangeDetectorRef
} from '@angular/core';
import { ShellDialogAnchorDirective } from '../shell-dialog-anchor.directive';
import { Subject } from 'rxjs';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material';
import { ShellDialogConfig } from '../models';

@Component({
  selector: 'adm-shell-dialog',
  templateUrl: './shell-dialog.component.html',
  styleUrls: ['./shell-dialog.component.scss']
})
export class ShellDialogComponent implements AfterViewInit, OnDestroy {

  public componentRef: ComponentRef<any>;

  @ViewChild(ShellDialogAnchorDirective, { static: true })
  insertionPoint: ShellDialogAnchorDirective;

  private readonly destroyed = new Subject();

  constructor(
    private readonly componentFactoryResolver: ComponentFactoryResolver,
    @Inject(MAT_DIALOG_DATA)
    public readonly data: ShellDialogConfig,
    private readonly changeDetectionRef: ChangeDetectorRef,
    private readonly dialogRef: MatDialogRef<ShellDialogComponent>,
  ) { }

  ngAfterViewInit() {
    this.loadChildComponent(this.data.component);
    this.changeDetectionRef.detectChanges();
  }

  close(result?: any) {
    this.dialogRef.close(result);
  }

  loadChildComponent(componentType: Type<any>) {
    const componentFactory = this.componentFactoryResolver.resolveComponentFactory(componentType);

    const viewContainerRef = this.insertionPoint.viewContainerRef;
    viewContainerRef.clear();

    this.componentRef = viewContainerRef.createComponent(componentFactory);
  }

  ngOnDestroy() {
    this.destroyed.next();
    this.destroyed.complete();
    if (this.componentRef) {
      this.componentRef.destroy();
    }
  }
}
