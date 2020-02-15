import { Component, Input, forwardRef, ViewChild, TemplateRef } from '@angular/core';
import { Observable } from 'rxjs';
import { NG_VALUE_ACCESSOR, ControlValueAccessor, FormControlDirective, FormControl, ControlContainer } from '@angular/forms';

@Component({
  selector: 'adm-config-selector-card',
  templateUrl: './config-selector-card.component.html',
  styleUrls: ['./config-selector-card.component.scss'],
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => ConfigSelectorCardComponent),
      multi: true
    }
  ]
})
export class ConfigSelectorCardComponent<T> implements ControlValueAccessor {

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

  @ViewChild(FormControlDirective, { static: true })
  formControlDirective: FormControlDirective;

  @Input()
  formControl: FormControl;

  @Input()
  formControlName: string;

  /* get hold of FormControl instance no matter formControl or
  formControlName is given. If formControlName is given,
  then this.controlContainer.control is the parent FormGroup (or FormArray) instance. */
  get control() {
    return this.formControl || this.controlContainer.control.get(this.formControlName);
  }

  constructor(private controlContainer: ControlContainer) {
  }

  clearInput() {
    this.control.setValue('');
  }

  registerOnTouched(fn: any): void {
    this.formControlDirective.valueAccessor.registerOnTouched(fn);
  }

  registerOnChange(fn: any): void {
    this.formControlDirective.valueAccessor.registerOnChange(fn);
  }

  writeValue(obj: any): void {
    this.formControlDirective.valueAccessor.writeValue(obj);
  }

  setDisabledState(isDisabled: boolean): void {
    this.formControlDirective.valueAccessor.setDisabledState(isDisabled);
  }

}
