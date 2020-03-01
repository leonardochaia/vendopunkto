import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ConfigInvoiceEditComponent } from './config-invoice-edit.component';

describe('ConfigInvoiceEditComponent', () => {
  let component: ConfigInvoiceEditComponent;
  let fixture: ComponentFixture<ConfigInvoiceEditComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ConfigInvoiceEditComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ConfigInvoiceEditComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
