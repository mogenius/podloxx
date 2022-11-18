import { Component, OnDestroy, OnInit } from '@angular/core';
import { StatsService } from '@lox/services/stats.service';
import { log } from 'console';
import { debounceTime, filter, merge, mergeMap, Subject, Subscription, switchMap, take, tap } from 'rxjs';

@Component({
  selector: 'lox-dashboard-page',
  templateUrl: './dashboard-page.component.html',
  styleUrls: ['./dashboard-page.component.scss']
})
export class DashboardPageComponent implements OnInit, OnDestroy {
  private _subscriptions: Subscription;
  private _showRaw: boolean = false;

  private _isReady: Subject<void> = new Subject<void>();
  private _ready: boolean = false;

  private _refreshingOverview: boolean = false;
  private _refreshingFlow: boolean = false;

  constructor(private readonly _statsService: StatsService) {}

  ngOnInit(): void {
    if (this._subscriptions) {
      this._subscriptions.unsubscribe();
    }
    this._subscriptions = new Subscription();

    this.refreshData();

    this._subscriptions.add(
      this._isReady
        .pipe(
          filter(() => {
            return this._statsService.records.podNames.length > 0;
          })
        )
        .subscribe(() => {
          this._ready = true;
        })
    );
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
        .subscribe((data) => {
          this._isReady.next();

          setTimeout(() => {
            this.refreshTotalStats();
          }, 10000);
        })
    );
  }

  //ff refresh Flow Stats every 2 seconds
  private refreshFlow(): void {
    this._subscriptions.add(
      this._statsService
        .statsFlow()
        .pipe(take(1), debounceTime(1000))
        .subscribe((data) => {
          this._isReady.next();

          this._refreshingFlow = true;
          setTimeout(() => {
            this._refreshingFlow = false;
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
          this._isReady.next();

          this._refreshingOverview = true;
          setTimeout(() => {
            this._refreshingOverview = false;
            this.refreshOverviewStats();
          }, 10000);
        })
    );
  }

  get refreshingOverview(): boolean {
    return this._refreshingOverview;
  }

  get refreshingFlow(): boolean {
    return this._refreshingFlow;
  }

  get ready(): boolean {
    return this._ready;
  }
}
