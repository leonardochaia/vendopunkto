import { Component, Input, TemplateRef, OnInit, Self, Optional } from '@angular/core';
import { ControlValueAccessor, FormControl, NgControl } from '@angular/forms';

@Component({
  selector: 'adm-config-selector-card',
  templateUrl: './config-selector-card.component.html',
  styleUrls: ['./config-selector-card.component.scss'],
})
export class ConfigSelectorCardComponent<T>
  implements OnInit, ControlValueAccessor {

  @Input()
  public title: string;

  @Input()
  public description: string;

  @Input()
  public items: T[];

  @Input()
  public itemNameTemplate: TemplateRef<T>;

  @Input()
  public itemIdProp = 'id';

  public innerControl = new FormControl(null, this.ngControl.validator);

  private initialValue: T;

  private onChange: (newValue: T) => void;

  constructor(
    @Self() @Optional()
    private readonly ngControl: NgControl,
  ) {
    ngControl.valueAccessor = this;
  }

  public ngOnInit() {
    this.initialValue = this.ngControl.value;
  }

  public reset() {
    this.innerControl.reset();
    this.innerControl.setValue(this.initialValue);
  }

  public startSaving() {
    this.innerControl.markAsPristine();
    this.ngControl.control.markAsPristine();
    this.initialValue = this.innerControl.value;
    this.onChange(this.innerControl.value);
  }

  registerOnTouched(fn: any): void {
  }

  registerOnChange(fn: any): void {
    this.onChange = fn;
  }

  writeValue(obj: any): void {
    this.innerControl.setValue(obj);
  }

  setDisabledState(isDisabled: boolean): void {
    if (isDisabled) {
      this.innerControl.disable();
    } else {
      this.innerControl.enable();
    }
  }
}
