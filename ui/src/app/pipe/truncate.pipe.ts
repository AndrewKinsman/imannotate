import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'truncate'
})
export class TruncatePipe implements PipeTransform {

  transform(value: string, limit = 25, completeWords = true, ellipsis = '…'): string {
    if (value.length <= limit) {
      return value;
    }

    const val = value.substr(0, limit);
    if (completeWords) {
      const r = new RegExp(`[ ,;]+`, 'g');
      while(r.test(val)) {
        limit = r.lastIndex
      }
    }

    return val.substr(0, limit) + ellipsis;
  }

}
