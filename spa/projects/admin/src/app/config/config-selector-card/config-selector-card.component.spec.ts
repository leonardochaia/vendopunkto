import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ConfigSelectorCardComponent } from './config-selector-card.component';

describe('ConfigSelectorCardComponent', () => {
  let component: ConfigSelectorCardComponent;
  let fixture: ComponentFixture<ConfigSelectorCardComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ConfigSelectorCardComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ConfigSelectorCardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
