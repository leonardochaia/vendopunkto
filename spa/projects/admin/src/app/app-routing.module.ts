import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { InvoiceListContainerComponent } from './invoices/invoice-list-container/invoice-list-container.component';
import { InvoiceCreationContainerComponent } from './invoices/invoice-creation-container/invoice-creation-container.component';


const routes: Routes = [

  {
    path: 'invoices',
    component: InvoiceListContainerComponent,
  },
  {
    path: 'invoices/new',
    component: InvoiceCreationContainerComponent,
  },
  {
    path: '',
    pathMatch: 'prefix',
    redirectTo: 'invoices'
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
