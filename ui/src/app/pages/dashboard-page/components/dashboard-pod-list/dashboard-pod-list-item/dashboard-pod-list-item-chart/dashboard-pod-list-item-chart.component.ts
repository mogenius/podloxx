import { Component, Input, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { StatsRecordModel } from '@lox/models/stats-record.model';
import { StatsService } from '@lox/services/stats.service';
import { ChartConfiguration, ChartDataset } from 'chart.js';
import { BaseChartDirective } from 'ng2-charts';
import { Subscription } from 'rxjs';
import { tableChartConfig } from './table-chart-settings';

@Component({
  selector: 'lox-dashboard-pod-list-item-chart',
  templateUrl: './dashboard-pod-list-item-chart.component.html',
  styleUrls: ['./dashboard-pod-list-item-chart.component.scss']
})
export class DashboardPodListItemChartComponent implements OnInit, OnDestroy {
  @Input() key: string = '';
  @ViewChild(BaseChartDirective) chart?: BaseChartDirective;

  protected _subscriptions: Subscription;

  //! CHART OPTIONS
  protected _data: ChartDataset[] = [];
  protected _lineChartLabels: string[] = [];
  protected _chartConfig: ChartConfiguration = tableChartConfig;

  constructor(protected readonly _statsService: StatsService) {}

  ngOnInit(): void {
    /// INIT SUBSCRIPTION
    if (this._subscriptions) {
      this._subscriptions.unsubscribe();
    }
    this._subscriptions = new Subscription();

    this.initUpdateEvent();
  }

  ngOnDestroy(): void {
    /// CLEANUP SUBSCRIPTION
    if (this._subscriptions) {
      this._subscriptions.unsubscribe();
    }
  }

  protected initUpdateEvent(): void {
    this._subscriptions.add(
      this._statsService.records.updateFlowEvent.subscribe({
        next: () => {
          // NEW DATA ARRIVED

          this.updateData(this._statsService.records);
          this.chart?.update();
        },
        error: (err) => {
          console.log(err);
        }
      })
    );
  }

  //!!! NEEDS TO BE IMPLEMENTED
  protected updateData(data: StatsRecordModel): void {
    const newData = data.podList[this.key].records[data.podList[this.key].records.length - 1];
    const index = data.podList[this.key].index;

    if (!!this._data[0]) {
      // Add Value to existing Entry
      this._data[0].data.push(newData.receivedBytes);
      this._lineChartLabels.push(newData.timeStamp.toLocaleString());
    } else {
      // generate Array of Transferred Bytes
      const totalBytesArray = data.podList[this.key].records.map((record) => record.receivedBytes);
      const dataLabelArray = data.podList[this.key].records.map((record) => record.timeStamp.toLocaleString());

      if (totalBytesArray.length === 1) {
        // Create new Entry with data
        this._data[0] = {
          data: [...totalBytesArray, ...totalBytesArray],
          label: this.key,
          borderColor: '#009bc5',
          fill: true,
          backgroundColor: 'rgba(0, 155, 197, 0.1)'
        };

        this._lineChartLabels = [...dataLabelArray, ...dataLabelArray];
      } else {
        // Create new Entry with data
        this._data[0] = {
          data: [...totalBytesArray],
          label: this.key,
          borderColor: '#009bc5',
          fill: true,
          backgroundColor: 'rgba(0, 155, 197, 0.1)'
        };

        this._lineChartLabels = [...dataLabelArray];
      }
    }
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
