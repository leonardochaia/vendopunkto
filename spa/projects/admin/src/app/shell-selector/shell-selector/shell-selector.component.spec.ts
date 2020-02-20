import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ShellSelectorComponent } from './shell-selector.component';

describe('ShellSelectorComponent', () => {
  let component: ShellSelectorComponent;
  let fixture: ComponentFixture<ShellSelectorComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ShellSelectorComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ShellSelectorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
