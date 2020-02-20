import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ShellNotificationListContainerComponent } from './shell-notification-list-container.component';

describe('ShellNotificationListContainerComponent', () => {
  let component: ShellNotificationListContainerComponent;
  let fixture: ComponentFixture<ShellNotificationListContainerComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ShellNotificationListContainerComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ShellNotificationListContainerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
