import { TitleFromCamelPipe } from './title-from-camel.pipe';

describe('TitleFromCamelPipe', () => {
  it('create an instance', () => {
    const pipe = new TitleFromCamelPipe();
    expect(pipe).toBeTruthy();
  });
});
