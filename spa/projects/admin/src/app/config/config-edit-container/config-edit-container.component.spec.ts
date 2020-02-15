import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ConfigEditContainerComponent } from './config-edit-container.component';

describe('ConfigEditContainerComponent', () => {
  let component: ConfigEditContainerComponent;
  let fixture: ComponentFixture<ConfigEditContainerComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ConfigEditContainerComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ConfigEditContainerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
