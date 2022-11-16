import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'lox-double-status-tile',
  templateUrl: './double-status-tile.component.html',
  styleUrls: ['./double-status-tile.component.scss']
})
export class DoubleStatusTileComponent implements OnInit {
  @Input() valueIn: number;
  @Input() suffixIn: string;
  @Input() tileName: string;

  @Input() valueOut: number;
  @Input() suffixOut: string;

  constructor() {}

  ngOnInit(): void {}
}
