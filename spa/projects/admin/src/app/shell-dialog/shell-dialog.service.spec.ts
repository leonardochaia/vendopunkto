import { TestBed } from '@angular/core/testing';

import { ShellDialogService } from './shell-dialog.service';

describe('ShellDialogService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: ShellDialogService = TestBed.get(ShellDialogService);
    expect(service).toBeTruthy();
  });
});
