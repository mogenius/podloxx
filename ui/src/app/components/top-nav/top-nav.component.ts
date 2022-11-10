import { Component, OnInit } from '@angular/core';
import { StatsService } from '@lox/services/stats.service';
import { Subscription, take } from 'rxjs';

@Component({
  selector: 'lox-top-nav',
  templateUrl: './top-nav.component.html',
  styleUrls: ['./top-nav.component.scss']
})
export class TopNavComponent implements OnInit {
  private _subscriptions: Subscription;
  public searchString: string;
  private _loadingStats: boolean = false;

  constructor(private readonly _statsService: StatsService) {}

  ngOnInit(): void {
    if (this._subscriptions) {
      this._subscriptions.unsubscribe();
    }
    this._subscriptions = new Subscription();
  }

  ngOnDestroy(): void {
    this._subscriptions.unsubscribe();
  }

  public selectPod(pod: string): void {
    if (this.isSelected(pod)) {
      this._statsService.records.deSelectPod(pod);
    } else {
      this._statsService.records.selectPod(pod);
    }
  }

  public isSelected(pod: string): boolean {
    return !!this._statsService.records.selectedpodNames.find((selectedPod) => selectedPod === pod);
  }

  public clearSelection(): void {
    this._statsService.records.clearSelection();
  }

  get podNameList(): string[] {
    if (this.searchString?.length > 0) {
      return this._statsService.records.podNames.filter((podName) => podName.includes(this.searchString));
    } else {
      return this._statsService.records.podNames;
    }
  }

  get selectedPodNameList(): string[] {
    return this._statsService.records.selectedpodNames;
  }

  get loadingStats(): boolean {
    return this._loadingStats;
  }
}
