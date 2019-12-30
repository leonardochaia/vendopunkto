import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { InvoicePreviewHowToDialogComponent } from './invoice-preview-how-to-dialog.component';

describe('InvoicePreviewHowToDialogComponent', () => {
  let component: InvoicePreviewHowToDialogComponent;
  let fixture: ComponentFixture<InvoicePreviewHowToDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ InvoicePreviewHowToDialogComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(InvoicePreviewHowToDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
