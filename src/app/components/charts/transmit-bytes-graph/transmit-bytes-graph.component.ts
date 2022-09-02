import { Component } from '@angular/core';
import { StatsRecordModel } from '@lox/models/stats-record.model';
import { BaseChartComponent } from '../base-chart.component';

@Component({
  selector: 'lox-transmit-bytes-graph',
  templateUrl: './transmit-bytes-graph.component.html',
  styleUrls: ['./transmit-bytes-graph.component.scss']
})
export class TransmitBytesGraphComponent extends BaseChartComponent {
  //i OVERRIDE FROM BaseChartComponent
  protected updateData(data: StatsRecordModel): void {
    for (const [key, value] of Object.entries(data.pods)) {
      if (!!this._data[value.index]) {
        // Add Value to existing Entry
        this._data[value.index].data = value.transmitBytes;
      } else {
        // Create new Entry with data
        this._data[value.index] = {
          data: [...value.transmitBytes],
          label: key,
          borderColor: this.generateRandomHex()
        };
      }
    }
  }
}
