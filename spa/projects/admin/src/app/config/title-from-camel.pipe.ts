import { Pipe, PipeTransform } from '@angular/core';
import { TitleCasePipe } from '@angular/common';

@Pipe({
  name: 'titleFromCamel'
})
export class TitleFromCamelPipe implements PipeTransform {
  private readonly titleCase = new TitleCasePipe();

  transform(value: any, ...args: any[]): any {
    return value && this.titleCase.transform(
      value.split('_').join(' ')
    );
  }

}
