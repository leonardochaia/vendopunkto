import { Component, OnInit, Input, Self } from '@angular/core';
import { ControlValueAccessor, NgControl } from '@angular/forms';
import {
  ShellSelectorItem, ShellSelectorDialogComponent,
  ShellSelectorDialogData
} from '../shell-selector-dialog/shell-selector-dialog.component';
import { isObject, isArray, isString, isNumber } from 'util';
import { ShellDialogService } from '../../shell-dialog/shell-dialog.service';
import { of } from 'rxjs';

@Component({
  selector: 'adm-shell-selector',
  templateUrl: './shell-selector.component.html',
  styleUrls: ['./shell-selector.component.scss'],
})
export class ShellSelectorComponent implements OnInit, ControlValueAccessor {

  @Input()
  public items: ShellSelectorItem[];

  @Input()
  public title: string;

  public currentItemIds: (string | number)[] = [];

  public get selectedItems() {
    return this.items.filter(i => this.currentItemIds.indexOf(i.id) >= 0);
  }

  public isDisabled = false;

  protected onChanged: (items: (string | number)[]) => void;

  constructor(
    @Self()
    private readonly ngControl: NgControl,
    private readonly shellDialog: ShellDialogService
  ) {
    ngControl.valueAccessor = this;
  }

  public addItems() {
    this.shellDialog.open({
      component: ShellSelectorDialogComponent,
      title: this.title,
      extra: {
        formControl: this.ngControl.control,
        items$: of(this.items)
      } as ShellSelectorDialogData
    });
  }

  public removeItem(item: ShellSelectorItem) {
    this.currentItemIds.splice(this.currentItemIds.indexOf(item.id), 1);
    this.onChanged(this.currentItemIds);
  }

  public ngOnInit() {
  }

  public registerOnTouched(fn: () => void): void {
  }

  public registerOnChange(fn: (items: (string | number)[]) => void): void {
    this.onChanged = fn;
  }

  public writeValue(items: (string | number)[]): void {
    if (isArray(items)) {
      if (!items.length) {
        this.currentItemIds = [];
        return;
      }

      if (isString(items[0]) || isNumber(items[0])) {
        this.currentItemIds = items as string[];
      } else {
        throw new Error('Expected items to be an string | number array');
      }
    } else {
      throw new Error('Expected items to be an string | number array');
    }
  }

  public setDisabledState(isDisabled: boolean): void {
    this.isDisabled = isDisabled;
  }

}
