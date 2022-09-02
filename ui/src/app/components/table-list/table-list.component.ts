import { Component, OnInit } from '@angular/core';
import { IStatsRecord } from '@lox/interfaces/stats-record.interface';
import { StatsService } from '@lox/services/stats.service';

@Component({
  selector: 'lox-table-list',
  templateUrl: './table-list.component.html',
  styleUrls: ['./table-list.component.scss']
})
export class TableListComponent implements OnInit {
  constructor(private readonly _statsService: StatsService) {}

  ngOnInit(): void {}

  get pods(): any {
    return Object.entries(this._statsService.records.pods);
  }
}
