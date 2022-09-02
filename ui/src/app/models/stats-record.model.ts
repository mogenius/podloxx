import { IStatsRecord } from '@lox/interfaces/stats-record.interface';
import { IStatsResponse } from '@lox/interfaces/stats-response.interface';
import { StatsService } from '@lox/services/stats.service';
import moment, { Moment } from 'moment';
import { BehaviorSubject, ReplaySubject } from 'rxjs';

export class StatsRecordModel {
  private _recordCount: number = 0;
  private _podCount: number = 0;
  private _lastUpdate: ReplaySubject<Moment> = new ReplaySubject<Moment>();
  private _updateFilter: ReplaySubject<void> = new ReplaySubject<void>();
  private _record: IStatsRecord = {};
  private _filteredRecord: IStatsRecord = {};

  constructor() {}

  public addRecord(rawData: IStatsResponse[]) {
    rawData.forEach((record) => {
      if (!!this._record[record.podName]) {
        this._record[record.podName].receiveBytes.push(+record.receiveBytes);
        this._record[record.podName].transmitBytes.push(+record.transmitBytes);
      } else {
        this._record[record.podName] = {
          index: this._podCount++,
          // If Data appears later then the rest, fill the previous values with 0
          receiveBytes: [...Array.from(new Array(this._recordCount), () => 0), +record.receiveBytes],
          transmitBytes: [...Array.from(new Array(this._recordCount), () => 0), +record.transmitBytes]
        };
      }
    });
    this._lastUpdate.next(moment());

    this._recordCount++;
  }

  public selectPod(pod: string): void {
    this._filteredRecord[pod] = this._record[pod];
    this._updateFilter.next();
  }

  public deSelectPod(pod: string): void {
    delete this._filteredRecord[pod];
    this._updateFilter.next();
  }

  public clearSelection(): void {
    Object.keys(this._filteredRecord).forEach((pod: string) => {
      delete this._filteredRecord[pod];
    });
    this._updateFilter.next();
  }

  get podNames(): string[] {
    return Object.keys(this._record);
  }

  get selectedpodNames(): string[] {
    return Object.keys(this._filteredRecord);
  }

  get pods(): IStatsRecord {
    return Object.keys(this._filteredRecord).length > 0 ? this._filteredRecord : this._record;
  }

  get lastUpdate(): ReplaySubject<Moment> {
    return this._lastUpdate;
  }

  get updateFilter(): ReplaySubject<void> {
    return this._updateFilter;
  }
}
