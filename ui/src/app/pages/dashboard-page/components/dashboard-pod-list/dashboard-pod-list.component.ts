import { Component, OnInit } from '@angular/core';
import { IPods } from '../../../../interfaces/pod.interface';
import { StatsService } from '../../../../services/stats.service';

@Component({
  selector: 'lox-dashboard-pod-list',
  templateUrl: './dashboard-pod-list.component.html',
  styleUrls: ['./dashboard-pod-list.component.scss']
})
export class DashboardPodListComponent implements OnInit {
  constructor(private readonly _statsService: StatsService) {}

  ngOnInit(): void {}

  get podNames(): string[] {
    return this._statsService.records.podNames;
  }

  get pods(): IPods {
    return this._statsService.records.podList;
  }
}
