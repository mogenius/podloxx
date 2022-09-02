import { Pipe, PipeTransform } from '@angular/core';

/**
 * example:
 * 102348213 -> 102.3 M or 4102348213 -> 4.1 B
 */
@Pipe({
  name: 'formatBytes'
})
export class formatBytes implements PipeTransform {
  transform(value: number, decimals: number = 2): unknown {
    if (value === 0) return '0 Bytes';
    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

    const i = Math.floor(Math.log(value) / Math.log(k));

    return parseFloat((value / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
  }
}
