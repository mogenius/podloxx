import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'lox-status-tile',
  templateUrl: './status-tile.component.html',
  styleUrls: ['./status-tile.component.scss']
})
export class StatusTileComponent implements OnInit {
  @Input() value: number;
  @Input() suffix: string;
  @Input() tileName: string;

  constructor() {}

  ngOnInit(): void {}
}
