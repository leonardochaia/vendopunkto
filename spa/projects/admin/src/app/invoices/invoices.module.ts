import { EffectsModule } from '@ngrx/effects';
import { InvoiceEffects } from './+state/invoice.effects';
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { StoreModule } from '@ngrx/store';
import * as fromInvoice from './+state/invoice.reducer';
import { InvoiceListContainerComponent } from './invoice-list-container/invoice-list-container.component';
import { MatCardModule, MatChipsModule, MatFormFieldModule, MatInputModule, MatButtonModule, MatSelectModule } from '@angular/material';
import { FlexLayoutModule } from '@angular/flex-layout';
import { InvoiceCreationContainerComponent } from './invoice-creation-container/invoice-creation-container.component';
import { ReactiveFormsModule } from '@angular/forms';

@NgModule({
  declarations: [InvoiceListContainerComponent, InvoiceCreationContainerComponent],
  imports: [
    CommonModule,
    ReactiveFormsModule,


    FlexLayoutModule,
    MatCardModule,
    MatChipsModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatSelectModule,

    StoreModule.forFeature(fromInvoice.InvoiceFeatureKey, fromInvoice.reducer),
    EffectsModule.forFeature([InvoiceEffects]),
  ]
})
export class InvoicesModule { }
