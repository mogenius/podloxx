import { Component, OnInit } from '@angular/core';
import { BaseChartComponent } from '../../../../../../components/charts/base-chart.component';
import { StatsRecordModel } from '../../../../../../models/stats-record.model';

@Component({
  selector: 'lox-dashboard-pod-list-item-chart',
  templateUrl: './dashboard-pod-list-item-chart.component.html',
  styleUrls: ['./dashboard-pod-list-item-chart.component.scss']
})
export class DashboardPodListItemChartComponent extends BaseChartComponent {
  //i OVERRIDE FROM BaseChartComponent
  protected updateData(data: StatsRecordModel): void {
    console.log(data);
    for (const [key, value] of Object.entries(data.pods)) {
      if (!!this._data[value.index]) {
        // Add Value to existing Entry
        this._data[value.index].data = value.receiveBytes;
      } else {
        // Create new Entry with data
        this._data[value.index] = {
          data: [...value.receiveBytes],
          label: key,
          borderColor: this.generateRandomHex()
        };
      }
    }
  }
}
