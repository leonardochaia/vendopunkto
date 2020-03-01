import { Component, OnInit, Input, Self, Output, EventEmitter, Optional } from '@angular/core';
import { ControlValueAccessor, NgControl } from '@angular/forms';
import {
  ShellSelectorItem, ShellSelectorDialogComponent,
  ShellSelectorDialogData
} from '../shell-selector-dialog/shell-selector-dialog.component';
import { isArray, isString, isNumber } from 'util';
import { ShellDialogService } from '../../shell-dialog/shell-dialog.service';
import { Observable, BehaviorSubject, combineLatest } from 'rxjs';
import { map } from 'rxjs/operators';

@Component({
  selector: 'adm-shell-selector',
  templateUrl: './shell-selector.component.html',
  styleUrls: ['./shell-selector.component.scss'],
})
export class ShellSelectorComponent implements OnInit, ControlValueAccessor {

  /** This is the list of all available items. */
  @Input()
  public items$: Observable<ShellSelectorItem[]>;

  /** The title for the shell-dialog */
  @Input()
  public title: string;

  /** When the dialog is opened this event is fired */
  @Output()
  public dialogOpened = new EventEmitter<void>();

  public selectedItems$: Observable<ShellSelectorItem[]>;

  public isDisabled = false;

  protected currentItemsSubject = new BehaviorSubject<(string | number)[]>([]);

  protected onChanged: (items: (string | number)[]) => void;

  constructor(
    @Self() @Optional()
    private readonly ngControl: NgControl,
    private readonly shellDialog: ShellDialogService
  ) {
    // instead of providing the CVA through DI, we're inverting the operations
    // and injecting the ngControl and setting the CVA manually.
    // this allows us to get the AbstractController currently in use. 
    ngControl.valueAccessor = this;
  }

  public ngOnInit() {
    // the selector model's is the list of IDs
    // whenever the model or the items$ change, keep the view up to date
    this.selectedItems$ = combineLatest(this.currentItemsSubject, this.items$)
      .pipe(
        map(([ids, items]) => {
          const r: ShellSelectorItem[] = [];
          for (const id of ids) {
            const item = items.find(i => i.id === id);
            if (item) {
              r.push(item);
            }
          }
          return r;
        })
      );
  }

  public addItems() {
    this.dialogOpened.emit();
    this.shellDialog.open({
      component: ShellSelectorDialogComponent,
      title: this.title,
      extra: {
        formControl: this.ngControl.control,
        items$: this.items$,
      } as ShellSelectorDialogData
    });
  }

  public removeItem(item: ShellSelectorItem) {
    const n = [...this.currentItemsSubject.value];
    n.splice(
      this.currentItemsSubject.value.indexOf(item.id), 1);
    this.currentItemsSubject.next(n);
    this.onChanged(this.currentItemsSubject.value);
  }

  public registerOnTouched(fn: () => void): void {
  }

  public registerOnChange(fn: (items: (string | number)[]) => void): void {
    this.onChanged = fn;
  }

  public writeValue(items: (string | number)[]): void {
    if (isArray(items)) {
      if (!items.length) {
        this.currentItemsSubject.next([]);
        return;
      }

      if (isString(items[0]) || isNumber(items[0])) {
        this.currentItemsSubject.next(items);
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
