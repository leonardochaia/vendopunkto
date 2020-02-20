import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ShellNotificationCardComponent } from './shell-notification-card.component';

describe('ShellNotificationCardComponent', () => {
  let component: ShellNotificationCardComponent;
  let fixture: ComponentFixture<ShellNotificationCardComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ShellNotificationCardComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ShellNotificationCardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
