import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { StatsService } from '@lox/services/stats.service';
import { Subscription, take } from 'rxjs';

@Component({
  selector: 'lox-headline',
  templateUrl: './headline.component.html',
  styleUrls: ['./headline.component.scss']
})
export class HeadlineComponent implements OnInit {
  @Input() refreshing: boolean = false;

  private _subscriptions: Subscription;
  private _showRaw: boolean = false;

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

  get showRaw(): boolean {
    return this._showRaw;
  }
}
