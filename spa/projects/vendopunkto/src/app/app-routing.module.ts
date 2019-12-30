import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { InvoiceSearchContainerComponent } from './invoice-search-container/invoice-search-container.component';


const routes: Routes = [
  // {
  //   path: '',
  //   component:
  // },
  {
    path: 'invoices/:invoiceID',
    component: InvoiceSearchContainerComponent,
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
