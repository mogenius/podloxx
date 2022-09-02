import { Component, OnDestroy, OnInit } from '@angular/core';
import { StatsService } from '@lox/services/stats.service';
import { Subscription, take } from 'rxjs';

@Component({
  selector: 'lox-dashboard-page',
  templateUrl: './dashboard-page.component.html',
  styleUrls: ['./dashboard-page.component.scss']
})
export class DashboardPageComponent implements OnInit, OnDestroy {
  private _subscriptions: Subscription;
  constructor(private readonly _statsService: StatsService) {}

  ngOnInit(): void {
    if (this._subscriptions) {
      this._subscriptions.unsubscribe();
    }
    this._subscriptions = new Subscription();

    this.refreshData();
  }

  ngOnDestroy(): void {
    this._subscriptions.unsubscribe();
  }

  public refreshData(): void {
    this._subscriptions.add(
      this._statsService
        .stats()
        .pipe(take(1))
        .subscribe({
          next: () => {},
          error: (err) => {
            console.log(err);
          }
        })
    );
  }

  get totalPods(): number {
    return this._statsService.records.selectedpodNames.length > 0
      ? this._statsService.records.selectedpodNames.length
      : this._statsService.records.podNames.length;
  }

  get totalTransmitted(): number {
    let total = 0;
    for (const [key, value] of Object.entries(this._statsService.records.pods)) {
      total = total + value.transmitBytes[value.transmitBytes.length - 1];
    }
    return total;
  }

  get totalReceived(): number {
    let total = 0;
    for (const [key, value] of Object.entries(this._statsService.records.pods)) {
      total = total + value.receiveBytes[value.receiveBytes.length - 1];
    }
    return total;
  }
}
