import { Component, OnDestroy, OnInit } from '@angular/core';
import { StatsService } from '@lox/services/stats.service';
import { Subscription, take } from 'rxjs';

@Component({
  selector: 'lox-tile-overview',
  templateUrl: './tile-overview.component.html',
  styleUrls: ['./tile-overview.component.scss']
})
export class TileOverviewComponent implements OnInit {
  private _subscriptions: Subscription;

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
