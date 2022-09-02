import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'lox-table-list-item',
  templateUrl: './table-list-item.component.html',
  styleUrls: ['./table-list-item.component.scss']
})
export class TableListItemComponent implements OnInit {
  @Input() podData: { transmitBytes: number[]; receiveBytes: number[] };
  @Input() podName: string;
  constructor() {}

  ngOnInit(): void {}
}
