import { Component, OnInit, Inject, ElementRef, ViewChild, Renderer2, Renderer, AfterViewInit } from '@angular/core';
import { MAT_DIALOG_DATA } from '@angular/material';
import { ShellDialogConfig } from '../../shell-dialog/models';
import { Observable } from 'rxjs';
import { FormControl } from '@angular/forms';
import { map } from 'rxjs/operators';
import { ShellDialogComponent } from '../../shell-dialog/shell-dialog/shell-dialog.component';

export interface ShellSelectorDialogData {
  items$: Observable<ShellSelectorItem[]>;
  formControl: FormControl;
}

export interface ShellSelectorItem {
  id: string | number;
  name: string;
  initials?: string;
  description?: string;
  imageURL?: string;
}

interface ViewSelectorItem extends ShellSelectorItem {
  selected: boolean;
}

@Component({
  selector: 'adm-shell-selector-dialog',
  templateUrl: './shell-selector-dialog.component.html',
  styleUrls: ['./shell-selector-dialog.component.scss']
})
export class ShellSelectorDialogComponent implements AfterViewInit {

  public items$: Observable<ViewSelectorItem[]>;

  public formControl: FormControl;

  public searchFormControl = new FormControl();

  @ViewChild('searchInput', { static: true })
  searchInputElm: ElementRef;

  protected selectedItems: ViewSelectorItem[] = [];

  constructor(
    private readonly shellDialog: ShellDialogComponent,
    @Inject(MAT_DIALOG_DATA)
    readonly data: ShellDialogConfig,
  ) {
    const extra: ShellSelectorDialogData = data.extra;
    this.formControl = extra.formControl || new FormControl();

    if (!this.formControl.value) {
      this.formControl.setValue([], { emitEvent: false });
    }

    this.items$ = extra.items$.pipe(
      map(items => items
        .filter(i => this.formControl.value.indexOf(i.id) === -1)
        .map(i => ({ ...i, selected: false } as ViewSelectorItem)))
    );
  }

  public ngAfterViewInit() {
    this.searchInputElm.nativeElement.focus();
  }

  public toggleItemSelection(item: ViewSelectorItem) {
    if (item.selected) {
      this.selectedItems.splice(this.selectedItems.indexOf(item), 1);
    } else {
      this.selectedItems.push(item);
    }

    item.selected = !item.selected;
  }

  public confirm() {
    this.formControl.setValue(this.formControl.value.concat(this.selectedItems.map(i => i.id)));
    this.shellDialog.close();
  }

  public discard() {
    this.shellDialog.close();
  }
}
