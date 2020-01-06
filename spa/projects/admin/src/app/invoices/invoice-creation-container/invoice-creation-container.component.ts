import { Component, OnInit } from '@angular/core';
import { InvoiceFacade } from '../+state/invoice.facade';
import { FormBuilder, Validators } from '@angular/forms';
import { startCreateInvoice } from '../+state/invoice.actions';

@Component({
  selector: 'adm-invoice-creation-container',
  templateUrl: './invoice-creation-container.component.html',
  styleUrls: ['./invoice-creation-container.component.scss']
})
export class InvoiceCreationContainerComponent implements OnInit {

  public readonly currencies$ = this.facade.currencies$;

  public readonly totalControl = this.fb.control(1, Validators.required);
  public readonly currencyControl = this.fb.control(null, Validators.required);

  public readonly form = this.fb.group({
    total: this.totalControl,
    currency: this.currencyControl
  });

  constructor(
    private readonly facade: InvoiceFacade,
    private readonly fb: FormBuilder) { }

  ngOnInit() {
  }

  public submit() {
    this.facade.dispatch(startCreateInvoice({
      invoice: this.form.value
    }));
  }

}
