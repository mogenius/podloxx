import { Component, OnDestroy, OnInit } from '@angular/core';
import { StatsService } from '@lox/services/stats.service';
import { debounceTime, merge, mergeMap, Subscription, switchMap, take } from 'rxjs';

@Component({
  selector: 'lox-dashboard-page',
  templateUrl: './dashboard-page.component.html',
  styleUrls: ['./dashboard-page.component.scss']
})
export class DashboardPageComponent implements OnInit, OnDestroy {
  private _subscriptions: Subscription;
  private _showRaw: boolean = false;
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
    this.refreshTotalStats();
    this.refreshFlow();
    this.refreshOverviewStats();
  }

  //ff refresh Total Stats every 10 seconds
  private refreshTotalStats(): void {
    this._subscriptions.add(
      this._statsService
        .statsTotal()
        .pipe(take(1), debounceTime(1000))
        .subscribe(() => {
          setTimeout(() => {
            this.refreshTotalStats();
          }, 20000);
        })
    );
  }

  //ff refresh Flow Stats every 2 seconds
  private refreshFlow(): void {
    this._subscriptions.add(
      this._statsService
        .statsFlow()
        .pipe(take(1), debounceTime(1000))
        .subscribe(() => {
          setTimeout(() => {
            this.refreshFlow();
          }, 10000);
        })
    );
  }

  //ff refresh overview stats every 10 seconds
  private refreshOverviewStats(): void {
    this._subscriptions.add(
      this._statsService
        .statsOverview()
        .pipe(take(1), debounceTime(1000))
        .subscribe(() => {
          setTimeout(() => {
            this.refreshOverviewStats();
          }, 10000);
        })
    );
  }

  public toggleRawData(state?: boolean): void {
    if (!!state) {
      this._showRaw = state;
    } else {
      this._showRaw = !this.showRaw;
    }

    if (this._showRaw) {
      document.body.classList.add('unscroll');
      document.getElementsByTagName('html')[0].classList.add('unscroll');
    } else {
      document.body.classList.remove('unscroll');
      document.getElementsByTagName('html')[0].classList.remove('unscroll');
    }
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

  get showRaw(): boolean {
    return this._showRaw;
  }
}
