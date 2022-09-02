import { Component, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { chartConfig } from '@lox/constants/chart.config';
import { StatsRecordModel } from '@lox/models/stats-record.model';
import { StatsService } from '@lox/services/stats.service';
import { ChartConfiguration, ChartDataset } from 'chart.js';
import { Moment } from 'moment';
import { BaseChartDirective } from 'ng2-charts';
import { Subscription } from 'rxjs';

@Component({
  template: ''
})
export abstract class BaseChartComponent implements OnInit, OnDestroy {
  @ViewChild(BaseChartDirective) chart?: BaseChartDirective;

  protected _subscriptions: Subscription;

  //! CHART OPTIONS
  protected _data: ChartDataset[] = [];
  protected _lineChartLabels: string[] = [];
  protected _chartConfig: ChartConfiguration = chartConfig;

  constructor(protected readonly _statsService: StatsService) {}

  ngOnInit(): void {
    /// INIT SUBSCRIPTION
    if (this._subscriptions) {
      this._subscriptions.unsubscribe();
    }
    this._subscriptions = new Subscription();

    this.initDataSubject();
    this.initFilterSubject();
  }

  ngOnDestroy(): void {
    /// CLEANUP SUBSCRIPTION
    if (this._subscriptions) {
      this._subscriptions.unsubscribe();
    }
  }

  protected initDataSubject(): void {
    this._subscriptions.add(
      this._statsService.records.lastUpdate.subscribe({
        next: (time: Moment) => {
          // NEW DATA ARRIVED
          const data = this._statsService.records;
          this._lineChartLabels.push(time.format('DD MMM HH:mm'));
          this.updateData(data);
          this.chart?.update();
        },
        error: (err) => {
          console.log(err);
        }
      })
    );
  }

  protected initFilterSubject(): void {
    this._subscriptions.add(
      this._statsService.records.updateFilter.subscribe({
        next: () => {
          // NEW DATA ARRIVED
          const data = this._statsService.records;
          this._data.forEach((set) => {
            set.data = [];
          });
          console.log(data);
          this.updateData(data);
          this.chart?.update();
        },
        error: (err) => {
          console.log(err);
        }
      })
    );
  }

  //!!! NEEDS TO BE IMPLEMENTED
  protected updateData(data: StatsRecordModel): void {}

  protected generateRandomHex(): string {
    return `#${Math.floor(Math.random() * 16777215).toString(16)}`;
  }

  public filterByByteRange(from: number, to: number): void {
    if (!!this._chartConfig?.options?.scales) {
      this._chartConfig!.options!.scales['y-axis-0']!.min = from;
      this._chartConfig!.options!.scales['y-axis-0']!.max = to;
    }
    this.chart?.render();
  }

  get chartConfig(): ChartConfiguration {
    return this._chartConfig;
  }
  get data(): ChartDataset[] {
    return this._data;
  }
  get lineChartLabels(): string[] {
    return this._lineChartLabels;
  }
}
