import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ShellSelectorDialogComponent } from './shell-selector-dialog.component';

describe('ShellSelectorDialogComponent', () => {
  let component: ShellSelectorDialogComponent;
  let fixture: ComponentFixture<ShellSelectorDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ShellSelectorDialogComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ShellSelectorDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
