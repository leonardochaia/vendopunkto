import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ShellNotificationsDropdownComponent } from './shell-notifications-dropdown.component';

describe('ShellNotificationsDropdownComponent', () => {
  let component: ShellNotificationsDropdownComponent;
  let fixture: ComponentFixture<ShellNotificationsDropdownComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ShellNotificationsDropdownComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ShellNotificationsDropdownComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
