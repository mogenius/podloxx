import { Component, Input, OnInit } from '@angular/core';
import moment from 'moment';

@Component({
  selector: 'lox-table-list-item',
  templateUrl: './table-list-item.component.html',
  styleUrls: ['./table-list-item.component.scss']
})
export class TableListItemComponent implements OnInit {
  @Input() podData: { transmitBytes: number[]; receiveBytes: number[]; timeStamps: Date[] };
  @Input() podName: string;
  constructor() {}

  public formatShortDate(date: Date): string {
    return moment(date).format('HH:mm');
  }

  public formatFullDate(date: Date): string {
    return moment(date).format('YYYY.MM.DD - HH:mm');
  }

  ngOnInit(): void {}
}
