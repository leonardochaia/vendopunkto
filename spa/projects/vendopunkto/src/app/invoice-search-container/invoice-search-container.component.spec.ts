import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { InvoiceSearchContainerComponent } from './invoice-search-container.component';

describe('InvoiceSearchContainerComponent', () => {
  let component: InvoiceSearchContainerComponent;
  let fixture: ComponentFixture<InvoiceSearchContainerComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ InvoiceSearchContainerComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(InvoiceSearchContainerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
