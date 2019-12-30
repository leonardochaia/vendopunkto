import { TestBed } from '@angular/core/testing';

import { VendopunktoApiService } from './vendopunkto-api.service';

describe('VendopunktoApiService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: VendopunktoApiService = TestBed.get(VendopunktoApiService);
    expect(service).toBeTruthy();
  });
});
