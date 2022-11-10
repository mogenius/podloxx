import { Component, OnInit } from '@angular/core';
import { IPods } from '../../../../interfaces/pod.interface';
import { StatsService } from '../../../../services/stats.service';

@Component({
  selector: 'lox-dashboard-pod-list',
  templateUrl: './dashboard-pod-list.component.html',
  styleUrls: ['./dashboard-pod-list.component.scss']
})
export class DashboardPodListComponent implements OnInit {
  private _currentSortKey: string = 'podName';
  private _currentSortDirection: 'ASC' | 'DSC' = 'ASC';

  constructor(private readonly _statsService: StatsService) {}

  ngOnInit(): void {}

  sortPods(key: string): void {
    if (this._currentSortKey === key) {
      this._currentSortDirection = this._currentSortDirection === 'ASC' ? 'DSC' : 'ASC';
    } else {
      this._currentSortKey = key;
      this._currentSortDirection = 'ASC';
    }

    this._statsService.records.sortPods(key, this._currentSortDirection);
  }

  get podNames(): string[] {
    return this._statsService.records.sortedNameList ?? this._statsService.records.podNames;
  }

  get pods(): IPods {
    return this._statsService.records.podList;
  }

  get currentSortKey(): string {
    return this._currentSortKey;
  }

  get currentSortDirection(): 'ASC' | 'DSC' {
    return this._currentSortDirection;
  }
}
