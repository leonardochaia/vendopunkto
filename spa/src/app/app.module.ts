import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { ReactiveFormsModule } from '@angular/forms';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { CommonModule } from '@angular/common';
import { InvoicePreviewComponent } from './invoice-preview/invoice-preview.component';

import { TextFieldModule } from '@angular/cdk/text-field';
import { FlexLayoutModule } from '@angular/flex-layout';
import {
  MatButtonModule,
  MatSnackBarModule,
  MatDialogModule,
  MatToolbarModule,
  MatIconModule,
  MatStepperModule
} from '@angular/material';

import { NgxKjuaModule } from 'ngx-kjua';
import { InvoiceSearchContainerComponent } from './invoice-search-container/invoice-search-container.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { InvoicePreviewHowToDialogComponent } from './invoice-preview-how-to-dialog/invoice-preview-how-to-dialog.component';

@NgModule({
  declarations: [
    AppComponent,
    InvoicePreviewComponent,
    InvoiceSearchContainerComponent,
    InvoicePreviewHowToDialogComponent
  ],
  imports: [
    CommonModule,
    HttpClientModule,
    ReactiveFormsModule,

    TextFieldModule,
    FlexLayoutModule,
    MatToolbarModule,
    MatIconModule,
    MatButtonModule,
    MatSnackBarModule,
    MatDialogModule,
    MatStepperModule,

    NgxKjuaModule,

    BrowserModule,
    AppRoutingModule,
    BrowserAnimationsModule
  ],
  providers: [],
  bootstrap: [AppComponent, InvoicePreviewHowToDialogComponent]
})
export class AppModule { }
