import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { ReactiveFormsModule } from '@angular/forms';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { CommonModule } from '@angular/common';
import { InvoicePreviewComponent } from './invoice-preview/invoice-preview.component';

import { NgxKjuaModule } from 'ngx-kjua';
import { InvoiceSearchContainerComponent } from './invoice-search-container/invoice-search-container.component';

@NgModule({
  declarations: [
    AppComponent,
    InvoicePreviewComponent,
    InvoiceSearchContainerComponent
  ],
  imports: [
    CommonModule,
    HttpClientModule,
    ReactiveFormsModule,

    NgxKjuaModule,

    BrowserModule,
    AppRoutingModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
