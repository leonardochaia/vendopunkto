import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { InvoiceListContainerComponent } from './invoice-list-container.component';

describe('InvoiceListContainerComponent', () => {
  let component: InvoiceListContainerComponent;
  let fixture: ComponentFixture<InvoiceListContainerComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ InvoiceListContainerComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(InvoiceListContainerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
