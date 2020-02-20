import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ShellAvatarComponent } from './shell-avatar.component';

describe('ShellAvatarComponent', () => {
  let component: ShellAvatarComponent;
  let fixture: ComponentFixture<ShellAvatarComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ShellAvatarComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ShellAvatarComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
