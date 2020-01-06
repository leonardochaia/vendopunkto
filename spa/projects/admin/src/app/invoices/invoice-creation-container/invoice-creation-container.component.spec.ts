import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { InvoiceCreationContainerComponent } from './invoice-creation-container.component';

describe('InvoiceCreationContainerComponent', () => {
  let component: InvoiceCreationContainerComponent;
  let fixture: ComponentFixture<InvoiceCreationContainerComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ InvoiceCreationContainerComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(InvoiceCreationContainerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
